// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net"
	"terraform-provider-linux/internal/lib"
	"terraform-provider-linux/internal/provider/user"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/crypto/ssh"
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

type LinuxProviderModel struct {
	Host       types.String `tfsdk:"host"`
	Username   types.String `tfsdk:"username"`
	PrivateKey types.String `tfsdk:"private_key"`
}

func (p *LinuxProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Required: true,
			},
			"username": schema.StringAttribute{
				Required: true,
			},
			"private_key": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *LinuxProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config LinuxProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Host unknown",
			"Host is unknown",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Username unknown",
			"Username is unknown",
		)
	}

	if config.PrivateKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("private_key"),
			"PrivateKey unknown",
			"PrivateKey is unknown",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	var host string
	var username string
	var privateKey string

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.PrivateKey.IsNull() {
		privateKey = config.PrivateKey.ValueString()
	}

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing host",
			"Host is required",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing username",
			"Username is required",
		)
	}

	if privateKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("private_key"),
			"Missing private_key",
			"PrivateKey is required",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	key, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		resp.Diagnostics.AddError("Wrong private_key", "Failed to parse private key")
		return
	}

	sshClientConfig := &ssh.ClientConfig{
		User:            username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}
	sshClient, err := ssh.Dial("tcp", net.JoinHostPort(host, "22"), sshClientConfig)
	if err != nil {
		resp.Diagnostics.AddError("Failed to dial host", "Failed to dial host")
		return
	}

	session, err := sshClient.NewSession()
	if err != nil {
		resp.Diagnostics.AddError("Failed to create ssh session", "Failed to create ssh session")
		return
	}

	resp.DataSourceData = &lib.CustomSsh{
		Session: session,
	}
	resp.ResourceData = &lib.CustomSsh{
		Session: session,
	}
}

func (p *LinuxProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		user.NewUserDataSource,
	}
}

func (p *LinuxProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		user.NewUserResource,
	}
}
