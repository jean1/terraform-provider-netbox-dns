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
			"view": schema.SingleNestedAttribute{
				MarkdownDescription: "DNS View",
				Required:            true,
				Attributes: (*NestedView)(nil).SchemaAttributes(),
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Zone name",
				Required:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: `"active" or "inactive"`,
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
                                MarkdownDescription: `One of: "active", "reserved", "deprecated", "parked", "dynamic"`,

			},
        		"nameserver_ids": schema.ListAttribute{
                                Required: true,
				MarkdownDescription: `List of nameservers`,
				ElementType: types.Int64Type,
			},
			"defaul_ttl": schema.Int64Attribute{
                                Optional: true,
				MarkdownDescription: `Default TTL`,
			},
			"soa_ttl": schema.Int64Attribute{
                                Required: true,
				MarkdownDescription: `Default TTL`,
			},
			"soa_mname": schema.SingleNestedAttribute{
                                Required: true,
				MarkdownDescription: "Primary nameserver",
				Attributes: (*NestedNameserver)(nil).SchemaAttributes(),
			}
			"soa_rname": schema.StringAttribute{
                                Required: true,
				MarkdownDescription: `zone administrator email address`,
			},
			"soa_serial": schema.Int64Attribute{
                                Optional: true,
				MarkdownDescription: `Default TTL`,
			},
			"soa_refresh": schema.Int64Attribute{
                                Required: true,
				MarkdownDescription: `Default TTL`,
			},
			"soa_retry": schema.Int64Attribute{
                                Required: true,
				MarkdownDescription: `Default TTL`,
			},
			"soa_expire": schema.Int64Attribute{
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
