package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/spiarh/gojo/pkg/core"
	"github.com/spiarh/gojo/pkg/provider"
	"github.com/spiarh/gojo/pkg/util"
)

func Facts() (*cobra.Command, error) {
	var command = &cobra.Command{
		Use:   "facts",
		Short: "Find or List the latest facts of a build image",
		Example: `gojo versions list --image-dir ~/image-git-dir --image nextcloud
gojo versions find --image-dir ~/image-git-dir --image nextcloud
`,
		SilenceUsage:      true,
		PersistentPreRunE: SetGlobalLogLevel,
	}

	if err := AddCommonPersistentFlags(listCommand); err != nil {
		return nil, err
	}
	if err := AddCommonPersistentFlags(getCommand); err != nil {
		return nil, err
	}

	command.AddCommand(listCommand)
	command.AddCommand(getCommand)

	return command, nil
}

var listCommand = &cobra.Command{
	Use:          core.ListAction,
	RunE:         func(cmd *cobra.Command, args []string) error { return facts(cmd, args) },
	SilenceUsage: true,
}

var getCommand = &cobra.Command{
	Use:          core.GetAction,
	RunE:         func(cmd *cobra.Command, args []string) error { return facts(cmd, args) },
	SilenceUsage: true,
}

func facts(command *cobra.Command, args []string) error {
	action := command.Use
	flagSet := command.Flags()
	opt, err := getOptions(flagSet)
	if err != nil {
		return err
	}

	log.Info().Bool(core.EnabledKey, opt.dryRun).Msg(core.DryRunFlag)

	build, err := core.NewBuildFromManifest(opt.buildFilePath)
	if err != nil {
		return err
	}

	if len(build.Spec.Sources) == 0 {
		log.Warn().Msg("no value sources defined, nothing to do")
		return nil
	}

	if err := build.ValidatePreProcess(); err != nil {
		log.Fatal().AnErr(core.ErrKey, err).Msg("")
	}

	if err = setFacts(flagSet, build.Spec.Facts, build.Spec.Sources); err != nil {
		log.Fatal().AnErr(core.ErrKey, err).Msg("retrieve facts")
	}

	if build.Image.Tag, err = core.BuildTag(build.Spec.Facts, build.Spec.TagFormat, build.Image.Context); err != nil {
		return err
	}

	if opt.dryRun || (action == core.ListAction) {
		return nil
	}

	return build.WriteToFile(build.Image.BuildfilePath)
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

				if fact.Value, err = repo.GetLatest(fact.Semver); err != nil {
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
