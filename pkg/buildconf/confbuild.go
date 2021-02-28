package buildconf

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/lcavajani/gojo/pkg/util"
	"gopkg.in/yaml.v2"
)

type ImageMeta struct {
	// Registry is the registry with the path of the image.
	Registry string `yaml:"registry"`
	// Name is the name of the image.
	Name string `yaml:"name"`
	// Tag is the tag of the image.
	Tag string `yaml:"tag"`

	// Context is the context for the build.
	Context string `yaml:"-"`
	// Path is the path of the image build file.
	Path string `yaml:"-"`
}

func NewImageMetaFromFullName(image string) (ImageMeta, error) {
	var ibc = ImageMeta{}
	registry, name, tag, err := util.ParseImageFullName(image)
	if err != nil {
		return ibc, err
	}
	ibc = NewImageMeta(registry, name, tag)

	return ibc, nil
}

func NewImageMeta(registry, name, tag string) ImageMeta {
	return ImageMeta{
		Registry: registry,
		Name:     name,
		Tag:      tag,
	}
}

func (i *ImageMeta) GetFullName() string {
	return fmt.Sprintf("%s/%s:%s", i.Registry, i.Name, i.Tag)
}

func (i *ImageMeta) GetFullNameWithTag(tag string) string {
	return fmt.Sprintf("%s/%s:%s", i.Registry, i.Name, tag)
}

type AlpineSource struct {
	Package    string `yaml:"package"`
	Repository string `yaml:"repository"`
	VersionId  string `yaml:"versionId"`
	Arch       string `yaml:"arch,omitempty"`
	Mirror     string `yaml:"mirror,omitempty"`
}

type GitHubSource struct {
	Owner      string       `yaml:"owner"`
	Repository string       `yaml:"repository"`
	Object     GitHubObject `yaml:"object"`
}

type TagBuild struct {
	Semver  string  `yaml:"semver,omitempty"`
	Type    TagType `yaml:"type"`
	Version string  `yaml:"version"`
	Source  string  `yaml:"source,omitempty"`
}

type ImageSpec struct {
	FromImage        *ImageMeta `yaml:"fromImage"`
	FromImageBuilder *ImageMeta `yaml:"fromImageBuilder,omitempty"`
	TagBuild         *TagBuild  `yaml:"tagBuild,omitempty"`
	Sources          []Source   `yaml:"sources,omitempty"`
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

type Image struct {
	Metadata *ImageMeta `yaml:"metadata"`
	Spec     *ImageSpec `yaml:"spec"`
}

func (i *Image) GetBuildArgs() map[string]string {
	buildArgs := map[string]string{}

	//TODO: add build args
	if i.Spec.TagBuild.Version != "" {
		buildArgs["VERSION"] = i.Spec.TagBuild.Version
	}

	if i.Spec.FromImage != nil {
		buildArgs["FROM_IMAGE"] = i.Spec.FromImage.GetFullName()
	}
	if i.Spec.FromImageBuilder != nil {
		buildArgs["FROM_IMAGE_BUILDER"] = i.Spec.FromImageBuilder.GetFullName()
	}

	return buildArgs
}

func (i *Image) WriteToFile(path string) error {
	data, err := Encode(i)
	if err != nil {
		return err
	}

	return util.WriteToFile(path, data, 0644)
}

func NewImage(image, fromImage string) (*Image, error) {
	imageMeta, err := NewImageMetaFromFullName(image)
	if err != nil {
		return nil, err
	}
	fromImageMeta, err := NewImageMetaFromFullName(fromImage)
	if err != nil {
		return nil, err
	}

	return &Image{
		Metadata: &imageMeta,
		Spec: &ImageSpec{
			FromImage: &fromImageMeta,
		},
	}, nil
}

// NewImageFromFile returns a new decoded Image struct from a file.
func NewImageFromFile(path string) (*Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	configBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return Decode(configBytes)
}

func Decode(in []byte) (*Image, error) {
	ibc := &Image{}
	if err := yaml.Unmarshal(in, &ibc); err != nil {
		return nil, err
	}
	return ibc, nil
}

func Encode(i *Image) ([]byte, error) {
	return yaml.Marshal(i)
}

// TODO: update when adding tagType date, git etc
func NewTagBuild(tagType TagType, source string) *TagBuild {
	return &TagBuild{
		Type: tagType,
		// TODO: add type
		Source: source,
	}
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
