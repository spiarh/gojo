package cmd

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/lcavajani/gojo/pkg/core"
	"github.com/lcavajani/gojo/pkg/util"
)

var imageFQIN, fromImage, fromImageBuilder string

func Scaffold() (*cobra.Command, error) {
	var command = &cobra.Command{
		Use:               "scaffold",
		Short:             "Scaffold a new image project",
		Example:           "  gojo new github / gojo new alpine",
		SilenceUsage:      true,
		PersistentPreRunE: SetGlobalLogLevel,
	}

	if err := AddCommonPersistentFlags(command); err != nil {
		return nil, err
	}

	command.PersistentFlags().StringVar(&imageFQIN, core.ImageFQINFlag, "", "Fully qualified image name (FQIN)")
	command.PersistentFlags().StringVar(&fromImage, core.FromImageFlag, "", "Image to use for the base image")
	command.PersistentFlags().StringVar(&fromImageBuilder, core.FromImageBuilderFlag, "", "Image to use for the builder")

	if err := command.MarkPersistentFlagRequired(core.ImageFQINFlag); err != nil {
		return nil, err
	}
	if err := command.MarkPersistentFlagRequired(core.FromImageFlag); err != nil {
		return nil, err
	}

	newAlpineCommand.PersistentFlags().StringP(core.PkgFlag, "p", "", "Alpine package to use for the version")
	if err := newAlpineCommand.MarkPersistentFlagRequired(core.PkgFlag); err != nil {
		return nil, err
	}
	newAlpineCommand.PersistentFlags().StringP(core.RepoFlag, "r", "", "Alpine repository to use for the package")
	if err := newAlpineCommand.MarkPersistentFlagRequired(core.RepoFlag); err != nil {
		return nil, err
	}
	newAlpineCommand.PersistentFlags().StringP(core.VersionIDFlag, "v", "", "Alpine version ID of the repository")
	if err := newAlpineCommand.MarkPersistentFlagRequired(core.VersionIDFlag); err != nil {
		return nil, err
	}

	newGitHubCommand.PersistentFlags().StringP(core.OwnerFlag, "o", "", "GitHub owner of the repository")
	if err := newGitHubCommand.MarkPersistentFlagRequired(core.OwnerFlag); err != nil {
		return nil, err
	}
	newGitHubCommand.PersistentFlags().StringP(core.RepoFlag, "r", "", "GitHub repository")
	if err := newGitHubCommand.MarkPersistentFlagRequired(core.RepoFlag); err != nil {
		return nil, err
	}

	command.AddCommand(newAlpineCommand)
	command.AddCommand(newGitHubCommand)
	command.AddCommand(newSimpleCommand)

	return command, nil
}

func scaffoldProject(command *cobra.Command, args []string) error {
	var err error
	providerType := command.Use

	flagSet := command.Flags()
	opt, err := getOptions(flagSet)
	if err != nil {
		return err
	}
	log.Info().Bool(core.EnabledKey, opt.dryRun).Msg(core.DryRunFlag)

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
		build.Spec.FromImages = append(build.Spec.FromImages, core.FromImage{Image: fromImageBuilderMeta})
	}

	switch providerType {
	case core.Alpine:
		err = addAlpineProvider(build, flagSet)
		if err != nil {
			return err
		}
	case core.Github:
		err = addGitHubProvider(build, flagSet)
		if err != nil {
			return err
		}
	}

	containerfile, err := core.TemplateContainerfile(build)
	if err != nil {
		return err
	}
	log.Info().Str(core.NameKey, core.BuildFileName).
		Msg("gojo image build file creation")
	fmt.Println(build)

	log.Info().Str(core.NameKey, core.ContainerfileName).
		Msg("Containerfile creation")
	fmt.Println(containerfile)

	if opt.dryRun {
		return nil
	}

	err = util.MakeDir(opt.imageDir, 0755)
	if err != nil {
		return err
	}

	if err := build.WriteToFile(opt.buildFilePath); err != nil {
		return err
	}

	if err := containerfile.WriteToFile(opt.containerFilePath); err != nil {
		return err
	}

	return nil
}

// Simple
var newSimpleCommand = &cobra.Command{
	Use:          core.Simple,
	RunE:         func(cmd *cobra.Command, args []string) error { return scaffoldProject(cmd, args) },
	SilenceUsage: true,
}

// Alpine
var newAlpineCommand = &cobra.Command{
	Use:          core.Alpine,
	RunE:         func(cmd *cobra.Command, args []string) error { return scaffoldProject(cmd, args) },
	SilenceUsage: true,
}

func addAlpineProvider(build *core.Build, flagSet *pflag.FlagSet) error {
	var pkg, repo, versionId string
	var err error

	if repo, err = flagSet.GetString(core.RepoFlag); err != nil {
		return err
	}
	if pkg, err = flagSet.GetString(core.PkgFlag); err != nil {
		return err
	}
	if versionId, err = flagSet.GetString(core.VersionIDFlag); err != nil {
		return err
	}

	source := core.NewAlpineSource(core.Alpine, pkg, repo, versionId)
	build.Spec.Sources = []core.Source{source}

	fact := core.NewFact(core.VersionKey, "", core.Alpine, core.VersionFactKind)
	build.Spec.BuildArgs = append(build.Spec.BuildArgs, core.VersionKey)
	build.Spec.Facts = append(build.Spec.Facts, fact)
	build.Spec.TagFormat = core.TagFormatVersion

	return nil
}

// GitHub
var newGitHubCommand = &cobra.Command{
	Use:          core.Github,
	RunE:         func(cmd *cobra.Command, args []string) error { return scaffoldProject(cmd, args) },
	SilenceUsage: true,
}

func addGitHubProvider(build *core.Build, flagSet *pflag.FlagSet) error {
	var repo, owner string
	var err error

	if repo, err = flagSet.GetString(core.RepoFlag); err != nil {
		return err
	}
	if owner, err = flagSet.GetString(core.OwnerFlag); err != nil {
		return err
	}

	source := core.NewGitHubSource(core.Github, repo, owner)
	build.Spec.Sources = []core.Source{source}

	fact := core.NewFact(core.VersionKey, "", core.Github, core.VersionFactKind)
	build.Spec.BuildArgs = append(build.Spec.BuildArgs, core.VersionKey)
	build.Spec.Facts = append(build.Spec.Facts, fact)
	build.Spec.TagFormat = core.TagFormatVersion

	return nil
}
