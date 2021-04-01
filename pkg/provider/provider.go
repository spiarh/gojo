package provider

import (
	"fmt"

	"github.com/spf13/pflag"

	"github.com/spiarh/gojo/pkg/core"
)

const (
	defaultArch = "x86_64"
)

type provider interface {
	GetLatest(semverRange string) (string, error)
}

var _ provider = &Alpine{}
var _ provider = &GitHub{}

func New(pflagSet *pflag.FlagSet, source core.Source) (provider, error) {
	switch {
	case source.Provider.Alpine != nil:
		a := source.Provider.Alpine
		setDefaultsAlpine(a)
		prvdr := NewAlpine(a.Mirror, a.Arch, a.VersionId, a.Repository, a.Package)
		return prvdr, nil
	case source.Provider.GitHub != nil:
		g := source.Provider.GitHub
		prvdr := NewGitHub(g.Owner, g.Repository, g.Object)
		return prvdr, nil
	}

	return nil, fmt.Errorf("provider type not recognized: %s", source.Name)
}

func setDefaultsAlpine(repo *core.AlpineSource) {
	if repo.Mirror == "" {
		repo.Mirror = alpineDefaultMirror
	}
	if repo.Arch == "" {
		repo.Arch = defaultArch
	}
}
