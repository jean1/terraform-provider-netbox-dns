Configuration:

```terraform
provider "netboxbgp" {
  server_url = "https://netbox.my-company.net"
  api_token  = var.netbox_api_token
}
```

Example usage of the Record resource:

```terraform
resource "netboxdns_record" "test" {
        name    = "www"
        zone    = "example.com"
        view    = "_default_"
        type    = "A"
        value   = "192.0.2.1"
        status  = "active"
}
```
