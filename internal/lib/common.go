package lib

import "github.com/hashicorp/terraform-plugin-framework/diag"

type CommonError struct {
	Error       error
	Diagnostics diag.Diagnostics
}
