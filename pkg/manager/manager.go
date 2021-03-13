package manager

import (
	"fmt"

	"github.com/spf13/pflag"

	"github.com/lcavajani/gojo/pkg/core"
	"github.com/lcavajani/gojo/pkg/execute"
	"github.com/lcavajani/gojo/pkg/util"
)

type manager interface {
	Build(ibc *core.Build) error
}

var _ manager = &Podman{}

func New(flagSet *pflag.FlagSet, mgrType string) (manager, error) {
	push, err := flagSet.GetBool(core.PushFlag)
	if err != nil {
		return nil, err
	}
	tagLatest, err := flagSet.GetBool(core.TagLatestFlag)
	if err != nil {
		return nil, err
	}
	dryRun, err := flagSet.GetBool(core.DryRunFlag)
	if err != nil {
		return nil, err
	}

	streamStdio := false
	if util.IsTTYAllocated() {
		streamStdio = true
	}

	switch mgrType {
	case string(BuildahType):
		b, err := NewBuildah(push, tagLatest, dryRun, streamStdio)
		if err != nil {
			return nil, err
		}
		return b, nil
	case string(BuildkitType):
		opt, err := getBuildkitOptions(flagSet)
		if err != nil {
			return nil, err
		}
		b, err := NewBuildkit(push, tagLatest, dryRun, streamStdio, opt)
		if err != nil {
			return nil, err
		}
		return b, nil
	case string(PodmanType):
		p, err := NewPodman(push, tagLatest, dryRun, streamStdio)
		if err != nil {
			return nil, err
		}
		return p, nil
	case string(KanikoType):
		k, err := NewKaniko(push, tagLatest, dryRun, streamStdio)
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
