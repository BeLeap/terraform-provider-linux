package file

import (
	"terraform-provider-linux/internal/util"
	sshUtil "terraform-provider-linux/internal/util/ssh"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LinuxFile struct {
	Path string
}

type LinuxFileModel struct {
	Path types.String `tfsdk:"path"`
}

func NewLinuxFileModel(linuxDirectory *LinuxFile) LinuxFileModel {
	return LinuxFileModel{
		Path: types.StringValue(linuxDirectory.Path),
	}
}

func Get(linuxCtx util.LinuxContext, path string) (*LinuxFile, *util.CommonError) {
	errorhandler := func(out []byte, err error) (util.Status, *util.CommonError) {
		switch err.Error() {
		case "Process exited with status 1":
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
		default:
			return sshUtil.DefaultErrorHandler(out, err)
		}
	}
	_, commonError := sshUtil.RunCommand(linuxCtx, "getfacl"+" "+path, errorhandler)
	if commonError != nil {
		return nil, commonError
	}
	return &LinuxFile{
		Path: path,
	}, nil
}
