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
var _ datasource.DataSource = &NameserverDataSource{}

func NewNameserverDataSource() datasource.DataSource {
	return &NameserverDataSource{}
}

type NameserverDataSource struct {
	client *client.Client
}

type NameserverDataSourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	Name        types.String `tfsdk:"name"`
}

func (m *NameserverDataSourceModel) FillFromAPIModel(ctx context.Context, resp *client.NameServer, diags diag.Diagnostics) {
	m.ID = maybeInt64Value(resp.Id)
	m.Description = maybeStringValue(resp.Description)
	m.Name = maybeStringValue(resp.Name)
}

func (d *NameserverDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nameserver"
}

var nameserverDataSchema = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		MarkdownDescription: "ID of the resource in Netbox to use for lookup",
		Required:            true,
	},
	"description": schema.StringAttribute{
		Computed: true,
	},
	"name": schema.StringAttribute{
		Computed: true,
	},
}

func (d *NameserverDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DNS Nameserver data source",
		Attributes:          nameserverDataSchema,
	}
}

func (d *NameserverDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = configureDataSourceClient(req, resp)
}

func (d *NameserverDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data NameserverDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	httpRes, err := d.client.PluginsNetboxDnsNameserversRetrieve(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to retrieve nameserver: %s", err))
		return
	}
	res, err := client.ParsePluginsNetboxDnsNameserversRetrieveResponse(httpRes)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to parse nameserver: %s", err))
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
