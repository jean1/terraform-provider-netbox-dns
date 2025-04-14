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
        Slug    types.String `tfsdk:"slug"`
        URL     types.String `tfsdk:"url"`
}

func (tfo NestedZone) ToAPIModel() client.NestedZone {
        return client.NestedZone{
                Id:      toIntPointer(tfo.ID.ValueInt64Pointer()),
                Url:     tfo.URL.ValueStringPointer(),
                Display: tfo.Display.ValueStringPointer(),
                Name:    tfo.Name.ValueString(),
                Slug:    tfo.Slug.ValueString(),
        }
}

func NestedZoneFromAPI(resp *client.NestedZone) *NestedZone {
        if resp == nil {
                return nil
        }
        tfo := &NestedZone{}
        tfo.ID = types.Int64Value(int64(*resp.Id))
        tfo.URL = maybeStringValue(resp.Url)
        tfo.Display = maybeStringValue(resp.Display)
        tfo.Name = types.StringValue(resp.Name)
        tfo.Slug = types.StringValue(resp.Slug)
        return tfo
}

func (*NestedZone) SchemaAttributes() map[string]schema.Attribute {
        return map[string]schema.Attribute{
                "id": schema.Int64Attribute{
                        Computed: true,
                },
                "display": schema.StringAttribute{
                        Computed: true,
                },
                "url": schema.StringAttribute{
                        Optional: true,
                },
                "name": schema.StringAttribute{
                        Required: true,
                },
                "slug": schema.StringAttribute{
                        Required: true,
                },
        }
}

