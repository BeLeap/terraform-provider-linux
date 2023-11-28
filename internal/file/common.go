package file

import (
	"errors"
	"fmt"
	"strconv"
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
	User  *FaclLine
	Group *FaclLine
	Other *FaclLine
}
type FaclModel struct {
	User  FaclLineModel `tfsdk:"user"`
	Group FaclLineModel `tfsdk:"group"`
	Other FaclLineModel `tfsdk:"other"`
}

func newFaclModel(facl *Facl) FaclModel {
	return FaclModel{
		User:  newFaclLineModel(facl.User),
		Group: newFaclLineModel(facl.Group),
		Other: newFaclLineModel(facl.Other),
	}
}

type FaclLine struct {
	Id        int64
	Permisson int64
}
type FaclLineModel struct {
	Id        types.Int64 `tfsdk:"id"`
	Permisson types.Int64 `tfsdk:"permisson"`
}

func newFaclLineModel(faclLine *FaclLine) FaclLineModel {
	return FaclLineModel{
		Id:        types.Int64Value(faclLine.Id),
		Permisson: types.Int64Value(faclLine.Permisson),
	}
}

type PermissionType int64

const (
	Read         PermissionType = 4
	Write        PermissionType = 2
	Execute      PermissionType = 1
	NoPermission PermissionType = 0
	Invalid      PermissionType = -1
)

func mapPermissionTypeToString(permissionType PermissionType) (string, error) {
	switch permissionType {
	case Read:
		return "r", nil
	case Write:
		return "w", nil
	case Execute:
		return "x", nil
	case NoPermission:
		return "-", nil
	default:
		return "", errors.New(fmt.Sprintf("There is no string to map %d", permissionType))
	}
}
func parsePermissionTypeString(in string, permissionType PermissionType) (string, PermissionType, error) {
	permissionTypeString, err := mapPermissionTypeToString(permissionType)
	if err != nil {
		return "", 0, err
	}
	after, found := strings.CutPrefix(in, permissionTypeString)
	if found {
		return after, permissionType, nil
	}

	noPermissionString, err := mapPermissionTypeToString(NoPermission)
	if err != nil {
		return "", 0, err
	}
	after, found = strings.CutPrefix(in, noPermissionString)
	if found {
		return after, NoPermission, nil
	}

	return "", 0, errors.New(fmt.Sprintf("Invalid permission type string \"%s\" provided", in))
}

func parseFacl(content string) (*Facl, error) {
	lines := strings.Split(content, "\n")

	parsePermissionString := func(in string) (PermissionType, error) {
		permissionString := in
		var err error

		permission := NoPermission
		permissionForType := NoPermission

		permissionString, permissionForType, err = parsePermissionTypeString(permissionString, Read)
		if err != nil {
			return Invalid, err
		}
		permission = permission + permissionForType

		permissionString, permissionForType, err = parsePermissionTypeString(permissionString, Write)
		if err != nil {
			return Invalid, err
		}
		permission = permission + permissionForType

		permissionString, permissionForType, err = parsePermissionTypeString(permissionString, Execute)
		if err != nil {
			return Invalid, err
		}
		permission = permission + permissionForType

		return permission, nil
	}

	var userAcl *FaclLine
	var groupAcl *FaclLine
	var otherAcl *FaclLine

	for _, line := range lines {
		if after, found := strings.CutPrefix(line, "USER"); found {
			splitted := strings.Split(after, " ")

			id, err := strconv.ParseInt(splitted[0], 10, 64)
			if err != nil {
				return nil, err
			}

			permission, err := parsePermissionString(splitted[1])
			if err != nil {
				return nil, err
			}

			userAcl = &FaclLine{
				Id:        id,
				Permisson: int64(permission),
			}
		}
		if after, found := strings.CutPrefix(line, "GROUP"); found {
			splitted := strings.Split(after, " ")

			id, err := strconv.ParseInt(splitted[0], 10, 64)
			if err != nil {
				return nil, err
			}

			permission, err := parsePermissionString(splitted[1])
			if err != nil {
				return nil, err
			}

			groupAcl = &FaclLine{
				Id:        id,
				Permisson: int64(permission),
			}
		}
		if after, found := strings.CutPrefix(line, "other"); found {
			splitted := strings.Split(after, " ")

			permission, err := parsePermissionString(splitted[1])
			if err != nil {
				return nil, err
			}

			otherAcl = &FaclLine{
				Id:        -1,
				Permisson: int64(permission),
			}
		}
	}

	return &Facl{
		User:  userAcl,
		Group: groupAcl,
		Other: otherAcl,
	}, nil
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
	stdout, commonError := sshUtil.RunCommand(linuxCtx, "getfacl -nt"+" "+file.Path, errorhandler)
	if commonError != nil {
		return nil, commonError
	}
	acl, err := parseFacl(stdout)
	if err != nil {
		return nil, &util.CommonError{
			Error: err,
			Diagnostics: diag.Diagnostics{
				diag.NewErrorDiagnostic("Failed to parse facl", fmt.Sprintf("Failed to parse facl content: %s", stdout)),
			},
		}
	}

	return &LinuxFile{
		Path: file.Path,
		Type: file.Type,
		Acl:  acl,
	}, nil
}
