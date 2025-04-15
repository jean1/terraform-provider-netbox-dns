package provider

import (
	"context"
	"fmt"

	"github.com/jean1/terraform-provider-netbox-dns/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ViewDataSource{}

func NewViewDataSource() datasource.DataSource {
	return &ViewDataSource{}
}

type ViewDataSource struct {
	client *client.Client
}

type ViewDataSourceModel struct {
	ID            types.Int64        `tfsdk:"id"`
	Name          types.String       `tfsdk:"name"`
	Description   types.String       `tfsdk:"description"`
}

func (m *ViewDataSourceModel) FillFromAPIModel(ctx context.Context, resp *client.View, diags diag.Diagnostics) {
	m.ID = maybeInt64Value(resp.Id)
	m.Name = maybeStringValue(&resp.Name)
	m.Description = maybeStringValue(resp.Description)
}

func (d *ViewDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_view"
}

var viewDataSchema = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		MarkdownDescription: "ID of the resource in Netbox to use for lookup",
		Required:            true,
	},
	"name": schema.StringAttribute{
		MarkdownDescription: "view name",
		Computed: true,
	},
	"description": schema.StringAttribute{
		Computed: true,
		MarkdownDescription: `View description`,
	},
}

func (d *ViewDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DNS View data source",
		Attributes:          viewDataSchema,
	}
}

func (d *ViewDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = configureDataSourceClient(req, resp)
}

func (d *ViewDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ViewDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	httpRes, err := d.client.PluginsNetboxDnsViewsRetrieve(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to retrieve view: %s", err))
		return
	}
	res, err := client.ParsePluginsNetboxDnsViewsRetrieveResponse(httpRes)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to parse view: %s", err))
		return
	}
	if res.JSON200 == nil {
		resp.Diagnostics.AddError("Client Error", httpError(httpRes, res.Body))
		return
	}

	data.FillFromAPIModel(ctx, res.JSON200, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
