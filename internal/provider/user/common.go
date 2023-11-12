package user

import (
	"fmt"
	"strconv"
	"strings"
	"terraform-provider-linux/internal/lib"
	"terraform-provider-linux/internal/lib/commonssh"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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

	stdout, commonError := commonssh.RunCommand(linuxCtx, "getent passwd"+" "+username)
	if commonError != nil {
		return nil, commonError
	}

	getent := strings.Split(stdout, ":")

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
		Uid: uid,
		Gid: gid,
	}, nil
}
