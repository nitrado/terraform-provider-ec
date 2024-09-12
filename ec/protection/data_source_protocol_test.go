package protection_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
)

func TestDataSourceProtocols(t *testing.T) {
	name := "my-protocol"
	pf, _ := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceProtocolsConfigBasic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_protection_protocol.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_protection_protocol.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_protection_protocol.test", "spec.0.description", "My Protocol"),
					resource.TestCheckResourceAttr("ec_protection_protocol.test", "spec.0.mitigation_name", "my-mitigation"),
					resource.TestCheckResourceAttr("ec_protection_protocol.test", "spec.0.protocol", "UDP"),
				),
			},
			{
				Config: testDataSourceProtocolsConfigBasic(name) +
					testDataSourceProtocolConfigRead(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ec_protection_protocol.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("data.ec_protection_protocol.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("data.ec_protection_protocol.test", "spec.0.description", "My Protocol"),
					resource.TestCheckResourceAttr("data.ec_protection_protocol.test", "spec.0.mitigation_name", "my-mitigation"),
					resource.TestCheckResourceAttr("data.ec_protection_protocol.test", "spec.0.protocol", "UDP"),
				),
			},
		},
	})
}

func testDataSourceProtocolsConfigBasic(name string) string {
	return fmt.Sprintf(`resource "ec_protection_protocol" "test" {
  metadata {
    name = "%s"
  }
  spec {
    description = "My Protocol"
    mitigation_name = "my-mitigation"
    protocol = "UDP"
  }
}
`, name)
}

func testDataSourceProtocolConfigRead() string {
	return `data "ec_protection_protocol" "test" {
  metadata {
    name      = "${ec_protection_protocol.test.metadata.0.name}"
  }
}
`
}
