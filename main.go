package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/lcavajani/gojo/cmd"
	"github.com/lcavajani/gojo/pkg/util"
)

func main() {
	// Use pretty output if we are in a terminal, json otherwise.
	if util.IsTTYAllocated() {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	var rootCmd = &cobra.Command{
		Use: "gojo",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		SilenceUsage: true,
	}

	cmdBuild := cmd.Build()
	cmdCommit := cmd.Commit()
	cmdNew := cmd.New()
	cmdVersions := cmd.Versions()

	rootCmd.AddCommand(cmdBuild)
	rootCmd.AddCommand(cmdCommit)
	rootCmd.AddCommand(cmdNew)
	rootCmd.AddCommand(cmdVersions)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
