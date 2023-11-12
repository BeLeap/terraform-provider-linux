package lib

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/crypto/ssh"
)

type CustomSsh struct {
	Session *ssh.Session
}

func (s *CustomSsh) RunCommand(linuxCtx *LinuxContext, command string) (string, error) {
	var stdoutBuffer bytes.Buffer
	var stderrBuffer bytes.Buffer
	s.Session.Stdout = &stdoutBuffer
	s.Session.Stderr = &stderrBuffer

	tflog.Info(linuxCtx.Ctx, fmt.Sprintf("Running command \"%s\"", command))
	err := s.Session.Run(command)
	if err != nil {
		linuxCtx.Diagnostics.AddError(
			err.Error(), stderrBuffer.String(),
		)
		return "", err
	}

	return stdoutBuffer.String(), nil
}
