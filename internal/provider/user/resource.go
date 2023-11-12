package user

import (
	"context"
	"fmt"
	"terraform-provider-linux/internal/lib"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &userResource{}
	_ resource.ResourceWithConfigure = &userResource{}
)

func NewUserResource() resource.Resource {
	return &userResource{}
}

type userResource struct {
	session *lib.CustomSsh
}

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				Required: true,
			},
			"uid": schema.Int64Attribute{
				Computed: true,
				Optional: true,
			},
			"gid": schema.Int64Attribute{
				Computed: true,
				Optional: true,
			},
		},
	}
}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	command := "useradd"

	var username string
	if plan.Username.IsUnknown() || plan.Username.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Something wrong in username",
			"Something wrong in username",
		)
		return
	}
	username = plan.Username.ValueString()
	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Empty username is not allowed",
			"Empty username is not allowed",
		)
		return
	}

	command = command + " " + username

	if !plan.Uid.IsUnknown() && !plan.Uid.IsNull() {
		command = command + " " + "--uid" + " " + fmt.Sprintf("%d", plan.Uid.ValueInt64())
	}
	if !plan.Gid.IsUnknown() && !plan.Uid.IsNull() {
		command = command + " " + "--gid" + " " + fmt.Sprintf("%d", plan.Gid.ValueInt64())
	}

	_, sshErr := r.session.RunCommand(ctx, command)
	if sshErr != nil {
		resp.Diagnostics.AddError(sshErr.RawError.Error(), sshErr.Error())
		return
	}
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	session, ok := req.ProviderData.(*lib.CustomSsh)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *lib.CustomSsh, got: %T. Please report this this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.session = session
}

type userResourceModel struct {
	Username types.String `tfsdk:"username"`
	Uid      types.Int64  `tfsdk:"uid"`
	Gid      types.Int64  `tfsdk:"gid"`
}