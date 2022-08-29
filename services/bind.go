package services

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/svex99/bind-api/pkg/setting"
)

type BindService struct {
	ctx          context.Context
	DockerCli    *client.Client
	ContainerId  string
	ZoneFilePath string
}

var Bind = &BindService{}

func init() {
	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		panic(err)
	}

	Bind.ctx = context.Background()
	Bind.DockerCli = cli
	Bind.ContainerId = setting.Bind.ContainerId
	Bind.ZoneFilePath = setting.Bind.ConfPath + "named.conf.local"
}

func (bs *BindService) exec(command []string) error {
	// was used as reference for this method the docker-cli exec command implementation
	// https://github.com/docker/cli/blob/1163b4609978e0e6f2b2629b59c4a62d348e1466/cli/command/container/exec.go#L99

	if _, err := bs.DockerCli.ContainerInspect(bs.ctx, bs.ContainerId); err != nil {
		return err
	}

	execCreateConfig := &types.ExecConfig{
		User:         "bind",
		Privileged:   false,
		Tty:          false,
		AttachStdin:  false,
		AttachStderr: false,
		AttachStdout: false,
		Detach:       true,
		DetachKeys:   "",
		Env:          []string{},
		WorkingDir:   "/",
		Cmd:          command,
	}

	response, err := bs.DockerCli.ContainerExecCreate(bs.ctx, bs.ContainerId, *execCreateConfig)
	if err != nil {
		return err
	}
	if response.ID == "" {
		return errors.New("exec ID empty")
	}

	execStartConfig := &types.ExecStartCheck{
		Detach: execCreateConfig.Detach,
		Tty:    execCreateConfig.Tty,
	}

	if err := bs.DockerCli.ContainerExecStart(bs.ctx, response.ID, *execStartConfig); err != nil {
		return err
	}

	return nil
}

// Reloads the configuration files for the bind service running inside a container
func (bs *BindService) Reload() error {
	return bs.exec([]string{"service", "named", "reload"})
}
