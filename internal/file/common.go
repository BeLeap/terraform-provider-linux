package file

import (
	"strings"
	"terraform-provider-linux/internal/util"
	sshUtil "terraform-provider-linux/internal/util/ssh"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LinuxFile struct {
	Path string
	Type string
	Acl  *Facl
}

type LinuxFileModel struct {
	Path types.String `tfsdk:"path"`
	Type types.String `tfsdk:"type"`
	Acl  FaclModel    `tfsdk:"acl"`
}

func NewLinuxFileModel(linuxFile *LinuxFile) LinuxFileModel {
	return LinuxFileModel{
		Path: types.StringValue(linuxFile.Path),
		Type: types.StringValue(linuxFile.Type),
		Acl:  newFaclModel(linuxFile.Acl),
	}
}

type Facl struct {
	User  int64
	Group int64
	Other int64
}
type FaclModel struct {
	User  types.Int64 `tfsdk:"user"`
	Group types.Int64 `tfsdk:"group"`
	Other types.Int64 `tfsdk:"other"`
}

func newFaclModel(facl *Facl) FaclModel {
	return FaclModel{
		User:  types.Int64Value(facl.User),
		Group: types.Int64Value(facl.Group),
		Other: types.Int64Value(facl.Other),
	}
}

func Get(linuxCtx util.LinuxContext, file *LinuxFile) (*LinuxFile, *util.CommonError) {
	errorhandler := func(out []byte, err error) (util.Status, *util.CommonError) {
		switch err.Error() {
		case "Process exited with status 1":
			diagnostic := diag.NewErrorDiagnostic(
				"Path not found",
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
	stdout, commonError := sshUtil.RunCommand(linuxCtx, "getfacl -t"+" "+file.Path, errorhandler)
	if commonError != nil {
		return nil, commonError
	}

	lines := strings.Split(stdout, "\n")

	var userAclString string
	var groupAclString string
	var otherAclString string
	for _, line := range lines {
		if strings.HasPrefix(line, "USER") {
		}
	}

	return &LinuxFile{
		Path: file.Path,
		Type: file.Type,
	}, nil
}
