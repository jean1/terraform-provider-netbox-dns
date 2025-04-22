package provider

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandIPv4() string {
	buf := make([]byte, 4)
	ip := rand.Uint32()
	binary.LittleEndian.PutUint32(buf, ip)
	return net.IP(buf).String()
}

func TestAccRecordResource(t *testing.T) {
	// var ipaddress = "192.0.2.1"
	var hostname = RandString(8)
	var newhostname = RandString(8)
	var ipaddress = RandIPv4()
	var ipv6address = "2001:db6::" + RandIPv4()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(`
                                        resource "netboxdns_record" "test" {
                                                name    = "%s"
                                                zone    = "u-strasbg.fr"
                                                view    = "_default_"
                                                type    = "A"
                                                value   = "%s"
                                                status  = "active"
                                        }`, hostname, ipaddress),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netboxdns_record.test", "id"),
					resource.TestCheckResourceAttr("netboxdns_record.test", "name", hostname),
					resource.TestCheckResourceAttr("netboxdns_record.test", "value", ipaddress),
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
						name              = "%s"
                                                zone    = "u-strasbg.fr"
                                                view    = "_default_"
                                                type    = "AAAA"
                                                value   = "%s"
                                                status  = "active"
					}
				`, newhostname, ipv6address),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netboxdns_record.test", "name", newhostname),
					resource.TestCheckResourceAttr("netboxdns_record.test", "type", "AAAA"),
					resource.TestCheckResourceAttr("netboxdns_record.test", "value", ipv6address),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
