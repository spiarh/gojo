package cmd

import (
	"github.com/lcavajani/gojo/pkg/buildconf"
	"github.com/lcavajani/gojo/pkg/manager"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Build() *cobra.Command {
	var command = &cobra.Command{
		Use:              "build",
		Short:            "Build a container image",
		Example:          "gojo build --push --tag-latest haproxy",
		SilenceUsage:     true,
		TraverseChildren: true,
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

	ibc, err := buildconf.NewImageFromFile(opt.ibcPath)
	if err != nil {
		return err
	}
	ibc.Metadata.Context = opt.imageDir
	ibc.Metadata.Path = opt.ibcPath

	// TODO: Add image validation
	mgrType := command.Use
	mgr, err := manager.New(flagSet, mgrType)
	if err != nil {
		return err
	}

	err = mgr.Build(ibc)
	if err != nil {
		return err
	}

	return nil
}
