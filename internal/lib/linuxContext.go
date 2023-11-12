package lib

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type LinuxContext struct {
	Ctx         context.Context
	Ssh         *CustomSsh
	Diagnostics diag.Diagnostics
}

func NewLinuxContext(ctx context.Context, ssh *CustomSsh, diagnostics diag.Diagnostics) *LinuxContext {
	return &LinuxContext{
		Ctx:         ctx,
		Ssh:         ssh,
		Diagnostics: diagnostics,
	}
}
