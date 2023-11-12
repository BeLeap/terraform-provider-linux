package user

import (
	"context"
	"terraform-provider-linux/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/melbahja/goph"
)

var (
	_ datasource.DataSource              = &userDataSource{}
	_ datasource.DataSourceWithConfigure = &userDataSource{}
)

func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

type userDataSource struct {
	session *goph.Client
}

func (d *userDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *userDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				Required: true,
			},
			"uid": schema.Int64Attribute{
				Computed: true,
			},
			"gid": schema.Int64Attribute{
				Computed: true,
			},
		},
	}
}

func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	linuxCtx := util.NewLinuxContext(ctx, d.session)

	var state LinuxUserModel

	diags := req.Config.Get(linuxCtx.Ctx, &state)
	resp.Diagnostics.Append(diags...)

	if state.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Username unknown",
			"Username unknown",
		)
		return
	}

	var username string
	if !state.Username.IsNull() {
		username = state.Username.ValueString()
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing username",
			"Missing username",
		)
		return
	}

	user, commonError := GetUser(linuxCtx, username)
	if commonError != nil {
		resp.Diagnostics.Append(commonError.Diagnostics...)
		return
	}
	if user == nil {
		resp.Diagnostics.AddError("User not found", "Check user exists on server")
		return
	}

	state = NewLinuxUserModel(user)

	diags = resp.State.Set(linuxCtx.Ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	session, ok := req.ProviderData.(*goph.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Unexpected Data Source Configure Type",
		)

		return
	}

	d.session = session
}
