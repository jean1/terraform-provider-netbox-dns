// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"crypto/tls"

	"net/http"

	"github.com/jean1/terraform-provider-netbox-dns/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sethvargo/go-envconfig"
)

// Ensure NetboxDNSProvider satisfies various provider interfaces.
var _ provider.Provider = &NetboxDNSProvider{}






// NetboxDNSProvider defines the provider implementation.
type NetboxDNSProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// NetboxDNSProviderModel describes the provider data model.
type NetboxDNSProviderModel struct {
	ServerURL          types.String `tfsdk:"server_url"`
	APIToken           types.String `tfsdk:"api_token"`
	AllowInsecureHTTPS types.Bool   `tfsdk:"allow_insecure_https"`
	Headers            types.Map    `tfsdk:"headers"`
	RequestTimeout     types.Int64  `tfsdk:"request_timeout"`
}

type NetboxDNSProviderEnvModel struct {
	ServerURL          string `env:"NETBOX_SERVER_URL"`
	APIToken           string `env:"NETBOX_API_TOKEN"`
	AllowInsecureHTTPS *bool  `env:"NETBOX_ALLOW_INSECURE_HTTPS"`
	RequestTimeout     int64  `env:"NETBOX_REQUEST_TIMEOUT"`
}

func (p *NetboxDNSProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "netboxdns"
	resp.Version = p.version
}

func (p *NetboxDNSProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: providerDocs,
		Attributes: map[string]schema.Attribute{
			"server_url": schema.StringAttribute{
				MarkdownDescription: "Location of Netbox server including scheme (http or https) and optional port. Can be set via the `NETBOX_SERVER_URL` environment variable.",
				Optional:            true,
			},
			"api_token": schema.StringAttribute{
				MarkdownDescription: "Netbox API authentication token. Can be set via the `NETBOX_API_TOKEN` environment variable.",
				Optional:            true,
			},
			"allow_insecure_https": schema.BoolAttribute{
				MarkdownDescription: "Flag to set whether to allow https with invalid certificates. Can be set via the `NETBOX_ALLOW_INSECURE_HTTPS` environment variable. Defaults to `false`.",
				Optional:            true,
			},
			"headers": schema.MapAttribute{
				MarkdownDescription: "Set these header on all requests to Netbox. Can be set via the `NETBOX_HEADERS` environment variable.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"request_timeout": schema.Int64Attribute{
				MarkdownDescription: "Netbox API HTTP request timeout in seconds. Can be set via the `NETBOX_REQUEST_TIMEOUT` environment variable.",
				Optional:            true,
			},
		},
	}
}

func apiKeyAuth(token string) client.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Token "+token)
		return nil
	}
}

type configuredProvider struct {
	Client *client.Client
}

func (p *NetboxDNSProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data NetboxDNSProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// load defaults from env vars
	envData := NetboxDNSProviderEnvModel{}
	err := envconfig.ProcessWith(ctx, &envconfig.Config{
		Target:           &envData,
		DefaultNoInit:    true,
		DefaultOverwrite: true,
	})
	if err != nil {
		resp.Diagnostics.AddError("failed to parse environment variables", err.Error())
		return
	}
	if data.ServerURL.IsNull() && envData.ServerURL != "" {
		data.ServerURL = types.StringValue(envData.ServerURL)
	}
	if data.APIToken.IsNull() && envData.APIToken != "" {
		data.APIToken = types.StringValue(envData.APIToken)
	}
	if data.AllowInsecureHTTPS.IsNull() && envData.AllowInsecureHTTPS != nil {
		data.AllowInsecureHTTPS = types.BoolValue(*envData.AllowInsecureHTTPS)
	}
	if data.RequestTimeout.IsNull() && envData.RequestTimeout > 0 {
		data.RequestTimeout = types.Int64Value(envData.RequestTimeout)
	}

	// apply defaults
	if data.RequestTimeout.IsNull() {
		data.RequestTimeout = types.Int64Value(10)
	}

	if data.ServerURL.IsNull() {
		resp.Diagnostics.AddError("Missing required attribute", "Server URL is required")
	}
	if data.APIToken.IsNull() {
		resp.Diagnostics.AddError("Missing required attribute", "API token is required")
	}
	if resp.Diagnostics.HasError() {
		return
	}

	opts := []client.ClientOption{
		client.WithRequestEditorFn(apiKeyAuth(data.APIToken.ValueString())), // auth
	}
	if !data.AllowInsecureHTTPS.IsNull() && data.AllowInsecureHTTPS.ValueBool() {
		httpClient := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
		opts = append(opts, client.WithHTTPClient(httpClient))
	}

	client, err := client.NewClient(data.ServerURL.ValueString(), opts...)
	if err != nil {
		resp.Diagnostics.AddError("failed to create client", err.Error())
		return
	}

	providerData := configuredProvider{
		Client: client,
	}
	resp.DataSourceData = &providerData
	resp.ResourceData = &providerData
}

func (p *NetboxDNSProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRecordResource,
	}
}

func (p *NetboxDNSProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewRecordDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &NetboxDNSProvider{
			version: version,
		}
	}
}

func configureResourceClient(req resource.ConfigureRequest, resp *resource.ConfigureResponse) *client.Client {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return nil
	}

	data, ok := req.ProviderData.(*configuredProvider)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *configuredProvider, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return nil
	}

	return data.Client
}

func configureDataSourceClient(req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) *client.Client {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return nil
	}

	data, ok := req.ProviderData.(*configuredProvider)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *configuredProvider, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return nil
	}

	return data.Client
}
