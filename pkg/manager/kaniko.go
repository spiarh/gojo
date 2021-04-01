package manager

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spiarh/gojo/pkg/core"
	"github.com/spiarh/gojo/pkg/execute"
)

type Kaniko struct {
	log      zerolog.Logger
	execTask execute.ExecTask

	push      bool
	tagLatest bool
}

func NewKaniko(push, tagLatest, dryRun, streamStdio bool) (*Kaniko, error) {
	logger := log.With().Str("manager", string(KanikoType)).Logger()
	return &Kaniko{
		log: logger,
		execTask: execute.ExecTask{
			Log:         logger,
			Command:     "/kaniko/executor",
			StreamStdio: streamStdio,
			DryRun:      dryRun,
		},
		push:      push,
		tagLatest: tagLatest,
	}, nil
}

func (k *Kaniko) Build(build *core.Build) error {
	task := k.execTask

	buildArgs := build.GetBuildArgs()
	for arg, val := range buildArgs {
		task.AddArgs("--build-arg", fmt.Sprintf("%s=%s", arg, val))
	}
	task.AddArgs("--context", build.Image.Context)
	task.AddArgs("--dockerfile", build.Image.Containerfile)
	task.AddArgs("--destination", build.Image.String())

	if k.tagLatest {
		task.AddArgs("--destination", build.Image.StringWithTagLatest())
	}

	if _, err := task.Execute(); err != nil {
		return err
	}

	return nil
}
