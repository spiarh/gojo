package manager

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spiarh/gojo/pkg/core"
	"github.com/spiarh/gojo/pkg/execute"
)

type Buildah struct {
	log      zerolog.Logger
	execTask execute.ExecTask

	push      bool
	tagLatest bool
}

func NewBuildah(push, tagLatest, dryRun, streamStdio bool) (*Buildah, error) {
	logger := log.With().Str("manager", string(BuildahType)).Logger()
	return &Buildah{
		log: logger,
		execTask: execute.ExecTask{
			Log:         logger,
			Command:     "buildah",
			StreamStdio: streamStdio,
			DryRun:      dryRun,
		},
		push:      push,
		tagLatest: tagLatest,
	}, nil
}

func (b *Buildah) Build(build *core.Build) error {
	task := b.execTask
	task.AddArgs("bud")

	task.AddArgs("-t", build.Image.String())
	if b.tagLatest {
		task.AddArgs("-t", build.Image.StringWithTagLatest())
	}

	buildArgs := build.GetBuildArgs()
	for arg, val := range buildArgs {
		task.AddArgs("--build-arg", fmt.Sprintf("%s=%s", arg, val))
	}
	task.AddArgs("-f", build.Image.Containerfile)
	task.AddArgs(build.Image.Context)

	if _, err := task.Execute(); err != nil {
		return err
	}

	if b.push {
		if err := b.Push(build.Image); err != nil {
			return err
		}
	}

	return nil
}

func (b *Buildah) Push(image *core.Image) error {
	task := b.execTask
	task.AddArgs("push")
	task.AddArgs(image.String())

	if _, err := task.Execute(); err != nil {
		return err
	}

	if b.tagLatest {
		task.Args[len(task.Args)-1] = image.StringWithTagLatest()
		if _, err := task.Execute(); err != nil {
			return err
		}
	}

	return nil
}
