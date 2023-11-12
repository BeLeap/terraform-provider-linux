package lib

import (
	"context"

	"github.com/melbahja/goph"
)

type LinuxContext struct {
	Ctx       context.Context
	SshClient *goph.Client
}

func NewLinuxContext(ctx context.Context, sshClient *goph.Client) LinuxContext {
	return LinuxContext{
		Ctx:       ctx,
		SshClient: sshClient,
	}
}
