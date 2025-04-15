package provider

import (
        "github.com/jean1/terraform-provider-netbox-dns/client"
        "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
        "github.com/hashicorp/terraform-plugin-framework/types"
)

type NestedZone struct {
        Display types.String `tfsdk:"display"`
        ID      types.Int64  `tfsdk:"id"`
        Name    types.String `tfsdk:"name"`
        URL     types.String `tfsdk:"url"`
}

func (tfo NestedZone) ToAPIModel() client.NestedZone {
        return client.NestedZone{
                Id:      toIntPointer(tfo.ID.ValueInt64Pointer()),
                Url:     tfo.URL.ValueStringPointer(),
                Display: tfo.Display.ValueStringPointer(),
                Name:    tfo.Name.ValueString(),
        }
}

func NestedZoneFromAPI(resp *client.NestedZone) *NestedZone {
        if resp == nil {
                return nil
        }
        tfo := &NestedZone{}
        tfo.ID = types.Int64Value(int64(*resp.Id))
        tfo.Name = types.StringValue(resp.Name)
        tfo.URL = maybeStringValue(resp.Url)
        tfo.Display = maybeStringValue(resp.Display)
        return tfo
}

func (*NestedZone) SchemaAttributes() map[string]schema.Attribute {
        return map[string]schema.Attribute{
                "id": schema.Int64Attribute{
                        Computed: true,
                },
                "name": schema.StringAttribute{
                        Required: true,
                },
                "url": schema.StringAttribute{
                        Optional: true,
                },
                "display": schema.StringAttribute{
                        Computed: true,
                },
        }
}

type NestedView struct {
        ID      types.Int64  `tfsdk:"id"`
        Name    types.String `tfsdk:"name"`
        URL     types.String `tfsdk:"url"`
        Display types.String `tfsdk:"display"`
}

func (tfo NestedView) ToAPIModel() client.BriefView {
        return client.BriefView{
                Id:      toIntPointer(tfo.ID.ValueInt64Pointer()),
                Name:    tfo.Name.ValueString(),
                Url:     tfo.URL.ValueStringPointer(),
                Display: tfo.Display.ValueStringPointer(),
        }
}

func NestedView(resp *client.BriefView) *NestedView {
        if resp == nil {
                return nil
        }
        tfo := &NestedView{}
        tfo.ID = types.Int64Value(int64(*resp.Id))
        tfo.Name = types.StringValue(resp.Name)
        tfo.URL = maybeStringValue(resp.Url)
        tfo.Display = maybeStringValue(resp.Display)
        return tfo
}

func (*NestedView) SchemaAttributes() map[string]schema.Attribute {
        return map[string]schema.Attribute{
                "id": schema.Int64Attribute{
                        Computed: true,
                },
                "name": schema.StringAttribute{
                        Required: true,
                },
                "url": schema.StringAttribute{
                        Optional: true,
                },
                "display": schema.StringAttribute{
                        Computed: true,
                },
        }
}

type NestedNameserver struct {
        ID      types.Int64  `tfsdk:"id"`
        Name    types.String `tfsdk:"name"`
        URL     types.String `tfsdk:"url"`
        Display types.String `tfsdk:"display"`
}

func (tfo NestedNameserver) ToAPIModel() client.BriefNameServer {
        return client.BriefNameServer{
                Id:      toIntPointer(tfo.ID.ValueInt64Pointer()),
                Name:    tfo.Name.ValueString(),
                Url:     tfo.URL.ValueStringPointer(),
                Display: tfo.Display.ValueStringPointer(),
        }
}

func NestedNameserverFromAPI(resp *client.BriefNameServer) *NestedNameserver {
        if resp == nil {
                return nil
        }
        tfo := &NestedNameserver{}
        tfo.ID = types.Int64Value(int64(*resp.Id))
        tfo.Name = types.StringValue(resp.Name)
        tfo.URL = maybeStringValue(resp.Url)
        tfo.Display = maybeStringValue(resp.Display)
        return tfo
}

func (*NestedNameserver) SchemaAttributes() map[string]schema.Attribute {
        return map[string]schema.Attribute{
                "id": schema.Int64Attribute{
                        Computed: true,
                },
                "name": schema.StringAttribute{
                        Required: true,
                },
                "url": schema.StringAttribute{
                        Optional: true,
                },
                "display": schema.StringAttribute{
                        Computed: true,
                },
        }
}
