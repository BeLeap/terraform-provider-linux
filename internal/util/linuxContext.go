package util

import (
	"context"
)

type LinuxContext struct {
	Ctx          context.Context
	ProviderData *LinuxProviderData
}

func NewLinuxContext(ctx context.Context, providerData *LinuxProviderData) LinuxContext {
	return LinuxContext{
		Ctx:          ctx,
		ProviderData: providerData,
	}
}
