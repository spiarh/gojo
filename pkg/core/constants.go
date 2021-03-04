package core

type GitHubObject string

const (
	GitHubObjectRelease GitHubObject = "release"
	GitHubObjectTag     GitHubObject = "tag"
)

const (
	TagFormatVersion = "{{ .VERSION }}"
)

const ContainerfileName = "Containerfile"
