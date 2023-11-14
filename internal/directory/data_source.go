package directory

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/melbahja/goph"
)

var (
	_ datasource.DataSource = &directoryDataSource{}
)

func NewDirectoryDataSource() datasource.DataSource {
	return &directoryDataSource{}
}

type directoryDataSource struct {
	client *goph.Client
}

// Metadata implements datasource.DataSource.
func (*directoryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_directory"
}

// Read implements datasource.DataSource.
func (*directoryDataSource) Read(context.Context, datasource.ReadRequest, *datasource.ReadResponse) {
	panic("unimplemented")
}

// Schema implements datasource.DataSource.
func (*directoryDataSource) Schema(context.Context, datasource.SchemaRequest, *datasource.SchemaResponse) {
	panic("unimplemented")
}
