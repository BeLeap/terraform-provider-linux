package customssh

import (
	"bytes"
	"fmt"
	"terraform-provider-linux/internal/lib"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func RunCommand(linuxCtx lib.LinuxContext, command string) (string, *lib.CommonError) {
	var stdoutBuffer bytes.Buffer
	var stderrBuffer bytes.Buffer
	linuxCtx.SshSession.Stdout = &stdoutBuffer
	linuxCtx.SshSession.Stderr = &stderrBuffer

	tflog.Info(linuxCtx.Ctx, fmt.Sprintf("Running command \"%s\"", command))
	err := linuxCtx.SshSession.Run(command)
	if err != nil {
		diagnostic := diag.NewErrorDiagnostic(err.Error(), stderrBuffer.String())
		return "", &lib.CommonError{
			Error:       err,
			Diagnostics: diag.Diagnostics{diagnostic},
		}
	}

	return stdoutBuffer.String(), nil
}
