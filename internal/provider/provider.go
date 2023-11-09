// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &LinuxProvider{}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &LinuxProvider{
			version: version,
		}
	}
}

type LinuxProvider struct {
	version string
}

func (p *LinuxProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "linux"
	resp.Version = p.version
}

func (p *LinuxProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *LinuxProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

func (p *LinuxProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (p *LinuxProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
