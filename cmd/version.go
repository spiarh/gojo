package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/spiarh/gojo/pkg/version"
)

func Version() *cobra.Command {
	var command = &cobra.Command{
		Use:          "version",
		Short:        "Display the version information",
		Example:      "gojo version",
		RunE:         func(cmd *cobra.Command, args []string) error { return getVersion(cmd, args) },
		SilenceUsage: true,
	}

	return command
}

func getVersion(command *cobra.Command, args []string) error {
	info := version.Get()

	v, err := info.Format()
	if err != nil {
		return err
	}

	fmt.Println(v)
	return nil
}
