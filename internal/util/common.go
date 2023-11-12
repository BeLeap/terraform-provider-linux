package util

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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
		time.Sleep(2 << count * time.Second)
	}
}
