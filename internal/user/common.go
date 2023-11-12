package user

import (
	"fmt"
	"strconv"
	"strings"
	"terraform-provider-linux/lib"
	"terraform-provider-linux/lib/commonssh"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LinuxUser struct {
	Username string
	Uid      int64
	Gid      int64
}

func GetUser(linuxCtx lib.LinuxContext, username string) (*LinuxUser, *lib.CommonError) {
	if username == "" {
		diagnoistic := diag.NewErrorDiagnostic("Empty username", "Please specify username")
		return nil, &lib.CommonError{
			Error:       nil,
			Diagnostics: diag.Diagnostics{diagnoistic},
		}
	}

	errorhandler := func(out []byte, err error) (lib.Status, error) {
		if err.Error() == "Process exited with status 2" {
			return lib.Success, nil
		}
		return lib.Failed, err
	}
	stdout, commonError := commonssh.RunCommand(linuxCtx, "getent passwd"+" "+username, errorhandler)
	if commonError != nil {
		return nil, commonError
	}

	getent := strings.Split(stdout, ":")
	if len(getent) != 7 {
		return nil, nil
	}

	uid, err := strconv.ParseInt(getent[2], 10, 64)
	if err != nil {
		diagnostic := diag.NewErrorDiagnostic("Failed to parse getent uid", fmt.Sprint(err.Error()))
		return nil, &lib.CommonError{
			Error:       err,
			Diagnostics: diag.Diagnostics{diagnostic},
		}
	}

	gid, err := strconv.ParseInt(getent[3], 10, 64)
	if err != nil {
		diagnostic := diag.NewErrorDiagnostic("Failed to parse getent gid", fmt.Sprint(err.Error()))
		return nil, &lib.CommonError{
			Error:       err,
			Diagnostics: diag.Diagnostics{diagnostic},
		}
	}

	return &LinuxUser{
		Username: username,
		Uid:      uid,
		Gid:      gid,
	}, nil
}

type LinuxUserModel struct {
	Username types.String `tfsdk:"username"`
	Uid      types.Int64  `tfsdk:"uid"`
	Gid      types.Int64  `tfsdk:"gid"`
}

func NewLinuxUserModel(user *LinuxUser) LinuxUserModel {
	return LinuxUserModel{
		Username: types.StringValue(user.Username),
		Uid:      types.Int64Value(user.Uid),
		Gid:      types.Int64Value(user.Gid),
	}
}
