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
var _ resource.Resource = &NameserverResource{}
var _ resource.ResourceWithImportState = &NameserverResource{}

func NewNameserverResource() resource.Resource {
	return &NameserverResource{}
}

// NameserverResource defines the resource implementation.
type NameserverResource struct {
	client *client.Client
}

// NameserverResourceModel describes the resource data model.
type NameserverResourceModel struct {
	ID          types.Int64 `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func (r *NameserverResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nameserver"
}

func (r *NameserverResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Nameserver resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Nameserver id in NetBox",
				PlanModifiers: []planmodifier.Int64{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "DNS Nameserver name",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Nameserver description",
				Optional:            true,
			},
		},
	}
}

func (r *NameserverResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = configureResourceClient(req, resp)
}

func (r *NameserverResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NameserverResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := data.ToAPIModel(ctx, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	httpRes, err := r.client.PluginsNetboxDnsNameserversCreate(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to create nameserver: %s", err))
		return
	}
	res, err := client.ParsePluginsNetboxDnsNameserversCreateResponse(httpRes)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to parse nameserver response: %s", err))
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

func (r *NameserverResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NameserverResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() {
		resp.Diagnostics.AddAttributeError(path.Root("id"), "Internal Error", "Missing ID value")
		return
	}

	httpRes, err := r.client.PluginsNetboxDnsNameserversRetrieve(ctx, int(data.ID.ValueInt64()))
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

func (r *NameserverResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data NameserverResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	params := data.ToAPIModel(ctx, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	httpRes, err := r.client.PluginsNetboxDnsNameserversUpdate(ctx, int(data.ID.ValueInt64()), params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to update nameserver: %s", err))
		return
	}
	res, err := client.ParsePluginsNetboxDnsNameserversUpdateResponse(httpRes)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to parse nameserver response: %s", err))
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

func (r *NameserverResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NameserverResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpRes, err := r.client.PluginsNetboxDnsNameserversDestroy(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to destroy nameserver: %s", err))
		return
	}
	res, err := client.ParsePluginsNetboxDnsNameserversDestroyResponse(httpRes)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to parse response: %s", err))
		return
	}
	if res.StatusCode() != http.StatusNoContent {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to destroy nameserver: %s", string(res.Body)))
		return
	}
}

func (r *NameserverResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importByInt64ID(ctx, req, resp)
}
