package directory

import (
	"context"
	"terraform-provider-linux/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/melbahja/goph"
)

var (
	_ datasource.DataSource              = &directoryDataSource{}
	_ datasource.DataSourceWithConfigure = &directoryDataSource{}
)

func NewDirectoryDataSource() datasource.DataSource {
	return &directoryDataSource{}
}

type directoryDataSource struct {
	client *goph.Client
}

// Configure implements datasource.DataSourceWithConfigure.
func (*directoryDataSource) Configure(context.Context, datasource.ConfigureRequest, *datasource.ConfigureResponse) {
	panic("unimplemented")
}

// Metadata implements datasource.DataSource.
func (*directoryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_directory"
}

// Read implements datasource.DataSource.
func (d *directoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	linuxCtx := util.NewLinuxContext(ctx, d.client)

	var state LinuxDirectoryModel

	diags := req.Config.Get(linuxCtx.Ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if state.Path.IsUnknown() || state.Path.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("path"),
			"Wrong path",
			"Wrong path",
		)
		return
	}

	directory_path := state.Path.ValueString()
	if directory_path == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("path"),
			"Empty path is not allowed",
			"Please speicfy path",
		)
		return
	}

	directory, commonError := Get(linuxCtx, directory_path)
	if commonError != nil {
		resp.Diagnostics.Append(commonError.Diagnostics...)
		return
	}

	state = NewLinuxDirectoryModel(directory)

	diags = resp.State.Set(linuxCtx.Ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Schema implements datasource.DataSource.
func (*directoryDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Required: true,
			},
		},
	}
}
