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
var _ datasource.DataSource = &RecordDataSource{}

func NewRecordDataSource() datasource.DataSource {
	return &RecordDataSource{}
}

type RecordDataSource struct {
	client *client.Client
}

type RecordDataSourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Zone        *NestedZone  `tfsdk:"zone"`
	Type        types.String `tfsdk:"type"`
	Value       types.String `tfsdk:"value"`
	Status      types.String `tfsdk:"status"`
	Description types.String `tfsdk:"description"`
	Comments    types.String `tfsdk:"comments"`
	TTL         types.Int64 `tfsdk:"ttl"`
}

func (m *RecordDataSourceModel) FillFromAPIModel(ctx context.Context, resp *client.Record, diags diag.Diagnostics) {
	m.ID = maybeInt64Value(resp.Id)
	m.Name = maybeStringValue(resp.Name)
	m.Zone = NestedZoneFromAPI(resp.Zone)
	m.Type = maybeStringValue(resp.Type)
	m.Value = maybeStringValue(resp.Value)
	m.Status = maybeStringValue(resp.Status)
	m.Description = maybeStringValue(resp.Description)
	m.Comments = maybeStringValue(resp.Comments)
	m.TTL = maybeInt64Value(resp.TTL)
}

func (d *RecordDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_record"
}

var recordDataSchema = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		MarkdownDescription: "ID of the resource in Netbox to use for lookup",
		Required:            true,
	},
	"name": schema.StringAttribute{
		Computed: true,
	},
	"zone": schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: (*NestedZone)(nil).SchemaAttributes(),
	},
	"type": schema.StringAttribute{
		Computed: true,
	},
	"value": schema.StringAttribute{
		Computed: true,
	},
	"status": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: `"active" or "inactive"`,
	},
	"description": schema.StringAttribute{
		Computed: true,
	},
	"comments": schema.StringAttribute{
		Computed: true,
	},
	"ttl": schema.Int64Attribute{
		Computed: true,
	},
}

func (d *RecordDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DNS Record data source",
		Attributes:          recordDataSchema,
	}
}

func (d *RecordDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = configureDataSourceClient(req, resp)
}

func (d *RecordDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RecordDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	httpRes, err := d.client.PluginsNetboxDnsRecordsRetrieve(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to retrieve record: %s", err))
		return
	}
	res, err := client.ParsePluginsNetboxDnsRecordsRetrieveResponse(httpRes)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to parse record: %s", err))
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
