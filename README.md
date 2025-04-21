# Netbox DNS Terraform Provider

_This template repository is built on the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework). The template repository built on the [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk) can be found at [terraform-provider-scaffolding](https://github.com/hashicorp/terraform-provider-scaffolding). See [Which SDK Should I Use?](https://developer.hashicorp.com/terraform/plugin/framework-benefits) in the Terraform documentation for additional information._

## Configure

Example configuration:

```tf
provider "netboxdns" {
  server_url = "https://netbox.example.com"
  api_token  = var.netbox_api_token
}
```
You can also set the provider config from environment variables:

- `NETBOX_SERVER_URL` in place of `server_url`
- `NETBOX_API_TOKEN` in place of `api_token`

For more details and additional properties, see [the docs](./docs/index.md).

## Development

### Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.22
- [Docker](https://docs.docker.com/desktop/)

### Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```
