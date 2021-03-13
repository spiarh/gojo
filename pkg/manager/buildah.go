package manager

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/lcavajani/gojo/pkg/core"
	"github.com/lcavajani/gojo/pkg/execute"
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

	if b.tagLatest {
		err := b.tag(build.Image, "latest")
		if err != nil {
			return err
		}
	}

	if b.push {
		err := b.Push(build.Image)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Buildah) tag(image *core.Image, tag string) error {
	task := b.execTask
	task.AddArgs("tag")
	task.AddArgs(image.String(), image.StringWithTag(tag))

	_, err := task.Execute()
	if err != nil {
		return err
	}

	return nil
}

func (b *Buildah) Push(image *core.Image) error {
	task := b.execTask
	task.AddArgs("push")
	task.AddArgs(image.String())

	_, err := task.Execute()
	if err != nil {
		return err
	}

	if b.tagLatest {
		task.Args[len(task.Args)-1] = image.StringWithTag("latest")
		_, err := task.Execute()
		if err != nil {
			return err
		}
	}

	return nil
}
