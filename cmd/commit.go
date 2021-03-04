package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/lcavajani/gojo/pkg/core"
	"github.com/lcavajani/gojo/pkg/util"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Commit() *cobra.Command {
	var command = &cobra.Command{
		Use:               "commit",
		Short:             "Commit changes from an image directory",
		Example:           "gojo commit haproxy",
		RunE:              func(cmd *cobra.Command, args []string) error { return commit(cmd, args) },
		SilenceUsage:      true,
		TraverseChildren:  true,
		PersistentPreRunE: SetGlobalLogLevel,
	}

	AddCommonPersistentFlags(command)

	command.PersistentFlags().StringP(nameFlag, "n", "", "Name of the git commit author")
	command.PersistentFlags().StringP(emailFlag, "e", "", "Email of the git commit author")

	command.MarkPersistentFlagRequired(nameFlag)
	command.MarkPersistentFlagRequired(emailFlag)

	return command
}

func getGitAuthorInfo(flagSet *pflag.FlagSet) (string, string, error) {
	var name, email string
	var err error

	if name, err = flagSet.GetString(nameFlag); err != nil {
		return name, email, err
	}
	if email, err = flagSet.GetString(emailFlag); err != nil {
		return name, email, err
	}

	return name, email, nil
}

func commit(command *cobra.Command, args []string) error {
	flagSet := command.Flags()

	name, email, err := getGitAuthorInfo(flagSet)
	if err != nil {
		return err
	}

	opt, err := getOptions(flagSet)
	if err != nil {
		return err
	}
	log.Info().Bool(enabledKey, opt.dryRun).Msg(dryRunFlag)

	if _, err := os.Stat(opt.buildFilePath); os.IsNotExist(err) {
		log.Warn().Str(fileKey, opt.buildFilePath).
			Msg("no build file found")
		return nil
	}

	repo, err := util.OpenGitRepo(opt.imageDir)
	if err != nil {
		return err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	status, err := wt.Status()
	if err != nil {
		return err
	}

	if status.IsClean() {
		log.Info().Str(repoKey, opt.imageDir).
			Msg("nothing to commit, working tree clean")
		return nil
	}

	root := wt.Filesystem.Root()
	buildFileRelPath, err := util.GetRelPathFromPathInTree(root, opt.buildFilePath)
	if err != nil {
		return err
	}

	// it file not in status output, this means it is
	// by default tracked and unmodified.
	if _, ok := (status)[buildFileRelPath]; !ok {
		log.Info().Str(fileKey, buildFileRelPath).
			Msg("nothing to commit, file is not in status")
		return nil
	}

	if unmod := util.GitIsFileClean(status, buildFileRelPath); unmod {
		log.Info().Str(fileKey, buildFileRelPath).
			Msg("nothing to commit, file working tree clean")
		return nil
	}

	fmt.Println(buildFileRelPath)
	if unmod := util.GitIsFileUnmodifiedWorktree(status, buildFileRelPath); !unmod {
		log.Info().Str(fileKey, buildFileRelPath).
			Msg("add file content to the index")
		if !opt.dryRun {
			if _, err := wt.Add(buildFileRelPath); err != nil {
				return err
			}
		}
	}

	if unmod := util.GitIsFileUnmodifiedStaging(status, buildFileRelPath); !unmod {
		log.Info().Msg("create new commit")

		msg, err := newCommitMessage(opt)
		if err != nil {
			return err
		}
		log.Info().Str(msgKey, msg).Msg("create commit message")

		if !opt.dryRun {
			commit, err := wt.Commit(msg, &git.CommitOptions{
				Author: &object.Signature{
					Name:  name,
					Email: email,
					When:  time.Now(),
				},
			})
			if err != nil {
				return err
			}
			log.Info().Str(hashKey, commit.String()).Msg("")

			obj, err := repo.CommitObject(commit)
			if err != nil {
				return err
			}
			log.Info().Str(commitKey, obj.String()).Msg("")

			log.Info().Msg("push to remote")
			if err = repo.Push(&git.PushOptions{}); err != nil {
				return err
			}
		}

	}

	return nil
}

func newCommitMessage(opt CommonOptions) (string, error) {
	var msg string
	build, err := core.NewBuildFromManifest(opt.buildFilePath)
	if err != nil {
		return msg, err
	}

	msg = fmt.Sprintf("[gojo] New build file, image=%s, tag=%s", build.Metadata.Name, build.Metadata.Tag)

	return msg, nil
}
