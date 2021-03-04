package cmd

// Flags
const (
	dryRunFlag    = "dry-run"
	buildFileFlag = "build-file"
	imageFlag     = "image"
	imagesDirFlag = "images-dir"
	logLevelFlag  = "log-level"

	imageFQINFlag        = "image-fqin"
	fromImageFlag        = "from-image"
	fromImageBuilderFlag = "from-image-builder"

	ownerFlag     = "owner"
	repoFlag      = "repo"
	versionIDFlag = "version-id"
	pkgFlag       = "pkg"

	pushFlag      = "push"
	tagLatestFlag = "tag-latest"

	nameFlag  = "name"
	emailFlag = "email"
)
const (
	defaultLogLevel = "info"
)

// Filenames
const (
	defaultBuildFileName     = ".build.yaml"
	defaultContainerfileName = "Containerfile"
)

// Keys
const (
	contentKey        = "content"
	repoKey           = "repo"
	fileKey           = "file"
	nameKey           = "name"
	enabledKey        = "enabled"
	msgKey     string = "message"
	hashKey           = "hash"
	commitKey         = "commit"

	versionKey = "VERSION"
)

// Actions
const (
	listAction = "list"
	findAction = "find"
)

const (
	alpine = "alpine"
	simple = "simple"
	github = "github"
)
