package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/spiarh/gojo/pkg/core"
	"github.com/spiarh/gojo/pkg/util"
)

func Commit() (*cobra.Command, error) {
	var command = &cobra.Command{
		Use:               "commit",
		Short:             "Commit changes from an image directory",
		Example:           "gojo commit haproxy",
		RunE:              func(cmd *cobra.Command, args []string) error { return commit(cmd, args) },
		SilenceUsage:      true,
		TraverseChildren:  true,
		PersistentPreRunE: SetGlobalLogLevel,
	}

	if err := AddCommonPersistentFlags(command); err != nil {
		return nil, err
	}

	command.PersistentFlags().StringP(core.NameFlag, "n", "", "Name of the git commit author")
	command.PersistentFlags().StringP(core.EmailFlag, "e", "", "Email of the git commit author")

	if err := command.MarkPersistentFlagRequired(core.NameFlag); err != nil {
		return nil, err
	}
	if err := command.MarkPersistentFlagRequired(core.EmailFlag); err != nil {
		return nil, err
	}

	return command, nil
}

func getGitAuthorInfo(flagSet *pflag.FlagSet) (string, string, error) {
	var name, email string
	var err error

	if name, err = flagSet.GetString(core.NameFlag); err != nil {
		return name, email, err
	}
	if email, err = flagSet.GetString(core.EmailFlag); err != nil {
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
	log.Info().Bool(core.EnabledKey, opt.dryRun).Msg(core.DryRunFlag)

	if _, err := os.Stat(opt.buildFilePath); os.IsNotExist(err) {
		log.Warn().Str(core.FileKey, opt.buildFilePath).
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
		log.Info().Str(core.RepoKey, opt.imageDir).
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
		log.Info().Str(core.FileKey, buildFileRelPath).
			Msg("nothing to commit, file is not in status")
		return nil
	}

	if unmod := util.GitIsFileClean(status, buildFileRelPath); unmod {
		log.Info().Str(core.FileKey, buildFileRelPath).
			Msg("nothing to commit, file working tree clean")
		return nil
	}

	fmt.Println(buildFileRelPath)
	if unmod := util.GitIsFileUnmodifiedWorktree(status, buildFileRelPath); !unmod {
		log.Info().Str(core.FileKey, buildFileRelPath).
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
		log.Info().Str(core.MsgKey, msg).Msg("create commit message")

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
			log.Info().Str(core.HashKey, commit.String()).Msg("")

			obj, err := repo.CommitObject(commit)
			if err != nil {
				return err
			}
			log.Info().Str(core.CommitKey, obj.String()).Msg("")

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

	msg = fmt.Sprintf("[gojo] New build file, image=%s, tag=%s", build.Image.Name, build.Image.Tag)

	return msg, nil
}
