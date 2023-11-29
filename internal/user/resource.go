package user

import (
	"context"
	"fmt"
	"terraform-provider-linux/internal/util"
	sshUtil "terraform-provider-linux/internal/util/ssh"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

func NewUserResource() resource.Resource {
	return &userResource{}
}

type userResource struct {
	providerData *util.LinuxProviderData
}

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
	var plan LinuxUserModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	linuxCtx := util.NewLinuxContext(ctx, r.providerData)

	command := "useradd"

	if plan.Username.IsUnknown() || plan.Username.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Something wrong in username",
			"Something wrong in username",
		)
		return
	}
	username := plan.Username.ValueString()
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

	_, _, commonError := sshUtil.RunCommand(linuxCtx, command, nil)
	if commonError != nil {
		resp.Diagnostics.Append(commonError.Diagnostics...)
		return
	}

	user, commonError := Get(linuxCtx, plan.Username.ValueString())
	if commonError != nil {
		resp.Diagnostics.Append(commonError.Diagnostics...)
		return
	}
	if user == nil {
		resp.Diagnostics.AddError("Failed to created user.", "User not exists after creation request")
		return
	}

	plan = NewLinuxUserModel(user)
	diags = resp.State.Set(linuxCtx.Ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	linuxCtx := util.NewLinuxContext(ctx, r.providerData)

	var state LinuxUserModel
	diags := req.State.Get(linuxCtx.Ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, commonError := Get(linuxCtx, state.Username.ValueString())
	if commonError != nil {
		resp.Diagnostics.Append(commonError.Diagnostics...)
		return
	}
	if user == nil {
		// TODO: Find way to update resource if user deleted outside
		resp.Diagnostics.AddError("User not found", "This indicates user created with terraform deleted outside")
		return
	}

	state = NewLinuxUserModel(user)

	diags = resp.State.Set(linuxCtx.Ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	linuxCtx := util.NewLinuxContext(ctx, r.providerData)

	var plan LinuxUserModel
	diags := req.Plan.Get(linuxCtx.Ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	command := "usermod"

	if plan.Username.IsUnknown() || plan.Username.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Something wrong in username",
			"Something wrong in username",
		)
		return
	}
	username := plan.Username.ValueString()
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

	_, _, commonError := sshUtil.RunCommand(linuxCtx, command, nil)
	if commonError != nil {
		resp.Diagnostics.Append(commonError.Diagnostics...)
		return
	}

	user, commonError := Get(linuxCtx, plan.Username.ValueString())
	if commonError != nil {
		resp.Diagnostics.Append(commonError.Diagnostics...)
		return
	}
	if user == nil {
		resp.Diagnostics.AddError("Failed to created user.", "User not exists after creation request")
		return
	}

	plan = NewLinuxUserModel(user)
	diags = resp.State.Set(linuxCtx.Ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	linuxCtx := util.NewLinuxContext(ctx, r.providerData)

	var state LinuxUserModel
	diags := req.State.Get(linuxCtx.Ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	command := "userdel"
	if state.Username.IsUnknown() || state.Username.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Something wrong in username",
			"Something wrong in username",
		)
		return
	}
	username := state.Username.ValueString()
	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Empty username is not allowed",
			"Empty username is not allowed",
		)
		return
	}

	command = command + " " + username
	_, _, commonError := sshUtil.RunCommand(linuxCtx, command, nil)
	if commonError != nil {
		resp.Diagnostics.Append(commonError.Diagnostics...)
		return
	}
}

func (r *userResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	providerData, commonError := util.ConvertProviderData(req.ProviderData)
	if providerData == nil && commonError == nil {
		return
	}
	if commonError != nil {
		resp.Diagnostics.Append(commonError.Diagnostics...)
		return
	}

	r.providerData = providerData
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("username"), req, resp)
}
