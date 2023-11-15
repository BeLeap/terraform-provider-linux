package directory

import (
	"terraform-provider-linux/internal/util"
	sshUtil "terraform-provider-linux/internal/util/ssh"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LinuxDirectory struct {
	Path string
}

type LinuxDirectoryModel struct {
	Path types.String `tfsdk:"path"`
}

func NewLinuxDirectoryModel(linuxDirectory *LinuxDirectory) LinuxDirectoryModel {
	return LinuxDirectoryModel{
		Path: types.StringValue(linuxDirectory.Path),
	}
}

func Get(linuxCtx util.LinuxContext, path string) (*LinuxDirectory, *util.CommonError) {
	errorhandler := func(out []byte, err error) (util.Status, *util.CommonError) {
		if err.Error() == "Process exited with status 1" {
			diagnostic := diag.NewErrorDiagnostic(
				"Directory not found",
				"Please check path",
			)
			return util.Success, &util.CommonError{
				Error: err,
				Diagnostics: diag.Diagnostics{
					diagnostic,
				},
			}
		}
		return sshUtil.DefaultErrorHandler(out, err)
	}
	_, commonError := sshUtil.RunCommand(linuxCtx, "getfacl"+" "+path, errorhandler)
	if commonError != nil {
		return nil, commonError
	}
	return &LinuxDirectory{
		Path: path,
	}, nil
}
