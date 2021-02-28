package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/lcavajani/gojo/cmd"
)

const (
	logLevelFlag    = "log-level"
	defaultLogLevel = "info"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	var logLevel string
	var rootCmd = &cobra.Command{
		Use: "gojo",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := setGlobalLogLevel(logLevel); err != nil {
				return err
			}
			return nil
		},
		SilenceUsage:     true,
		TraverseChildren: true,
	}
	rootCmd.Flags().StringVarP(&logLevel, logLevelFlag, "l", defaultLogLevel, "Log level {debug,info,warn,error}")

	cmdBuild := cmd.Build()
	cmdNew := cmd.New()
	cmdVersions := cmd.Versions()

	rootCmd.AddCommand(cmdBuild)
	rootCmd.AddCommand(cmdNew)
	rootCmd.AddCommand(cmdVersions)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func setGlobalLogLevel(level string) error {
	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warning":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		return fmt.Errorf("invalid log level: %s", level)
	}
	return nil
}
