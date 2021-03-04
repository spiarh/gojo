package cmd

import (
	"fmt"
	"path"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// AddCommonPersistentFlags adds some common flags to a cobra command.
func AddCommonPersistentFlags(command *cobra.Command) {
	command.PersistentFlags().Bool(dryRunFlag, false, "Do not write files nor execute commands")
	command.PersistentFlags().StringP(buildFileFlag, "f", defaultBuildFileName, "Name of the buildfile")
	command.PersistentFlags().StringP(imageFlag, "i", "", "Name of the image as in the images directory name")
	command.PersistentFlags().StringP(imagesDirFlag, "d", "", "Path to the container images directory")
	command.PersistentFlags().StringP(logLevelFlag, "l", defaultLogLevel, "Log level {debug,info,warn,error}")

	command.MarkPersistentFlagRequired(imageFlag)
	command.MarkPersistentFlagRequired(imagesDirFlag)
}

// AddCommonBuildFlags adds some common build flags to a cobra command.
func AddCommonBuildFlags(command *cobra.Command) {
	command.PersistentFlags().Bool(pushFlag, false, "Push the image after the build")
	command.PersistentFlags().Bool(tagLatestFlag, false, "Tag the built image as latest")
}

func SetGlobalLogLevel(cmd *cobra.Command, args []string) error {
	logLevel, err := cmd.Flags().GetString(logLevelFlag)
	if err != nil {
		return err
	}

	switch logLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warning":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		return fmt.Errorf("invalid log level: %s", logLevel)
	}
	return nil
}

type CommonOptions struct {
	imagesDir, imageDir, imageName string
	buildFileName, buildFilePath    string
	dryRun                         bool
}

func getOptions(flagSet *pflag.FlagSet) (CommonOptions, error) {
	var opt CommonOptions
	var err error

	if opt.imageName, err = flagSet.GetString(imageFlag); err != nil {
		return opt, err
	}
	if opt.imagesDir, err = flagSet.GetString(imagesDirFlag); err != nil {
		return opt, err
	}
	opt.imageDir = path.Join(opt.imagesDir, opt.imageName)

	if opt.buildFileName, err = flagSet.GetString(buildFileFlag); err != nil {
		return opt, err
	}
	opt.buildFilePath = path.Join(opt.imageDir, opt.buildFileName)

	if opt.dryRun, err = flagSet.GetBool(dryRunFlag); err != nil {
		return opt, err
	}

	return opt, nil
}
