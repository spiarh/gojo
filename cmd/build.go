package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/lcavajani/gojo/pkg/core"
	"github.com/lcavajani/gojo/pkg/manager"
)

func Build() (*cobra.Command, error) {
	var command = &cobra.Command{
		Use:               "build",
		Short:             "Build a container image",
		Example:           "gojo build --push --tag-latest haproxy",
		SilenceUsage:      true,
		TraverseChildren:  true,
		PersistentPreRunE: SetGlobalLogLevel,
	}

	// Buildah
	command.AddCommand(buildahCommand)
	if err := AddCommonPersistentFlags(buildahCommand); err != nil {
		return nil, err
	}
	AddCommonBuildFlags(buildahCommand)

	// Buildkit
	command.AddCommand(buildkitCommand)
	if err := AddCommonPersistentFlags(buildkitCommand); err != nil {
		return nil, err
	}
	AddCommonBuildFlags(buildkitCommand)
	AddBuildkitFlags(buildkitCommand)

	// Podman
	command.AddCommand(podmanCommand)
	if err := AddCommonPersistentFlags(podmanCommand); err != nil {
		return nil, err
	}
	AddCommonBuildFlags(podmanCommand)

	return command, nil
}

var buildkitCommand = &cobra.Command{
	Use:          string(manager.BuildkitType),
	RunE:         func(cmd *cobra.Command, args []string) error { return build(cmd, args) },
	SilenceUsage: true,
}

var buildahCommand = &cobra.Command{
	Use:          string(manager.BuildahType),
	RunE:         func(cmd *cobra.Command, args []string) error { return build(cmd, args) },
	SilenceUsage: true,
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
	log.Info().Bool(core.EnabledKey, opt.dryRun).Msg(core.DryRunFlag)

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
