package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/lcavajani/gojo/pkg/util"
	"gopkg.in/yaml.v2"
)

type Image struct {
	// Registry is the registry with the path of the image.
	Registry string `yaml:"registry"`
	// Name is the name of the image.
	Name string `yaml:"name"`
	// Tag is the tag of the image.
	Tag string `yaml:"tag"`

	// Context is the context for the build.
	Context string `yaml:"-"`
	// Path is the path of the build file.
	Path string `yaml:"-"`
}

func NewImageFromFQIN(fqin string) (Image, error) {
	var image = Image{}
	registry, name, tag, err := util.ParseImageFullName(fqin)
	if err != nil {
		return image, err
	}
	image = NewImage(registry, name, tag)

	return image, nil
}

func NewImage(registry, name, tag string) Image {
	return Image{
		Registry: registry,
		Name:     name,
		Tag:      tag,
	}
}

func (b *Image) String() string {
	return fmt.Sprintf("%s/%s:%s", b.Registry, b.Name, b.Tag)
}

func (b *Image) StringWithTag(tag string) string {
	return fmt.Sprintf("%s/%s:%s", b.Registry, b.Name, tag)
}

type AlpineSource struct {
	Package    string `yaml:"package"`
	Repository string `yaml:"repository"`
	VersionId  string `yaml:"versionId"`
	Arch       string `yaml:"arch,omitempty"`
	Mirror     string `yaml:"mirror,omitempty"`
}

type GitHubSource struct {
	Owner      string `yaml:"owner"`
	Repository string `yaml:"repository"`
	// TODO: Rename, Kind?
	Object GitHubObject `yaml:"object"`
}

type BuildArgs []string

type Fact struct {
	Name   string `yaml:"name"`
	Value  string `yaml:"value"`
	Source string `yaml:"source,omitempty"`
	// TODO: Set default, can't be empty ?
	Kind FactKind `yaml:"kind"`
	// Semver  string  `yaml:"semver,omitempty"`
}

type FactKind string

const (
	VersionFactKind FactKind = "version"
	StringFactKind  FactKind = "string"
)

type FactInternalName string

const (
	DateFactInternalName      FactInternalName = "date"
	GitCommitFactInternalName FactInternalName = "gitCommit"
)

type ImageSpec struct {
	FromImage        *Image    `yaml:"fromImage"`
	FromImageBuilder *Image    `yaml:"fromImageBuilder,omitempty"`
	BuildArgs        BuildArgs `yaml:"buildArgs"`
	TagFormat        string    `yaml:"tagFormat,omitempty"`
	Facts            []*Fact   `yaml:"facts,omitempty"`
	Sources          []Source  `yaml:"sources,omitempty"`
}

type Source struct {
	Name string `yaml:"name"`
	// TODO: Make sure only one is defined
	Provider `yaml:",inline"`
}

type Provider struct {
	Alpine *AlpineSource `yaml:"alpine,omitempty"`
	GitHub *GitHubSource `yaml:"github,omitempty"`
}

type Build struct {
	Metadata *Image     `yaml:"metadata"`
	Spec     *ImageSpec `yaml:"spec"`
}

func (b *Build) GetBuildArgs() map[string]string {
	buildArgs := map[string]string{}

	for _, arg := range b.Spec.BuildArgs {
		for _, fact := range b.Spec.Facts {
			if arg == fact.Name {
				buildArgs[fact.Name] = fact.Value
				continue
			}
		}

	}

	if b.Spec.FromImage != nil {
		buildArgs["FROM_IMAGE"] = b.Spec.FromImage.String()
	}
	if b.Spec.FromImageBuilder != nil {
		buildArgs["FROM_IMAGE_BUILDER"] = b.Spec.FromImageBuilder.String()
	}

	return buildArgs
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

func NewBuild(bmage, fromImage string) (*Build, error) {
	imageMetaObj, err := NewImageFromFQIN(bmage)
	if err != nil {
		return nil, err
	}
	fromImageObj, err := NewImageFromFQIN(fromImage)
	if err != nil {
		return nil, err
	}

	return &Build{
		Metadata: &imageMetaObj,
		Spec: &ImageSpec{
			FromImage: &fromImageObj,
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

	build.Metadata.Context = path.Dir(manifestPath)
	build.Metadata.Path = manifestPath

	return build, nil
}

func DecodeBuild(bn []byte) (*Build, error) {
	build := &Build{}
	if err := yaml.Unmarshal(bn, &build); err != nil {
		return nil, err
	}
	return build, nil
}

func EncodeBuild(b *Build) ([]byte, error) {
	return yaml.Marshal(b)
}

func (b *Build) ValidatePreProcess() error {
	// TODO: use field validation kubernetes
	// TODO: make is more explicit, which field is invalid
	err := util.EnsureStringSliceDuplicates(b.Spec.BuildArgs)
	if err != nil {
		return err
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

func NewAlpineSource(name, pkg, repo, versionId string) Source {
	return Source{
		Name: name,
		Provider: Provider{
			Alpine: &AlpineSource{
				Package:    pkg,
				Repository: repo,
				VersionId:  versionId,
			},
		},
	}
}

func NewGitHubSource(name, repo, owner string) Source {
	return Source{
		Name: name,
		Provider: Provider{
			GitHub: &GitHubSource{
				Owner:      owner,
				Repository: repo,
				Object:     GitHubObjectRelease,
			},
		},
	}

}

func NewFact(name, value, source string, kind FactKind) *Fact {
	return &Fact{
		Name:   name,
		Value:  value,
		Kind:   kind,
		Source: source,
	}

}

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
// func ValidateConfigPath(path string) error {
//     s, err := os.Stat(path)
//     if err != nil {
//         return err
//     }
//     if s.IsDir() {
//         return fmt.Errorf("'%s' is a directory, not a normal file", path)
//     }
//     return nil
// }