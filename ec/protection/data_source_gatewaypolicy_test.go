package protection_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
)

func TestDataSourceGatewayPolicy(t *testing.T) {
	name := "my-gateway-policy"
	pf, _ := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGatewayPolicyConfigBasic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_protection_gatewaypolicy.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_protection_gatewaypolicy.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_protection_gatewaypolicy.test", "spec.0.description", "My Gateway Policy"),
					resource.TestCheckResourceAttr("ec_protection_gatewaypolicy.test", "spec.0.destination_cidrs.#", "1"),
					resource.TestCheckResourceAttr("ec_protection_gatewaypolicy.test", "spec.0.destination_cidrs.0", "1.2.3.4/32"),
				),
			},
			{
				Config: testDataSourceGatewayPolicyConfigBasic(name) +
					testDataSourceGatewayPolicyConfigRead(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ec_protection_gatewaypolicy.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("data.ec_protection_gatewaypolicy.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("data.ec_protection_gatewaypolicy.test", "spec.0.description", "My Gateway Policy"),
					resource.TestCheckResourceAttr("data.ec_protection_gatewaypolicy.test", "spec.0.destination_cidrs.#", "1"),
					resource.TestCheckResourceAttr("data.ec_protection_gatewaypolicy.test", "spec.0.destination_cidrs.0", "1.2.3.4/32"),
				),
			},
		},
	})
}

func testDataSourceGatewayPolicyConfigBasic(name string) string {
	return fmt.Sprintf(`resource "ec_protection_gatewaypolicy" "test" {
  metadata {
    name = "%s"
  }
  spec {
    description = "My Gateway Policy"
    destination_cidrs = ["1.2.3.4/32"]
  }
}
`, name)
}

func testDataSourceGatewayPolicyConfigRead() string {
	return `data "ec_protection_gatewaypolicy" "test" {
  metadata {
    name      = "${ec_protection_gatewaypolicy.test.metadata.0.name}"
  }
}
`
}
