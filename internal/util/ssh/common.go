package ssh

import (
	"fmt"
	"terraform-provider-linux/internal/util"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func defaultErrorHandler(out []byte, err error) (util.Status, *util.CommonError) {
	if err != nil {
		return util.Failed, &util.CommonError{
			Error: err,
		}
	}
	return util.Success, nil
}

func RunCommand(linuxCtx util.LinuxContext, command string, errorhandler func([]byte, error) (util.Status, *util.CommonError)) (util.Status, string, *util.CommonError) {
	tflog.Info(linuxCtx.Ctx, fmt.Sprintf("Running command \"%s\"", command))
	var out []byte
	errors := []*util.CommonError{}

	fn := func() util.Status {
		var err error
		out, err = linuxCtx.ProviderData.SshClient.Run(command)

		status := util.Bottom
		var commonError *util.CommonError = nil

		if errorhandler != nil {
			status, commonError = errorhandler(out, err)
		}
		if status == util.Bottom {
			status, commonError = defaultErrorHandler(out, err)
		}

		if commonError != nil {
			errors = append(errors, commonError)
		}
		return status
	}
	status := util.BackoffRetry(fn, 3)
	if len(errors) != 0 {
		return status, "", util.FoldCommonError(errors)
	}

	return status, string(out), nil
}
