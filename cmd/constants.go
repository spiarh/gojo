package cmd

// Flags
const (
	dryRunFlag    = "dry-run"
	imageFlag     = "image"
	imagesDirFlag = "images-dir"

	imageFQINFlag        = "image-fqin"
	fromImageFlag        = "from-image"
	fromImageBuilderFlag = "from-image-builder"

	ownerFlag     = "owner"
	repoFlag      = "repo"
	versionIDFlag = "version-id"
	pkgFlag       = "pkg"

	pushFlag      = "push"
	tagLatestFlag = "tag-latest"
)

// Filenames
const (
	gojoFilename      = ".gojo.yaml"
	containerfileName = "Containerfile"
)

// Keys
const (
	contentKey = "content"
	nameKey    = "name"
	enabledKey = "enabled"
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
