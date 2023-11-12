package user

import (
	"context"
	"terraform-provider-linux/internal/lib"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

type userDataSourceModel struct {
	Username types.String `tfsdk:"username"`
	Uid      types.Int64  `tfsdk:"uid"`
	Gid      types.Int64  `tfsdk:"gid"`
}

func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	session := d.session
	var state userDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	linuxCtx := lib.NewLinuxContext(ctx, session)

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

	state.Uid = types.Int64Value(user.Uid)
	state.Gid = types.Int64Value(user.Gid)

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
