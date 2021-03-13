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

func NewPodman(push, tagLatest, dryRun bool) (*Podman, error) {
	logger := log.With().Str("manager", string(PodmanType)).Logger()
	return &Podman{
		log: logger,
		execTask: execute.ExecTask{
			Log:         logger,
			Command:     "podman",
			StreamStdio: true,
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

	buildArgs := build.GetBuildArgs()
	for arg, val := range buildArgs {
		task.AddArgs("--build-arg", fmt.Sprintf("%s=%s", arg, val))
	}
	task.AddArgs(build.Image.Context)
	task.AddArgs("-f", build.Image.Containerfile)

	_, err := task.Execute()
	if err != nil {
		return err
	}

	if p.tagLatest {
		err := p.tag(build.Image, "latest")
		if err != nil {
			return err
		}
	}

	if p.push {
		err := p.Push(build.Image)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Podman) tag(image *core.Image, tag string) error {
	task := p.execTask
	task.AddArgs("tag")
	task.AddArgs(image.String(), image.StringWithTag(tag))

	_, err := task.Execute()
	if err != nil {
		return err
	}

	return nil
}

func (p *Podman) Push(image *core.Image) error {
	task := p.execTask
	task.AddArgs("push")
	task.AddArgs(image.String())

	_, err := task.Execute()
	if err != nil {
		return err
	}

	if p.tagLatest {
		task.Args[len(task.Args)-1] = image.StringWithTag("latest")
		_, err := task.Execute()
		if err != nil {
			return err
		}
	}

	return nil
}
