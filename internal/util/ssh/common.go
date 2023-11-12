package ssh

import (
	"fmt"
	"terraform-provider-linux/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func defaultErrorHandler(out []byte, err error) (util.Status, error) { return util.Failed, err }

func RunCommand(linuxCtx util.LinuxContext, command string, errorhandler func([]byte, error) (util.Status, error)) (string, *util.CommonError) {
	errorHandlerCoerced := errorhandler
	if errorHandlerCoerced == nil {
		errorHandlerCoerced = defaultErrorHandler
	}

	tflog.Info(linuxCtx.Ctx, fmt.Sprintf("Running command \"%s\"", command))
	var out []byte
	var errors []error
	errors = []error{}

	fn := func() util.Status {
		var err error
		out, err = linuxCtx.SshClient.Run(command)
		if err != nil {
			status, err := errorHandlerCoerced(out, err)
			if err != nil {
				errors = append(errors, err)
			}
			return status
		}
		return util.Success
	}
	_ = util.BackoffRetry(fn, 3)
	if len(errors) != 0 {
		var diagnostics diag.Diagnostics
		var error error
		for i, err := range errors {
			diagnostics = append(diagnostics, diag.NewErrorDiagnostic(fmt.Sprintf("Retry %d: %v", i+1, err), string(out)))
			error = err
		}
		return "", &util.CommonError{
			Error:       error,
			Diagnostics: diagnostics,
		}
	}

	return string(out), nil
}
