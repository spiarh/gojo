package manager

import (
	"fmt"

	"github.com/lcavajani/gojo/pkg/core"
	"github.com/spf13/pflag"
)

type managerType string

const (
	BuildahType  managerType = "buildah"
	BuildkitType managerType = "buildkit"
	KanikoType   managerType = "kaniko"
	PodmanType   managerType = "podman"
)

type manager interface {
	Build(ibc *core.Build) error
}

var _ manager = &Podman{}

func New(flagSet *pflag.FlagSet, mgrType string) (manager, error) {
	switch mgrType {
	case string(BuildahType):
		fmt.Println("NOT IMPLEMENTED")
	case string(BuildkitType):
		fmt.Println("NOT IMPLEMENTED")
	case string(KanikoType):
		fmt.Println("NOT IMPLEMENTED")
	case string(PodmanType):
		push, err := flagSet.GetBool("push")
		if err != nil {
			return nil, err
		}
		tagLatest, err := flagSet.GetBool("tag-latest")
		if err != nil {
			return nil, err
		}
		dryRun, err := flagSet.GetBool("dry-run")
		if err != nil {
			return nil, err
		}

		p, err := NewPodman(push, tagLatest, dryRun)
		if err != nil {
			return nil, err
		}
		return p, nil
	}

	return nil, fmt.Errorf("Manager type not recognized: %s", mgrType)
}
