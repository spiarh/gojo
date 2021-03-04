package cmd

import (
	"fmt"
	"path"

	"github.com/lcavajani/gojo/pkg/core"
	"github.com/lcavajani/gojo/pkg/util"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var imageFQIN, fromImage, fromImageBuilder string

func New() *cobra.Command {
	var command = &cobra.Command{
		Use:               "scaffold",
		Short:             "Scaffold a new image project",
		Example:           "  gojo new github / gojo new alpine",
		SilenceUsage:      true,
		PersistentPreRunE: SetGlobalLogLevel,
	}

	AddCommonPersistentFlags(command)

	command.PersistentFlags().StringVar(&imageFQIN, imageFQINFlag, "", "Fully qualified image name (FQIN)")
	command.PersistentFlags().StringVar(&fromImage, fromImageFlag, "", "Image to use for the base image")
	command.PersistentFlags().StringVar(&fromImageBuilder, fromImageBuilderFlag, "", "Image to use for the builder")

	command.MarkPersistentFlagRequired(imageFQINFlag)
	command.MarkPersistentFlagRequired(fromImageFlag)

	newAlpineCommand.PersistentFlags().StringP(pkgFlag, "p", "", "Alpine package to use for the version")
	newAlpineCommand.MarkPersistentFlagRequired(pkgFlag)
	newAlpineCommand.PersistentFlags().StringP(repoFlag, "r", "", "Alpine repository to use for the package")
	newAlpineCommand.MarkPersistentFlagRequired(repoFlag)
	newAlpineCommand.PersistentFlags().StringP(versionIDFlag, "v", "", "Alpine version ID of the repository")
	newAlpineCommand.MarkPersistentFlagRequired(versionIDFlag)

	newGitHubCommand.PersistentFlags().StringP(ownerFlag, "o", "", "GitHub owner of the repository")
	newGitHubCommand.MarkPersistentFlagRequired(ownerFlag)
	newGitHubCommand.PersistentFlags().StringP(repoFlag, "r", "", "GitHub repository")
	newGitHubCommand.MarkPersistentFlagRequired(repoFlag)

	command.AddCommand(newAlpineCommand)
	command.AddCommand(newGitHubCommand)
	command.AddCommand(newSimpleCommand)

	return command
}

func scaffoldProject(command *cobra.Command, args []string) error {
	var err error
	providerType := command.Use

	flagSet := command.Flags()
	opt, err := getOptions(flagSet)
	if err != nil {
		return err
	}
	log.Info().Bool(enabledKey, opt.dryRun).Msg(dryRunFlag)

	// Objects
	build, err := core.NewBuild(imageFQIN, fromImage)
	if err != nil {
		return err
	}

	if fromImageBuilder != "" {
		fromImageBuilderMeta, err := core.NewImageFromFQIN(fromImageBuilder)
		if err != nil {
			return err
		}
		build.Spec.FromImageBuilder = &fromImageBuilderMeta
	}

	switch providerType {
	case alpine:
		err = addAlpineProvider(build, flagSet)
		if err != nil {
			return err
		}
	case github:
		err = addGitHubProvider(build, flagSet)
		if err != nil {
			return err
		}
	}

	containerfile, err := core.TemplateContainerfile(build)
	if err != nil {
		return err
	}
	log.Info().Str(nameKey, defaultBuildFileName).
		Msg("gojo image build file creation")
	fmt.Println(build)

	log.Info().Str(nameKey, defaultContainerfileName).
		Msg("Containerfile creation")
	fmt.Println(containerfile)

	if opt.dryRun {
		return nil
	}

	err = util.MakeDir(opt.imageDir, 0755)
	if err != nil {
		return err
	}

	build.WriteToFile(opt.buildFilePath)
	containerfile.WriteToFile(path.Join(opt.imageDir, defaultContainerfileName))

	return nil
}

// Simple
var newSimpleCommand = &cobra.Command{
	Use:          simple,
	RunE:         func(cmd *cobra.Command, args []string) error { return scaffoldProject(cmd, args) },
	SilenceUsage: true,
}

// Alpine
var newAlpineCommand = &cobra.Command{
	Use:          alpine,
	RunE:         func(cmd *cobra.Command, args []string) error { return scaffoldProject(cmd, args) },
	SilenceUsage: true,
}

func addAlpineProvider(build *core.Build, flagSet *pflag.FlagSet) error {
	var pkg, repo, versionId string
	var err error

	if repo, err = flagSet.GetString(repoFlag); err != nil {
		return err
	}
	if pkg, err = flagSet.GetString(pkgFlag); err != nil {
		return err
	}
	if versionId, err = flagSet.GetString(versionIDFlag); err != nil {
		return err
	}

	source := core.NewAlpineSource(alpine, pkg, repo, versionId)
	build.Spec.Sources = []core.Source{source}

	fact := core.NewFact(versionKey, "", alpine, core.VersionFactKind)
	build.Spec.BuildArgs = append(build.Spec.BuildArgs, versionKey)
	build.Spec.Facts = append(build.Spec.Facts, fact)
	build.Spec.TagFormat = core.TagFormatVersion

	return nil
}

// GitHub
var newGitHubCommand = &cobra.Command{
	Use:          github,
	RunE:         func(cmd *cobra.Command, args []string) error { return scaffoldProject(cmd, args) },
	SilenceUsage: true,
}

func addGitHubProvider(build *core.Build, flagSet *pflag.FlagSet) error {
	var repo, owner string
	var err error

	if repo, err = flagSet.GetString(repoFlag); err != nil {
		return err
	}
	if owner, err = flagSet.GetString(ownerFlag); err != nil {
		return err
	}

	source := core.NewGitHubSource(github, repo, owner)
	build.Spec.Sources = []core.Source{source}

	fact := core.NewFact(versionKey, "", github, core.VersionFactKind)
	build.Spec.BuildArgs = append(build.Spec.BuildArgs, versionKey)
	build.Spec.Facts = append(build.Spec.Facts, fact)
	build.Spec.TagFormat = core.TagFormatVersion

	return nil
}
