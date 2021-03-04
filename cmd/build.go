package cmd

import (
	"github.com/lcavajani/gojo/pkg/core"
	"github.com/lcavajani/gojo/pkg/manager"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Build() *cobra.Command {
	var command = &cobra.Command{
		Use:               "build",
		Short:             "Build a container image",
		Example:           "gojo build --push --tag-latest haproxy",
		SilenceUsage:      true,
		TraverseChildren:  true,
		PersistentPreRunE: SetGlobalLogLevel,
	}

	command.AddCommand(podmanCommand)
	AddCommonPersistentFlags(podmanCommand)
	AddCommonBuildFlags(podmanCommand)

	return command
}

var podmanCommand = &cobra.Command{
	Use:          string(manager.PodmanType),
	RunE:         func(cmd *cobra.Command, args []string) error { return build(cmd, args) },
	SilenceUsage: true,
}

func build(command *cobra.Command, args []string) error {
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

	if err := build.Validate(); err != nil {
		log.Fatal().AnErr("err", err).Msg("build manifest validation")
	}

	mgrType := command.Use
	mgr, err := manager.New(flagSet, mgrType)
	if err != nil {
		return err
	}

	err = mgr.Build(build)
	if err != nil {
		return err
	}

	return nil
}
