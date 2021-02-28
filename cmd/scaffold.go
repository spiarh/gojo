package cmd

import (
	"path"

	"github.com/lcavajani/gojo/pkg/buildconf"
	"github.com/lcavajani/gojo/pkg/util"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func New() *cobra.Command {
	var command = &cobra.Command{
		Use:          "scaffold",
		Short:        "Scaffold a new image project",
		Example:      "  gojo new github / gojo new alpine",
		SilenceUsage: true,
	}

	AddCommonPersistentFlags(command)

	command.PersistentFlags().String(imageFQINFlag, "", "Fully qualified image name (FQIN)")
	command.PersistentFlags().String(fromImageFlag, "", "Image to use for the base image")
	command.PersistentFlags().String(fromImageBuilderFlag, "", "Image to use for the builder")

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

func getScaffoldFlags(flagSet *pflag.FlagSet) (string, string, string, error) {
	var imageFQIN, fromImage, fromImageBuilder string
	var err error

	if imageFQIN, err = flagSet.GetString(imageFQINFlag); err != nil {
		return imageFQIN, fromImage, fromImageBuilder, err
	}
	if fromImage, err = flagSet.GetString(fromImageFlag); err != nil {
		return imageFQIN, fromImage, fromImageBuilder, err
	}
	if fromImageBuilder, err = flagSet.GetString(fromImageBuilderFlag); err != nil {
		return imageFQIN, fromImage, fromImageBuilder, err
	}
	return imageFQIN, fromImage, fromImageBuilder, nil
}

func scaffoldProject(command *cobra.Command, args []string) error {
	var err error
	providerType := command.Use

	flagSet := command.Flags()
	opt, err := getOptions(flagSet)
	if err != nil {
		return err
	}

	imageFQIN, fromImage, fromImageBuilder, err := getScaffoldFlags(flagSet)
	if err != nil {
		return err
	}

	log.Info().Bool(enabledKey, opt.dryRun).Msg(dryRunFlag)

	// Objects
	ibc, err := buildconf.NewImage(imageFQIN, fromImage)
	if err != nil {
		return err
	}

	if fromImageBuilder != "" {
		fromImageBuilderMeta, err := buildconf.NewImageMetaFromFullName(fromImageBuilder)
		if err != nil {
			return err
		}
		ibc.Spec.FromImageBuilder = &fromImageBuilderMeta
	}

	switch providerType {
	case alpine:
		err = addAlpineProvider(ibc, flagSet)
		if err != nil {
			return err
		}
	case github:
		err = addGitHubProvider(ibc, flagSet)
		if err != nil {
			return err
		}
	}

	ibcBytes, err := buildconf.Encode(ibc)
	if err != nil {
		return err
	}

	containerfile, err := buildconf.TemplateContainerfile(ibc)
	if err != nil {
		return err
	}

	log.Debug().Str(nameKey, gojoFilename).
		Bytes(contentKey, ibcBytes).
		Msg("gojo image build file creation")
	log.Debug().Str(nameKey, containerfileName).
		Bytes(contentKey, containerfile).
		Msg("Containerfile creation")

	if opt.dryRun {
		return nil
	}

	err = util.MakeDir(opt.imageDir, 0755)
	if err != nil {
		return err
	}

	ibc.WriteToFile(opt.ibcPath)
	containerfile.WriteToFile(path.Join(opt.imageDir, containerfileName))

	return nil
}

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

func addAlpineProvider(image *buildconf.Image, flagSet *pflag.FlagSet) error {
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

	tagBuild := buildconf.NewTagBuild(buildconf.Version, buildconf.AlpineProviderName)
	image.Spec.TagBuild = tagBuild

	source := buildconf.NewAlpineSource(buildconf.AlpineProviderName, pkg, repo, versionId)
	image.Spec.Sources = []buildconf.Source{source}

	return nil
}

// GitHub
var newGitHubCommand = &cobra.Command{
	Use:          github,
	RunE:         func(cmd *cobra.Command, args []string) error { return scaffoldProject(cmd, args) },
	SilenceUsage: true,
}

func addGitHubProvider(image *buildconf.Image, flagSet *pflag.FlagSet) error {
	var repo, owner string
	var err error

	if repo, err = flagSet.GetString(repoFlag); err != nil {
		return err
	}
	if owner, err = flagSet.GetString(ownerFlag); err != nil {
		return err
	}

	tagBuild := buildconf.NewTagBuild(buildconf.Version, buildconf.GitHubProviderName)
	image.Spec.TagBuild = tagBuild

	source := buildconf.NewGitHubSource(buildconf.GitHubProviderName, repo, owner)
	image.Spec.Sources = []buildconf.Source{source}

	return nil
}
