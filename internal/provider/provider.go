package provider

import (
	"context"
	"fmt"
	"terraform-provider-linux/internal/file"
	"terraform-provider-linux/internal/user"
	"terraform-provider-linux/internal/util"

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
	Port       types.Int64  `tfsdk:"port"`
	Username   types.String `tfsdk:"username"`
	PrivateKey types.String `tfsdk:"private_key"`
	Password   types.String `tfsdk:"password"`
}

func (p *LinuxProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Required: true,
			},
			"port": schema.Int64Attribute{
				Optional: true,
			},
			"username": schema.StringAttribute{
				Required: true,
			},
			"private_key": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
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

	if config.Port.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("port"),
			"Port unknown",
			"Port is unknown",
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

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Password unknown",
			"Password is unknown",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	var host string
	port := 22
	var username string
	var privateKey string
	var password string

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Port.IsNull() {
		port = int(config.Port.ValueInt64())
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.PrivateKey.IsNull() {
		privateKey = config.PrivateKey.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Empty host",
			"Please specify host",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Empty username",
			"Please specify username",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	authMethods := []ssh.AuthMethod{}
	if privateKey != "" {
		signer, err := ssh.ParsePrivateKey([]byte(privateKey))
		if err != nil {
			resp.Diagnostics.AddError("Failed to parse private key", err.Error())
			return
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}
	if password != "" {
		authMethods = append(authMethods, ssh.Password(password))
	}

	if len(authMethods) == 0 {
		resp.Diagnostics.AddError("Empty auth info", "Please specify either private key or password.")
	}

	sshClientConfig := &ssh.ClientConfig{
		User:            username,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), sshClientConfig)
	if err != nil {
		resp.Diagnostics.AddError("Failed to connect", err.Error())
		return
	}

	session, err := conn.NewSession()
	if err != nil {
		resp.Diagnostics.AddError("Failed to create session", err.Error())
		return
	}

	providerData := &util.LinuxProviderData{
		SshSession: session,
	}
	resp.DataSourceData = providerData
	resp.ResourceData = providerData
}

func (p *LinuxProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		user.NewUserDataSource,
		file.NewFileDataSource,
	}
}

func (p *LinuxProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		user.NewUserResource,
	}
}
