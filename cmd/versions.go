package cmd

import (
	"github.com/lcavajani/gojo/pkg/provider"
	"github.com/lcavajani/gojo/pkg/util"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Versions() *cobra.Command {
	var command = &cobra.Command{
		Use:   "versions",
		Short: "Find or List the latest versions of a package",
		Example: `gojo versions list --image-dir ~/image-git-dir --image nextcloud
gojo versions find --image-dir ~/image-git-dir --image nextcloud
`,
		SilenceUsage: true,
	}

	AddCommonPersistentFlags(listCommand)
	AddCommonPersistentFlags(findCommand)

	command.AddCommand(listCommand)
	command.AddCommand(findCommand)
	return command
}

var listCommand = &cobra.Command{
	Use:          listAction,
	RunE:         func(cmd *cobra.Command, args []string) error { return versions(cmd, args) },
	SilenceUsage: true,
}

var findCommand = &cobra.Command{
	Use:          findAction,
	RunE:         func(cmd *cobra.Command, args []string) error { return versions(cmd, args) },
	SilenceUsage: true,
}

func versions(command *cobra.Command, args []string) error {
	action := command.Use
	flagSet := command.Flags()
	opt, err := getOptions(flagSet)
	if err != nil {
		return err
	}

	log.Info().Bool(enabledKey, opt.dryRun).Msg(dryRunFlag)

	ibc, err := decodeImageFromFile(opt)
	if err != nil {
		return err
	}

	if ibc.Spec.TagBuild == nil {
		log.Warn().Msg("no tag build defined, nothing to do")
		return nil
	}

	repo, err := provider.New(flagSet, ibc)
	if err != nil {
		return err
	}
	v, err := repo.GetLatest()
	if err != nil {
		return err
	}

	version := util.SanitizeVersion(v)

	// TODO: Tag shoud have different strategies
	ibc.Metadata.Tag = version

	if ibc.Spec.TagBuild.Version != version {
		log.Info().Str("arg", "VERSION").
			Str("current", ibc.Spec.TagBuild.Version).
			Str("new", version).
			Msg("set new version")
		ibc.Spec.TagBuild.Version = version
	}

	if opt.dryRun || (action == listAction) {
		return nil
	}

	return ibc.WriteToFile(ibc.Metadata.Path)
}
