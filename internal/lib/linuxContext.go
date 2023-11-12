package lib

import (
	"context"

	"golang.org/x/crypto/ssh"
)

type LinuxContext struct {
	Ctx        context.Context
	SshSession *ssh.Session
}

func NewLinuxContext(ctx context.Context, sshSession *ssh.Session) LinuxContext {
	return LinuxContext{
		Ctx:        ctx,
		SshSession: sshSession,
	}
}
