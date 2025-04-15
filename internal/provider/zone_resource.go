// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jean1/terraform-provider-netbox-dns/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ZoneResource{}
var _ resource.ResourceWithImportState = &ZoneResource{}

func NewZoneResource() resource.Resource {
	return &ZoneResource{}
}

// ZoneResource defines the resource implementation.
type ZoneResource struct {
	client *client.Client
}

// ZoneResourceModel describes the resource data model.
type ZoneResourceModel struct {
	ID             types.Int64        `tfsdk:"id"`
	ViewID         types.Int64        `tfsdk:"view"`
	Name           types.String       `tfsdk:"name"`
	Status         types.String       `tfsdk:"status"`
	Nameservers    types.List         `tfsdk:"nameservers"`
	DefaultTTL     types.Int32        `tfsdk:"default_ttl"`
	SOATTL         types.Int32        `tfsdk:"soa_ttl"`
	SOAMNameID     types.Int32        `tfsdk:"soa_mname"`
	SOARName       types.String       `tfsdk:"soa_rname"`
	SOASerial      types.Int32        `tfsdk:"soa_serial"`
	SOAMinimum     types.Int32        `tfsdk:"soa_minimum"`
	SOARefresh     types.Int32        `tfsdk:"soa_refresh"`
	SOARetry       types.Int32        `tfsdk:"soa_retry"`
	SOAExpire      types.Int32        `tfsdk:"soa_expire"`
	SOASerialAuto  types.Bool       `tfsdk:"soa_serial_auto"`
	Description    types.String       `tfsdk:"description"`
}

// Write to API
func (m *ZoneResourceModel) ToAPIModel(ctx context.Context, diags diag.Diagnostics) client.WritableZoneRequest {
	p := client.WritableZoneRequest{}
	p.View = fromInt64Value(m.ViewID)
	p.Name = m.Name.ValueString()
	if !m.Status.IsNull() {
		zonestatus := client.WritableZoneRequestStatus(m.Status.ValueString())
		p.Status = &zonestatus
	}
        if !m.Nameservers.IsNull() {
		var nameservers []BriefNameServerRequest{}
		for _, element :=  range m.Nameservers {
			nameservers = append(nameservers, BriefNameServerRequest {Name:element} )
			diags.Append(diag.WithPath(path.Root("Nameservers"), d))
		}
		p.Nameservers = &nameservers
        }
	p.DefaultTtl = fromInt32Value(m.DefaultTTL)
	p.SoaTtl = fromInt32Value(m.SOATTL)
	p.SoaExpire = fromInt32Value(m.SOAExpire)
	p.SoaMinimum  = fromInt32Value(m.SOAMinimum)
	p.SoaRefresh = fromInt32Value(m.SOARefresh)
	p.SoaRetry = fromInt32Value(m.SOARetry)
	p.SoaMname = fromInt64Value(m.SOAMNameID)
	p.SoaRname = m.SOARName.ValueStringPointer()
	p.SoaSerial = fromInt32Value(m.SOASerial)
	p.SoaSerialAuto = fromBoolValue(m.SOASerialAuto)
	p.Description = m.Description.ValueStringPointer()

	return p
}

// Read from API to resource model
func (m *ZoneResourceModel) FillFromAPIModel(ctx context.Context, resp *client.Zone, diags diag.Diagnostics) {
        m.ID = maybeInt64Value(resp.Id)
	m.ViewID = maybeInt64Value(resp.View.id)
	m.Name = maybeStringValue(resp.Name)
	m.Status = maybeStringValue(resp.Status)
	if resp.Nameservers != nil && len(*resp.Nameservers) > 0 {
                var ds diag.Diagnostics
		var nameservers []string{}
		// api p.Nameservers is a []BriefNameServer
		for _, element := range p.Nameservers {
                	nameservers = append(nameservers, element.Name)
		}
		m.Nameservers = nameservers
                for _, d := range ds {
                        diags.Append(diag.WithPath(path.Root("nameserver_ids"), d))
                }
        }
	m.DefaultTTL = maybeInt32Value(resp.DefaultTtl)
	m.SOATTL = maybeInt32Value(resp.SoaTtl)
	m.SOAMNameID = maybeInt32Value(resp.SoaMname.id)
	m.SOARName = maybeStringValue(resp.SoaRname)
	m.SOASerial = maybeInt32Value(resp.SoaSerial)
	m.SOAMinimum = maybeInt32Value(resp.SoaMinimum)
	m.SOARefresh = maybeInt32Value(resp.SoaRefresh)
	m.SOARetry = maybeInt32Value(resp.SoaRetry)
	m.SOAExpire = maybeInt32Value(resp.DefaultTtl)
	m.SOASerialAuto = maybeBoolValue(resp.SoaSerialAuto)
	m.Description = maybeStringValue(resp.Description)
}

func (r *ZoneResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zone"
}

func (r *ZoneResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DNS Zone resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Zone id in NetBox",
				PlanModifiers: []planmodifier.Int64{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"view_id": schema.Int64Attribute{
				MarkdownDescription: "ID of the DNS View the zone belongs to",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Zone name",
				Required:            true,
			},
			"status": schema.StringAttribute{
                                Required: true,
                                Validators: []validator.String{
                                        stringvalidator.OneOf(
                                                string(client.ZoneStatusActive),
                                                string(client.ZoneStatusDeprecated),
                                                string(client.ZoneStatusDynamic),
                                                string(client.ZoneStatusEmpty),
                                                string(client.ZoneStatusParked),
                                                string(client.ZoneStatusReserved),
                                        ),
                                },
				MarkdownDescription: `one of "active", "deprecated", "dynamic", "empty", "parked" or "reserved"`,

			},
        		"nameserver_ids": schema.ListAttribute{
                                Required: true,
				MarkdownDescription: `List of nameservers`,
				ElementType: types.Int64Type,
			},
			"defaul_ttl": schema.Int32Attribute{
                                Optional: true,
				MarkdownDescription: `Default TTL`,
			},
			"soa_ttl": schema.Int32Attribute{
                                Required: true,
				MarkdownDescription: `Default TTL`,
			},
			"soa_mname": schema.SingleNestedAttribute{
                                Required: true,
				MarkdownDescription: "Primary nameserver",
				Attributes: (*NestedNameserver)(nil).SchemaAttributes(),
			},
			"soa_rname": schema.StringAttribute{
                                Required: true,
				MarkdownDescription: `zone administrator email address`,
			},
			"soa_serial": schema.Int32Attribute{
                                Optional: true,
				MarkdownDescription: `Default TTL`,
			},
			"soa_refresh": schema.Int32Attribute{
                                Required: true,
				MarkdownDescription: `Default TTL`,
			},
			"soa_retry": schema.Int32Attribute{
                                Required: true,
				MarkdownDescription: `Default TTL`,
			},
			"soa_expire": schema.Int32Attribute{
                                Required: true,
				MarkdownDescription: `Default TTL`,
			},
			"soa_serial_auto": schema.BoolAttribute{
				MarkdownDescription: `True if serial is generated automaticaly`,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Zone description",
				Optional:            true,
			},
		},
	}
}

func (r *ZoneResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = configureResourceClient(req, resp)
}

func (r *ZoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ZoneResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := data.ToAPIModel(ctx, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	httpRes, err := r.client.PluginsNetboxDnsZonesCreate(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to create zone: %s", err))
		return
	}
	res, err := client.ParsePluginsNetboxDnsZonesCreateResponse(httpRes)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to parse zone response: %s", err))
		return
	}
	if res.JSON201 == nil {
		resp.Diagnostics.AddError("Client Error", httpError(httpRes, res.Body))
		return
	}

	data.FillFromAPIModel(ctx, res.JSON201, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ZoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ZoneResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() {
		resp.Diagnostics.AddAttributeError(path.Root("id"), "Internal Error", "Missing ID value")
		return
	}

	httpRes, err := r.client.PluginsNetboxDnsZonesRetrieve(ctx, int(data.ID.ValueInt64()))
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

func (r *ZoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ZoneResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	params := data.ToAPIModel(ctx, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	httpRes, err := r.client.PluginsNetboxDnsZonesUpdate(ctx, int(data.ID.ValueInt64()), params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to update zone: %s", err))
		return
	}
	res, err := client.ParsePluginsNetboxDnsZonesUpdateResponse(httpRes)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to parse zone response: %s", err))
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

func (r *ZoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ZoneResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpRes, err := r.client.PluginsNetboxDnsZonesDestroy(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to destroy zone: %s", err))
		return
	}
	res, err := client.ParsePluginsNetboxDnsZonesDestroyResponse(httpRes)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to parse response: %s", err))
		return
	}
	if res.StatusCode() != http.StatusNoContent {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to destroy zone: %s", string(res.Body)))
		return
	}
}

func (r *ZoneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importByInt64ID(ctx, req, resp)
}
