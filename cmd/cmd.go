package cmd

import (
	"path"

	"github.com/lcavajani/gojo/pkg/buildconf"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// AddCommonPersistentFlags adds some common flags to a cobra command.
func AddCommonPersistentFlags(command *cobra.Command) {
	command.PersistentFlags().StringP(imageFlag, "i", "", "Name of the image as in the images directory name")
	command.PersistentFlags().StringP(imagesDirFlag, "d", "", "Path to the container images directory")
	command.PersistentFlags().Bool(dryRunFlag, false, "Do not write files nor execute commands")

	command.MarkPersistentFlagRequired(imageFlag)
	command.MarkPersistentFlagRequired(imagesDirFlag)
}

// AddCommonBuildFlags adds some common build flags to a cobra command.
func AddCommonBuildFlags(command *cobra.Command) {
	command.PersistentFlags().Bool(pushFlag, false, "Push the image after the build")
	command.PersistentFlags().Bool(tagLatestFlag, false, "Tag the built image as latest")
}

type CommonOptions struct {
	imagesDir, imageDir, imageName, ibcPath string
	dryRun                                  bool
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
	opt.ibcPath = path.Join(opt.imagesDir, opt.imageName, gojoFilename)

	if opt.dryRun, err = flagSet.GetBool(dryRunFlag); err != nil {
		return opt, err
	}

	return opt, nil
}

func decodeImageFromFile(opt CommonOptions) (*buildconf.Image, error) {
	ibc, err := buildconf.NewImageFromFile(opt.ibcPath)
	if err != nil {
		return nil, err
	}
	ibc.Metadata.Context = opt.imageDir
	ibc.Metadata.Path = opt.ibcPath

	return ibc, nil
}
