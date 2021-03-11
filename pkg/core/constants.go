package core

type GitHubObject string

const (
	GitHubObjectRelease GitHubObject = "release"
	GitHubObjectTag     GitHubObject = "tag"
)

const (
	TagFormatVersion = "{{ .VERSION }}"
)

// Flags
const (
	DryRunFlag    = "dry-run"
	BuildFileFlag = "build-file"
	ImageFlag     = "image"
	ImagesDirFlag = "images-dir"
	LogLevelFlag  = "log-level"

	ImageFQINFlag        = "image-fqin"
	FromImageFlag        = "from-image"
	FromImageBuilderFlag = "from-image-builder"

	OwnerFlag     = "owner"
	RepoFlag      = "repo"
	VersionIDFlag = "version-id"
	PkgFlag       = "pkg"

	AddrFlag          = "addr"
	FrontendFlag      = "frontend"
	TLSServerNameFlag = "tls-server-name"
	TLSCaCertFlag     = "tls-ca-cert"
	TLSCertFlag       = "tls-cert"
	TLSKeyFlag        = "tls-key"
	TLSDirFlag        = "tls-dir"

	PushFlag      = "push"
	TagLatestFlag = "tag-latest"

	NameFlag  = "name"
	EmailFlag = "email"
)
const (
	DefaultLogLevel = "info"
)

// Filenames
const (
	BuildFileName     = ".build.yaml"
	ContainerfileName = "Containerfile"
)

// Keys
const (
	// ContentKey        = "content"
	RepoKey    = "repo"
	FileKey    = "file"
	NameKey    = "name"
	EnabledKey = "enabled"
	MsgKey     = "message"
	HashKey    = "hash"
	CommitKey  = "commit"
	VersionKey = "VERSION"
)

// Actions
const (
	ListAction = "list"
	GetAction  = "get"
)

const (
	Alpine = "alpine"
	Simple = "simple"
	Github = "github"
)
