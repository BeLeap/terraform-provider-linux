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

type customSshError struct {
	Stderr   string
	RawError error
}

func (e *customSshError) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		return "Failed to marshal error to json"
	}
	return string(b)
}

func (s *CustomSsh) RunCommand(ctx context.Context, command string) (string, *customSshError) {
	var stdoutBuffer bytes.Buffer
	var stderrBuffer bytes.Buffer
	s.Session.Stdout = &stdoutBuffer
	s.Session.Stderr = &stderrBuffer

	tflog.Info(ctx, fmt.Sprintf("Running command \"%s\"", command))
	err := s.Session.Run(command)
	if err != nil {
		return "", &customSshError{
			Stderr:   stderrBuffer.String(),
			RawError: err,
		}
	}

	return stdoutBuffer.String(), nil
}
