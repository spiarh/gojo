package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/blang/semver/v4"

	"github.com/spiarh/gojo/pkg/util"
)

func NewBuild(bmage, fromImage string) (*Build, error) {
	imageObj, err := NewImageFromFQIN(bmage)
	if err != nil {
		return nil, err
	}
	fromImageObj, err := NewImageFromFQIN(fromImage)
	if err != nil {
		return nil, err
	}

	return &Build{
		Image: &imageObj,
		Spec: &ImageSpec{
			FromImages: []FromImage{
				{Image: fromImageObj},
			},
		},
	}, nil
}

// NewBuildFromManifest returns a new decoded Build struct from a manifest.
func NewBuildFromManifest(manifestPath string) (*Build, error) {
	file, err := os.Open(manifestPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	configBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	build, err := DecodeBuild(configBytes)
	if err != nil {
		return nil, err
	}

	if err := build.ValidatePreProcess(); err != nil {
		return nil, err
	}

	build.Image.Context = path.Dir(manifestPath)
	build.Image.BuildfilePath = manifestPath

	return build, nil
}

func DecodeBuild(b []byte) (*Build, error) {
	build := &Build{}
	if err := yaml.Unmarshal(b, &build); err != nil {
		return nil, err
	}
	return build, nil
}

func EncodeBuild(b *Build) ([]byte, error) {
	return yaml.Marshal(b)
}

func (b *Build) String() string {
	data, _ := EncodeBuild(b)
	return string(data)
}

func (b *Build) WriteToFile(path string) error {
	data, err := EncodeBuild(b)
	if err != nil {
		return err
	}
	return util.WriteToFile(path, data, 0644)
}

func (b *Build) GetBuildArgs() map[string]string {
	buildArgs := make(map[string]string)

	for _, arg := range b.Spec.BuildArgs {
		for _, fact := range b.Spec.Facts {
			if arg == fact.Name {
				buildArgs[fact.Name] = fact.Value
				continue
			}
		}
	}
	fromImageArg := "FROM_IMAGE"
	for _, fromImage := range b.Spec.FromImages {
		if fromImage.Target == "" {
			buildArgs[fromImageArg] = fromImage.Image.String()
			continue
		}
		targetArg := "_" + strings.ToUpper(fromImage.Target)
		buildArgs[fromImageArg+targetArg] = fromImage.Image.String()
	}

	return buildArgs
}

func (b *Build) ValidatePreProcess() error {
	err := util.EnsureStringSliceDuplicates(b.Spec.BuildArgs)
	if err != nil {
		return err
	}

	if b.Image == nil {
		return fmt.Errorf("image definition missing")
	}

	if len(b.Spec.FromImages) == 0 {
		return fmt.Errorf("at least one fromImage must be defined")
	}

	numProviders := 0
	for _, source := range b.Spec.Sources {
		if source.Alpine != nil {
			numProviders++
		}
		if source.GitHub != nil {
			numProviders++
		}
	}
	if numProviders != len(b.Spec.Sources) {
		return fmt.Errorf("too many providers specified for one source")
	}

	for _, fact := range b.Spec.Facts {
		if fact.Name == "" {
			return fmt.Errorf("empty fact name")
		}

		if fact.Name == string(DateFactInternalName) || fact.Name == string(GitCommitFactInternalName) {
			return fmt.Errorf("internal fact name used: %s", fact.Name)
		}
		if fact.Source == "" {
			continue
		}
		found := false
		for _, source := range b.Spec.Sources {
			if fact.Source == source.Name {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("Source not found: %s", fact.Source)
		}
		if fact.Kind != VersionFactKind && fact.Semver != "" {
			return fmt.Errorf("SemVer specified for non version fact kind")
		}
		if fact.Kind == VersionFactKind && fact.Semver != "" {
			if _, err := semver.ParseRange(fact.Semver); err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *Build) Validate() error {
	if err := b.ValidatePreProcess(); err != nil {
		return err
	}
	for _, fact := range b.Spec.Facts {
		if fact.Value == "" {
			return fmt.Errorf("empty fact value")
		}
	}

	return nil
}
