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
			if err := cmd.Help(); err != nil {
				log.Fatal().AnErr("err", err).Msg("")
			}
		},
		SilenceUsage: true,
	}

	var cmdBuild, cmdCommit, cmdFacts, cmdScaffold, cmdVersion *cobra.Command
	var err error

	if cmdBuild, err = cmd.Build(); err != nil {
		log.Fatal().AnErr("err", err).Msg("")
	}
	if cmdCommit, err = cmd.Commit(); err != nil {
		log.Fatal().AnErr("err", err).Msg("")
	}
	if cmdFacts, err = cmd.Facts(); err != nil {
		log.Fatal().AnErr("err", err).Msg("")
	}
	if cmdScaffold, err = cmd.Scaffold(); err != nil {
		log.Fatal().AnErr("err", err).Msg("")
	}
	cmdVersion = cmd.Version()

	rootCmd.AddCommand(cmdBuild)
	rootCmd.AddCommand(cmdCommit)
	rootCmd.AddCommand(cmdFacts)
	rootCmd.AddCommand(cmdScaffold)
	rootCmd.AddCommand(cmdVersion)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
