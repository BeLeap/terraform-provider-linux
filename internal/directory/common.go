package directory

import (
	"terraform-provider-linux/internal/util"
	sshUtil "terraform-provider-linux/internal/util/ssh"

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
	_, commonError := sshUtil.RunCommand(linuxCtx, "get facl"+" "+path, nil)
	if commonError != nil {
		return nil, commonError
	}
	return &LinuxDirectory{
		Path: path,
	}, nil
}
