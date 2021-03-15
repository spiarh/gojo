package manager

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/lcavajani/gojo/pkg/core"
	"github.com/lcavajani/gojo/pkg/execute"
)

type Podman struct {
	log      zerolog.Logger
	execTask execute.ExecTask

	push      bool
	tagLatest bool
}

func NewPodman(push, tagLatest, dryRun, streamStdio bool) (*Podman, error) {
	logger := log.With().Str("manager", string(PodmanType)).Logger()

	return &Podman{
		log: logger,
		execTask: execute.ExecTask{
			Log:         logger,
			Command:     "podman",
			StreamStdio: streamStdio,
			DryRun:      dryRun,
		},
		push:      push,
		tagLatest: tagLatest,
	}, nil
}

func (p *Podman) Build(build *core.Build) error {
	task := p.execTask
	task.AddArgs("build")

	task.AddArgs("-t", build.Image.String())
	if p.tagLatest {
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

	if p.push {
		if err := p.Push(build.Image); err != nil {
			return err
		}
	}

	return nil
}

func (p *Podman) Push(image *core.Image) error {
	task := p.execTask
	task.AddArgs("push")
	task.AddArgs(image.String())

	if _, err := task.Execute(); err != nil {
		return err
	}

	if p.tagLatest {
		task.Args[len(task.Args)-1] = image.StringWithTagLatest()
		if _, err := task.Execute(); err != nil {
			return err
		}
	}

	return nil
}
