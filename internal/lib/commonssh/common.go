package commonssh

import (
	"fmt"
	"terraform-provider-linux/internal/lib"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func RunCommand(linuxCtx lib.LinuxContext, command string) (string, *lib.CommonError) {
	tflog.Info(linuxCtx.Ctx, fmt.Sprintf("Running command \"%s\"", command))
	var out []byte
	var err error
	fn := func() lib.Status {
		err = nil
		out, err = linuxCtx.SshClient.Run(command)
		if err != nil {
			return lib.Failed
		}
		return lib.Success
	}
	_ = lib.BackoffRetry(fn, 3)
	if err != nil {
		diagnostic := diag.NewErrorDiagnostic(err.Error(), string(out))
		return "", &lib.CommonError{
			Error:       err,
			Diagnostics: diag.Diagnostics{diagnostic},
		}
	}

	return string(out), nil
}
