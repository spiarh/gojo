package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/spiarh/gojo/pkg/core"
	"github.com/spiarh/gojo/pkg/manager"
)

const (
	imagesDirEnv = "GOJO_IMAGES_DIR"
)

// AddCommonPersistentFlags adds some common flags to a cobra command.
func AddCommonPersistentFlags(command *cobra.Command) error {
	command.PersistentFlags().Bool(core.DryRunFlag, false, "Do not write files nor execute commands")
	command.PersistentFlags().StringP(core.ContainerfileFlag, "c", core.ContainerfileName, "Name of the Containerfile")
	command.PersistentFlags().StringP(core.BuildfileFlag, "f", core.BuildFileName, "Name of the buildfile")
	command.PersistentFlags().StringP(core.ImageFlag, "i", "", "Name of the image as in the images directory name")
	command.PersistentFlags().StringP(core.ImagesDirFlag, "d", os.Getenv(imagesDirEnv), "Path to the container images directory")
	command.PersistentFlags().StringP(core.LogLevelFlag, "l", core.DefaultLogLevel, "Log level {debug,info,warn,error}")

	if err := command.MarkPersistentFlagRequired(core.ImageFlag); err != nil {
		return err
	}
	return nil
}

// AddCommonBuildFlags adds some common build flags to a cobra command.
func AddCommonBuildFlags(command *cobra.Command) {
	command.PersistentFlags().Bool(core.PushFlag, false, "Push the image after the build")
	command.PersistentFlags().Bool(core.TagLatestFlag, false, "Tag the built image as latest")
}

// AddBuildkitFlags adds some buildkit flags to a cobra command.
func AddBuildkitFlags(command *cobra.Command) {
	command.PersistentFlags().String(core.AddrFlag, "", "Buildkitd address, e.g podman-container://buildkitd")
	command.PersistentFlags().String(core.FrontendFlag, manager.DefaultBuildkitFrontend, "Frontend to convert any build definition to LLB")

	command.PersistentFlags().String(core.TLSServerNameFlag, "", "Buildkitd server name for certificate validation")
	command.PersistentFlags().String(core.TLSCaCertFlag, "", "CA certificate for validation")
	command.PersistentFlags().String(core.TLSCertFlag, "", "Client certificate")
	command.PersistentFlags().String(core.TLSKeyFlag, "", "Client key")
	command.PersistentFlags().String(core.TLSDirFlag, "", "Directory containing CA certificate, client certificate, and client key")
}

func SetGlobalLogLevel(cmd *cobra.Command, args []string) error {
	logLevel, err := cmd.Flags().GetString(core.LogLevelFlag)
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
	imagesDir, imageDir, imageName       string
	buildFileName, buildFilePath         string
	containerFileName, containerFilePath string
	dryRun                               bool
}

func getOptions(flagSet *pflag.FlagSet) (CommonOptions, error) {
	var opt CommonOptions
	var err error

	if opt.imageName, err = flagSet.GetString(core.ImageFlag); err != nil {
		return opt, err
	}
	if opt.imagesDir, err = flagSet.GetString(core.ImagesDirFlag); err != nil {
		return opt, err
	}
	if opt.imageDir, err = filepath.Abs(filepath.Dir(opt.imagesDir)); err != nil {
		opt.imageDir = path.Join(opt.imagesDir, opt.imageName)
	}
	opt.imageDir = path.Join(opt.imagesDir, opt.imageName)

	if opt.buildFileName, err = flagSet.GetString(core.BuildfileFlag); err != nil {
		return opt, err
	}
	opt.buildFilePath = path.Join(opt.imageDir, opt.buildFileName)

	if opt.containerFileName, err = flagSet.GetString(core.ContainerfileFlag); err != nil {
		return opt, err
	}
	opt.containerFilePath = path.Join(opt.imageDir, opt.containerFileName)

	if opt.dryRun, err = flagSet.GetBool(core.DryRunFlag); err != nil {
		return opt, err
	}

	return opt, nil
}
