// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRecordResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(`
                                        resource "netboxdns_record" "test" {
                                                name    = "www"
                                                zone_id = 1
                                                type    = "A"
                                                value   = "192.0.2.1"
                                                status  = "active"
                                        }`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netboxdns_record.test", "id"),
					resource.TestCheckResourceAttr("netboxdns_record.test", "name", "www"),
					resource.TestCheckResourceAttr("netboxdns_record.test", "value", "192.0.2.1"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "netboxdns_record.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: fmt.Sprintf(`
					resource "netboxdns_record" "test" {
						name              = "newname"
                                                zone_id = 1
                                                type    = "AAAA"
                                                value   = "2001:db6::1"
                                                status  = "active"
					}
				`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netboxdns_record.test", "name", "newname"),
					resource.TestCheckResourceAttr("netboxdns_record.test", "type", "AAAA"),
					resource.TestCheckResourceAttr("netboxdns_record.test", "value", "2001:db6::1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
