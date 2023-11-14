package util

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/melbahja/goph"
)

type CommonError struct {
	Error       error
	Diagnostics diag.Diagnostics
}

type Status int64

const (
	Success Status = 0
	Failed  Status = 1
)

func BackoffRetry(fn func() Status, retry int) Status {
	count := 0

	for {
		result := fn()

		if result == Success {
			return Success
		}

		count = count + 1
		if count >= retry {
			return result
		}
		time.Sleep(time.Duration(2) * time.Second << count)
	}
}

type LinuxProviderData struct {
	SshClient *goph.Client
}

func ConvetProviderData(providerData any) (*LinuxProviderData, *CommonError) {
	if providerData == nil {
		diagnostic := diag.NewErrorDiagnostic(
			"Empty ProviderData",
			"ProviderData is empty. Please check.",
		)
		return nil, &CommonError{
			Diagnostics: diag.Diagnostics{diagnostic},
		}
	}

	linuxProviderData, ok := providerData.(*LinuxProviderData)
	if !ok {
		diagnostic := diag.NewErrorDiagnostic(
			"ProviderData type assertion failed",
			"Expected ProviderData to be *util.LinuxProviderData, got different type",
		)
		return nil, &CommonError{
			Diagnostics: diag.Diagnostics{diagnostic},
		}
	}
	return linuxProviderData, nil
}
