package user

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"terraform-provider-linux/internal/lib"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &userDataSource{}
	_ datasource.DataSourceWithConfigure = &userDataSource{}
)

func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

type userDataSource struct {
	session *lib.CustomSsh
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

	stdout, sshErr := session.RunCommand(ctx, fmt.Sprintf("getent passwd %s", username))

	if sshErr != nil {
		resp.Diagnostics.AddError(sshErr.RawError.Error(), sshErr.Error())
		return
	}

	getent := strings.Split(stdout, ":")

	uid, err := strconv.ParseInt(getent[2], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse getent uid", fmt.Sprint(err.Error()))
		return
	}

	gid, err := strconv.ParseInt(getent[3], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse getent gid", fmt.Sprint(err.Error()))
		return
	}

	state.Uid = types.Int64Value(uid)
	state.Gid = types.Int64Value(gid)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	session, ok := req.ProviderData.(*lib.CustomSsh)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Unexpected Data Source Configure Type",
		)

		return
	}

	d.session = session
}
