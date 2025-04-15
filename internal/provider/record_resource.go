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
var _ resource.Resource = &RecordResource{}
var _ resource.ResourceWithImportState = &RecordResource{}

func NewRecordResource() resource.Resource {
	return &RecordResource{}
}

// RecordResource defines the resource implementation.
type RecordResource struct {
	client *client.Client
}

// RecordResourceModel describes the resource data model.
type RecordResourceModel struct {
	ID          types.Int64 `tfsdk:"id"`
	
	Name        types.String `tfsdk:"name"`
	ZoneID      types.Int64 `tfsdk:"zone_id"`
	Type        types.String `tfsdk:"type"`
	Value       types.String `tfsdk:"value"`

	Status      types.String `tfsdk:"status"`

	Description types.String `tfsdk:"description"`
	Comments    types.String `tfsdk:"comments"`
	TTL         types.Int64 `tfsdk:"ttl"`

}

func (r *RecordResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_record"
}

func (r *RecordResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Record resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Record id in NetBox",
				PlanModifiers: []planmodifier.Int64{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "DNS Record name",
				Required:            true,
			},
			"zone_id": schema.Int64Attribute{
				MarkdownDescription: "DNS Zone id",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "DNS Record type (A, CNAME, etc.)",
				Required:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "DNS Record value",
				Required:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Record status (active or inactive)",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						string(client.RecordStatusActive),
						string(client.RecordStatusInactive),
					),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Record description",
				Optional:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Record description",
				Optional:            true,
			},
			"ttl": schema.Int64Attribute{
				MarkdownDescription: "Record TTL",
				Optional:            true,
			},

		},
	}
}

func (r *RecordResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = configureResourceClient(req, resp)
}

func (r *RecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RecordResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := data.ToAPIModel(ctx, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	httpRes, err := r.client.PluginsNetboxDnsRecordsCreate(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to create record: %s", err))
		return
	}
	res, err := client.ParsePluginsNetboxDnsRecordsCreateResponse(httpRes)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to parse record response: %s", err))
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

func (r *RecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RecordResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() {
		resp.Diagnostics.AddAttributeError(path.Root("id"), "Internal Error", "Missing ID value")
		return
	}

	httpRes, err := r.client.PluginsNetboxDnsRecordsRetrieve(ctx, int(data.ID.ValueInt64()))
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

func (r *RecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RecordResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	params := data.ToAPIModel(ctx, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	httpRes, err := r.client.PluginsNetboxDnsRecordsUpdate(ctx, int(data.ID.ValueInt64()), params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to update record: %s", err))
		return
	}
	res, err := client.ParsePluginsNetboxDnsRecordsUpdateResponse(httpRes)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to parse record response: %s", err))
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

func (r *RecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RecordResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpRes, err := r.client.PluginsNetboxDnsRecordsDestroy(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to destroy record: %s", err))
		return
	}
	res, err := client.ParsePluginsNetboxDnsRecordsDestroyResponse(httpRes)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to parse response: %s", err))
		return
	}
	if res.StatusCode() != http.StatusNoContent {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to destroy record: %s", string(res.Body)))
		return
	}
}

func (r *RecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importByInt64ID(ctx, req, resp)
}
