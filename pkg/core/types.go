package core

type Build struct {
	Image *Image     `yaml:"image"`
	Spec  *ImageSpec `yaml:"spec"`
}

type Image struct {
	// Registry is the registry with the path of the image.
	Registry string `yaml:"registry"`
	// Name is the name of the image.
	Name string `yaml:"name"`
	// Tag is the tag of the image.
	Tag string `yaml:"tag"`

	// Containerfile is the name of the Containerfile.
	Containerfile string `yaml:"-"`
	// Context is the context for the build.
	Context string `yaml:"-"`
	// BuildfilePath is the path of the build file.
	BuildfilePath string `yaml:"-"`
}

type ImageSpec struct {
	FromImages []FromImage `yaml:"fromImages"`
	BuildArgs  BuildArgs   `yaml:"buildArgs"`
	TagFormat  string      `yaml:"tagFormat,omitempty"`
	Facts      []*Fact     `yaml:"facts,omitempty"`
	Sources    []Source    `yaml:"sources,omitempty"`
}

type FromImage struct {
	Image  `yaml:",inline"`
	Target string `yaml:"target,omitempty"`
}

type BuildArgs []string

type Fact struct {
	Name   string   `yaml:"name"`
	Value  string   `yaml:"value"`
	Source string   `yaml:"source,omitempty"`
	Kind   FactKind `yaml:"kind"`
	Semver string   `yaml:"semver,omitempty"`
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

type Source struct {
	Name     string `yaml:"name"`
	Provider `yaml:",inline"`
}

type Provider struct {
	Alpine *AlpineSource `yaml:"alpine,omitempty"`
	GitHub *GitHubSource `yaml:"github,omitempty"`
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
