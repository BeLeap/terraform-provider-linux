package ssh

import (
	"fmt"
	"terraform-provider-linux/internal/util"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func defaultErrorHandler(out []byte, err error) (util.Status, *util.CommonError) {
	return util.Failed, &util.CommonError{
		Error: err,
	}
}

func RunCommand(linuxCtx util.LinuxContext, command string, errorhandler func([]byte, error) (util.Status, *util.CommonError)) (string, *util.CommonError) {
	errorHandlerCoerced := errorhandler
	if errorHandlerCoerced == nil {
		errorHandlerCoerced = defaultErrorHandler
	}

	tflog.Info(linuxCtx.Ctx, fmt.Sprintf("Running command \"%s\"", command))
	var out []byte
	errors := []*util.CommonError{}

	fn := func() util.Status {
		var err error
		out, err = linuxCtx.ProviderData.SshClient.Run(command)
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
		return "", util.FoldCommonError(errors)
	}

	return string(out), nil
}
