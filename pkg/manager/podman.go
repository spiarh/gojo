package manager

import (
	"fmt"

	"github.com/lcavajani/gojo/pkg/buildconf"
	"github.com/lcavajani/gojo/pkg/execute"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

func (p *Podman) Build(ibc *buildconf.Image) error {
	task := p.execTask
	task.AddArgs("build")
	task.AddArgs("-t", ibc.Metadata.GetFullName())

	buildArgs := ibc.GetBuildArgs()
	for arg, val := range buildArgs {
		task.AddArgs("--build-arg", fmt.Sprintf("%s=%s", arg, val))
	}
	task.AddArgs(ibc.Metadata.Context)

	_, err := task.Execute()
	if err != nil {
		return err
	}

	if p.tagLatest {
		err := p.tag(ibc.Metadata, "latest")
		if err != nil {
			return err
		}
	}

	if p.push {
		err := p.Push(ibc.Metadata)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Podman) tag(image *buildconf.ImageMeta, tag string) error {
	task := p.execTask
	task.AddArgs("tag")
	task.AddArgs(image.GetFullName(), image.GetFullNameWithTag(tag))

	_, err := task.Execute()
	if err != nil {
		return err
	}

	return nil
}

func (p *Podman) Push(image *buildconf.ImageMeta) error {
	task := p.execTask
	task.AddArgs("push")
	task.AddArgs(image.GetFullName())

	_, err := task.Execute()
	if err != nil {
		return err
	}

	if p.tagLatest {
		task.Args[len(task.Args)-1] = image.GetFullNameWithTag("latest")
		_, err := task.Execute()
		if err != nil {
			return err
		}
	}

	return nil
}
