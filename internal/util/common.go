package util

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/melbahja/goph"
)

type MutlipleErrorContainer struct {
	errors []error
}

func (e *MutlipleErrorContainer) Error() string {
	errorString := ""
	for _, err := range e.errors {
		errorString = errorString + "\n" + err.Error()
	}

	return errorString
}

type CommonError struct {
	Error       error
	Diagnostics diag.Diagnostics
}

func FoldCommonError(arr []*CommonError) *CommonError {
	var errors []error
	diagnostics := diag.Diagnostics{}

	for _, commonError := range arr {
		errors = append(errors, commonError.Error)
		diagnostics = append(diagnostics, commonError.Diagnostics...)
	}
	error := &MutlipleErrorContainer{
		errors: errors,
	}

	return &CommonError{
		Error:       error,
		Diagnostics: diagnostics,
	}
}

type Status int64

const (
	Success Status = 0
	Failed  Status = 1
	Bottom  Status = -1
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

func ConvertProviderData(providerData any) (*LinuxProviderData, *CommonError) {
	if providerData == nil {
		return nil, nil
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
