package buildconf

const (
	AlpineProviderName = "alpine"
	GitHubProviderName = "github"
)

type GitHubObject string

const (
	GitHubObjectRelease GitHubObject = "release"
	GitHubObjectTag     GitHubObject = "tag"
)

type TagType string

const (
	Tag         TagType = "TAG"
	Version     TagType = "VERSION"
	VersionDate TagType = "VERSION_DATE"
	VersionGit  TagType = "VERSION_GIT"
)
