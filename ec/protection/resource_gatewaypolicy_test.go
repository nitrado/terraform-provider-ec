package protection_test

import (
	"context"
	"fmt"
	"testing"

	metav1 "github.com/gamefabric/gf-apicore/apis/meta/v1"
	"github.com/gamefabric/gf-core/pkg/apiclient/clientset"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
)

func TestResourceGatewayPolicy(t *testing.T) {
	name := "my-gateway-policy"
	pf, cs := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		CheckDestroy:      testCheckGatewayPolicyDestroy(cs),
		Steps: []resource.TestStep{
			{
				Config: testResourceGatewayPolicyConfigBasic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_protection_gatewaypolicy.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_protection_gatewaypolicy.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_protection_gatewaypolicy.test", "spec.0.description", "My Gateway Policy"),
					resource.TestCheckResourceAttr("ec_protection_gatewaypolicy.test", "spec.0.destination_cidrs.#", "1"),
					resource.TestCheckResourceAttr("ec_protection_gatewaypolicy.test", "spec.0.destination_cidrs.0", "1.2.3.4/32"),
				),
			},
			{
				Config: testResourceGatewayPolicyConfigMultipleCIDRs(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_protection_gatewaypolicy.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_protection_gatewaypolicy.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_protection_gatewaypolicy.test", "spec.0.description", "My Gateway Policy"),
					resource.TestCheckResourceAttr("ec_protection_gatewaypolicy.test", "spec.0.destination_cidrs.#", "2"),
					resource.TestCheckResourceAttr("ec_protection_gatewaypolicy.test", "spec.0.destination_cidrs.0", "1.2.3.4/32"),
					resource.TestCheckResourceAttr("ec_protection_gatewaypolicy.test", "spec.0.destination_cidrs.1", "2.3.4.5/32"),
				),
			},
			{
				ResourceName:      "ec_protection_gatewaypolicy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testResourceGatewayPolicyConfigBasic(name string) string {
	return fmt.Sprintf(`resource "ec_protection_gatewaypolicy" "test" {
  metadata {
    name = "%s"
  }
  spec {
    description = "My Gateway Policy"
    destination_cidrs = ["1.2.3.4/32"]
  }
}`, name)
}

func testResourceGatewayPolicyConfigMultipleCIDRs(name string) string {
	return fmt.Sprintf(`resource "ec_protection_gatewaypolicy" "test" {
  metadata {
    name = "%s"
  }
  spec {
    description = "My Gateway Policy"
    destination_cidrs = ["1.2.3.4/32", "2.3.4.5/32"]
  }
}`, name)
}

func testCheckGatewayPolicyDestroy(cs clientset.Interface) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "ec_protection_gatewaypolicy" {
				continue
			}

			name := rs.Primary.ID
			resp, err := cs.ProtectionV1Alpha1().Protocols().Get(context.Background(), name, metav1.GetOptions{})
			if err == nil {
				if resp.Name == rs.Primary.ID {
					return fmt.Errorf("gateway policy still exists: %s", rs.Primary.ID)
				}
			}
		}
		return nil
	}
}
