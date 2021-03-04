package cmd

import (
	"github.com/lcavajani/gojo/pkg/core"
	"github.com/lcavajani/gojo/pkg/provider"
	"github.com/lcavajani/gojo/pkg/util"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Versions() *cobra.Command {
	var command = &cobra.Command{
		Use:   "versions",
		Short: "Find or List the latest versions of a package",
		Example: `gojo versions list --image-dir ~/image-git-dir --image nextcloud
gojo versions find --image-dir ~/image-git-dir --image nextcloud
`,
		SilenceUsage:      true,
		PersistentPreRunE: SetGlobalLogLevel,
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

	build, err := core.NewBuildFromManifest(opt.buildFilePath)
	if err != nil {
		return err
	}

	if len(build.Spec.Sources) == 0 {
		log.Warn().Msg("no value sources defined, nothing to do")
		return nil
	}

	if err := build.ValidatePreProcess(); err != nil {
		log.Fatal().AnErr("err", err).Msg("build manifest validation")
	}

	if err = setFacts(flagSet, build.Spec.Facts, build.Spec.Sources); err != nil {
		log.Fatal().AnErr("err", err).Msg("retrieve facts")
	}

	if build.Metadata.Tag, err = core.BuildTag(build.Spec.Facts, build.Spec.TagFormat, build.Metadata.Context); err != nil {
		return err
	}

	if opt.dryRun || (action == listAction) {
		return nil
	}

	return build.WriteToFile(build.Metadata.Path)
}

func setFacts(flagSet *pflag.FlagSet, facts []*core.Fact, sources []core.Source) error {
	for _, fact := range facts {
		if fact.Source == "" {
			continue
		}
		for _, src := range sources {
			if fact.Source == src.Name {
				repo, err := provider.New(flagSet, src)
				if err != nil {
					return err
				}

				if fact.Value, err = repo.GetLatest(); err != nil {
					return err
				}

				if fact.Kind == core.VersionFactKind {
					fact.Value = util.SanitizeVersion(fact.Value)
				}
			}
		}
		if fact.Value == "" {
			log.Fatal().Msgf("no value found for fact with name: %s", fact.Name)
		}
	}
	return nil
}
