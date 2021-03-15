package manager

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"

	"github.com/lcavajani/gojo/pkg/core"
	"github.com/lcavajani/gojo/pkg/execute"
)

type Buildkit struct {
	log      zerolog.Logger
	execTask execute.ExecTask

	push      bool
	tagLatest bool
}

type BuildKitOptions struct {
	frontend      string
	addr          string
	tlsServerName string
	tlsCaCert     string
	tlsCert       string
	tlsKey        string
	tlsDir        string
}

func NewBuildkit(push, tagLatest, dryRun, streamStdio bool, opt *BuildKitOptions) (*Buildkit, error) {
	logger := log.With().Str("manager", string(BuildkitType)).Logger()
	task := execute.ExecTask{
		Log:         logger,
		Command:     "buildctl",
		StreamStdio: streamStdio,
		DryRun:      dryRun,
	}

	// task.AddArgs("--debug")
	addArgsToTaskFromOptions(&task, "--addr", opt.addr)
	addArgsToTaskFromOptions(&task, "--tlsservername", opt.tlsServerName)
	addArgsToTaskFromOptions(&task, "--tlscacert", opt.tlsCaCert)
	addArgsToTaskFromOptions(&task, "--tlscert", opt.tlsCert)
	addArgsToTaskFromOptions(&task, "--tlskey", opt.tlsKey)
	addArgsToTaskFromOptions(&task, "--tlsdir", opt.tlsDir)

	return &Buildkit{
		log:       logger,
		execTask:  task,
		push:      push,
		tagLatest: tagLatest,
	}, nil
}

func (b *Buildkit) Build(build *core.Build) error {
	task := b.execTask
	task.AddArgs("build")

	task.AddArgs("--frontend", "dockerfile.v0")
	task.AddArgs("--opt", "filename="+build.Image.Containerfile)
	task.AddArgs("--local", "context="+build.Image.Context)
	task.AddArgs("--local", "dockerfile="+build.Image.Context)

	buildArgs := build.GetBuildArgs()
	for arg, val := range buildArgs {
		task.AddArgs("--opt", fmt.Sprintf("build-arg:%s=%s", arg, val))
	}

	task.AddArgs("--output",
		fmt.Sprintf("type=image,name=%s,push=%t", build.Image.String(), b.push))

	_, err := task.Execute()
	if err != nil {
		return err
	}

	if b.tagLatest {
		// rerun with a new tag
		task.Args = task.Args[:len(task.Args)-1]
		task.AddArgs(fmt.Sprintf("type=image,name=%s,push=%t", build.Image.StringWithTagLatest(), b.push))
		_, err := task.Execute()
		if err != nil {
			return err
		}
	}

	return nil
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
