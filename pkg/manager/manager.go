package manager

import (
	"fmt"

	"github.com/spf13/pflag"

	"github.com/lcavajani/gojo/pkg/core"
	"github.com/lcavajani/gojo/pkg/execute"
)

type manager interface {
	Build(ibc *core.Build) error
}

var _ manager = &Podman{}

func New(flagSet *pflag.FlagSet, mgrType string) (manager, error) {
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

	switch mgrType {
	case string(BuildahType):
		b, err := NewBuildah(push, tagLatest, dryRun)
		if err != nil {
			return nil, err
		}
		return b, nil
	case string(BuildkitType):
		opt, err := getBuildkitOptions(flagSet)
		if err != nil {
			return nil, err
		}
		b, err := NewBuildkit(push, tagLatest, dryRun, opt)
		if err != nil {
			return nil, err
		}
		return b, nil
	case string(PodmanType):
		p, err := NewPodman(push, tagLatest, dryRun)
		if err != nil {
			return nil, err
		}
		return p, nil
	case string(KanikoType):
		k, err := NewKaniko(push, tagLatest, dryRun)
		if err != nil {
			return nil, err
		}
		return k, nil
	}

	return nil, fmt.Errorf("Manager type not recognized: %s", mgrType)
}

func addArgsToTaskFromOptions(task *execute.ExecTask, args, val string) {
	if val != "" {
		task.AddArgs(args, val)
	}
}

func getBuildkitOptions(flagSet *pflag.FlagSet) (*BuildKitOptions, error) {
	var err error
	opt := &BuildKitOptions{}
	if opt.frontend, err = flagSet.GetString(core.FrontendFlag); err != nil {
		return nil, err
	}
	if opt.addr, err = flagSet.GetString(core.AddrFlag); err != nil {
		return nil, err
	}
	if opt.tlsServerName, err = flagSet.GetString(core.TLSServerNameFlag); err != nil {
		return nil, err
	}
	if opt.tlsCaCert, err = flagSet.GetString(core.TLSCaCertFlag); err != nil {
		return nil, err
	}
	if opt.tlsCert, err = flagSet.GetString(core.TLSCertFlag); err != nil {
		return nil, err
	}
	if opt.tlsKey, err = flagSet.GetString(core.TLSKeyFlag); err != nil {
		return nil, err
	}
	if opt.tlsDir, err = flagSet.GetString(core.TLSDirFlag); err != nil {
		return nil, err
	}
	return opt, nil
}
