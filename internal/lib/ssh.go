package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/crypto/ssh"
)

type CustomSsh struct {
	Session *ssh.Session
}

type CustomSshError struct {
	Stderr   string
	RawError error
}

func (e *CustomSshError) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		return "Failed to marshal error to json"
	}
	return string(b)
}

func (s *CustomSsh) RunCommand(ctx context.Context, command string) (string, error) {
	var stdoutBuffer bytes.Buffer
	var stderrBuffer bytes.Buffer
	s.Session.Stdout = &stdoutBuffer
	s.Session.Stderr = &stderrBuffer

	tflog.Info(ctx, fmt.Sprintf("Running command \"%s\"", command))
	err := s.Session.Run(command)
	if err != nil {
		return "", &CustomSshError{
			Stderr:   stderrBuffer.String(),
			RawError: err,
		}
	}

	return stdoutBuffer.String(), nil
}
