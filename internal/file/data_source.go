package file

import (
	"context"
	"terraform-provider-linux/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

var (
	_ datasource.DataSource              = &fileDataSource{}
	_ datasource.DataSourceWithConfigure = &fileDataSource{}
)

func NewFileDataSource() datasource.DataSource {
	return &fileDataSource{}
}

type fileDataSource struct {
	providerData *util.LinuxProviderData
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *fileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	providerData, commonError := util.ConvertProviderData(req.ProviderData)
	if providerData == nil && commonError == nil {
		return
	}
	if commonError != nil {
		resp.Diagnostics.Append(commonError.Diagnostics...)
		return
	}

	d.providerData = providerData
}

// Metadata implements datasource.DataSource.
func (*fileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

// Read implements datasource.DataSource.
func (d *fileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	linuxCtx := util.NewLinuxContext(ctx, d.providerData)

	var state LinuxFileModel

	diags := req.Config.Get(linuxCtx.Ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if state.Path.IsUnknown() || state.Path.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("path"),
			"Wrong path",
			"Invalid or unknown path provided",
		)
		return
	}

	directory_path := state.Path.ValueString()
	if directory_path == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("path"),
			"Empty path is not allowed",
			"Please specify a valid path",
		)
		return
	}

	directory, commonError := Get(linuxCtx, directory_path)
	if commonError != nil {
		resp.Diagnostics.Append(commonError.Diagnostics...)
		return
	}

	state = NewLinuxFileModel(directory)

	diags = resp.State.Set(linuxCtx.Ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Schema implements datasource.DataSource.
func (*fileDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Required: true,
			},
		},
	}
}
