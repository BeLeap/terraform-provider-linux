package commonssh

import (
	"fmt"
	"terraform-provider-linux/internal/lib"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func RunCommand(linuxCtx lib.LinuxContext, command string, errorhandler func([]byte, error) (lib.Status, error)) (string, *lib.CommonError) {
	tflog.Info(linuxCtx.Ctx, fmt.Sprintf("Running command \"%s\"", command))
	var out []byte
	var errors []error
	errors = []error{}

	fn := func() lib.Status {
		var err error
		out, err = linuxCtx.SshClient.Run(command)
		if err != nil {
			status, err := errorhandler(out, err)
			if err != nil {
				errors = append(errors, err)
			}
			return status
		}
		return lib.Success
	}
	_ = lib.BackoffRetry(fn, 3)
	if len(errors) != 0 {
		var diagnostics diag.Diagnostics
		var error error
		for i, err := range errors {
			diagnostics = append(diagnostics, diag.NewErrorDiagnostic(fmt.Sprintf("Retry %d: %v", i+1, err), string(out)))
			error = err
		}
		return "", &lib.CommonError{
			Error:       error,
			Diagnostics: diagnostics,
		}
	}

	return string(out), nil
}
