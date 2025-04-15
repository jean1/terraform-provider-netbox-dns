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
var _ datasource.DataSource = &ZoneDataSource{}

func NewZoneDataSource() datasource.DataSource {
	return &ZoneDataSource{}
}

type ZoneDataSource struct {
	client *client.Client
}

type ZoneDataSourceModel struct {
	ID            types.Int64        `tfsdk:"id"`
	View          *NestedView        `tfsdk:"view"`
	Name          types.String       `tfsdk:"name"`
	Status        types.String       `tfsdk:"status"`
	Nameservers   []types.Int64      `tfsdk:"nameserver_ids"`
	DefaultTTL    types.Int64        `tfsdk:"default_ttl"`
	SOATTL        types.Int64        `tfsdk:"soa_ttl"`
	SOAMName      *NestedNameserver  `tfsdk:"soa_mname"`
	SOARName      types.String       `tfsdk:"soa_rname"`
	SOASerial     types.Int64        `tfsdk:"soa_serial"`
	SOARefresh    types.Int64        `tfsdk:"soa_refresh"`
	SOARetry      types.Int64        `tfsdk:"soa_retry"`
	SOAExpire     types.Int64        `tfsdk:"soa_expire"`
	SOASerialAuto types.String       `tfsdk:"soa_serial_auto"`
	Description   types.String       `tfsdk:"description"`
}

func (m *ZoneDataSourceModel) FillFromAPIModel(ctx context.Context, resp *client.Zone, diags diag.Diagnostics) {
	m.ID = maybeInt64Value(resp.Id)
	m.View = NestedViewFromAPI(&resp.View)
	m.Name = maybeStringValue(resp.Name)
	m.Status = maybeStringValue(resp.Status)
	for _, id := range *resp.Nameservers {
		m.Nameservers = append(m.Nameservers, types.Int64Value(int64(id)))
	}
	m.DefaultTTL  = maybeInt64Value(resp.DefaultTtl)
	m.SOATTL  = maybeInt64Value(resp.SoaTtl)
	m.SOAMName = NestedNameserver(&resp.SoaMname)
	m.SOARName = maybeStringValue(resp.SoaRname)
	m.SOASerial  = maybeInt64Value(resp.SoaSerial)
	m.SOARefresh  = maybeInt64Value(resp.SoaRefresh)
	m.SOARetry  = maybeInt64Value(resp.SoaRetry)
	m.SOAExpire  = maybeInt64Value(resp.SoaExpire)
	m.SOASerialAuto = maybeBoolValue(resp.SoaSerialAuto)
	m.Description = maybeStringValue(resp.Description)
}

func (d *ZoneDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zone"
}

var zoneDataSchema = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		MarkdownDescription: "ID of the resource in Netbox to use for lookup",
		Required:            true,
	},
	"view": schema.SingleNestedAttribute{
		MarkdownDescription: "DNS View",
		Computed:   true,
		Attributes: (*NestedView)(nil).SchemaAttributes(),
	},
	"name": schema.StringAttribute{
		MarkdownDescription: "zone name",
		Computed: true,
	},
	"status": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: `"active" or "inactive"`,
	},
        "nameserver_ids": schema.ListAttribute{
                ElementType: types.Int64Type,
                Computed:    true,
		MarkdownDescription: `List of nameservers`,
        },
	"defaul_ttl": schema.Int64Attribute{
		Computed: true,
		MarkdownDescription: `Default TTL`,
	},
	"soa_ttl": schema.Int64Attribute{
		Computed: true,
		MarkdownDescription: `Default TTL`,
	},
	"soa_mname": schema.SingleNestedAttribute{
		MarkdownDescription: "Primary nameserver",
		Computed:   true,
		Attributes: (*NestedNameserver)(nil).SchemaAttributes(),
	},
	"soa_rname": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: `zone administrator email address`,
	},

	"soa_serial": schema.Int64Attribute{
		Computed: true,
		MarkdownDescription: `Default TTL`,
	},
	"soa_refresh": schema.Int64Attribute{
		Computed: true,
		MarkdownDescription: `Default TTL`,
	},
	"soa_retry": schema.Int64Attribute{
		Computed: true,
		MarkdownDescription: `Default TTL`,
	},
	"soa_expire": schema.Int64Attribute{
		Computed: true,
		MarkdownDescription: `Default TTL`,
	},
	"soa_serial_auto": schema.BoolAttribute{
		Computed: true,
		MarkdownDescription: `True if serial is generated automaticaly`,
	},
	"description": schema.StringAttribute{
		Computed: true,
		MarkdownDescription: `Zone description`,
	},
}

func (d *ZoneDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DNS Zone data source",
		Attributes:          zoneDataSchema,
	}
}

func (d *ZoneDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = configureDataSourceClient(req, resp)
}

func (d *ZoneDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ZoneDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	httpRes, err := d.client.PluginsNetboxDnsZonesRetrieve(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to retrieve zone: %s", err))
		return
	}
	res, err := client.ParsePluginsNetboxDnsZonesRetrieveResponse(httpRes)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to parse zone: %s", err))
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
