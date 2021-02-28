package provider

import (
	"fmt"

	"github.com/lcavajani/gojo/pkg/buildconf"
	"github.com/spf13/pflag"
)

const (
	defaultArch = "x86_64"
)

type provider interface {
	GetLatest() (string, error)
}

var _ provider = &Alpine{}
var _ provider = &GitHub{}

func New(pflagSet *pflag.FlagSet, conf *buildconf.Image) (provider, error) {
	var source buildconf.Source
	for _, src := range conf.Spec.Sources {
		if src.Name == conf.Spec.TagBuild.Source {
			source = src
			break
		}
	}

	if (source == buildconf.Source{}) {
		return nil, fmt.Errorf("no source found for name: %s", conf.Spec.TagBuild.Source)
	}

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

	return nil, fmt.Errorf("provider type not recognized: %s", conf.Spec.TagBuild.Source)
}

func setDefaultsAlpine(repo *buildconf.AlpineSource) {
	if repo.Mirror == "" {
		repo.Mirror = alpineDefaultMirror
	}
	if repo.Arch == "" {
		repo.Arch = defaultArch
	}
}
