package manager

type managerType string

const (
	BuildahType  managerType = "buildah"
	BuildkitType managerType = "buildkit"
	KanikoType   managerType = "kaniko"
	PodmanType   managerType = "podman"
)

const (
	DefaultBuildkitFrontend = "dockerfile.v0"
)
