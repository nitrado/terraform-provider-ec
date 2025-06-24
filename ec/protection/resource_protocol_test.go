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

func TestResourceProtocol(t *testing.T) {
	name := "my-protocol"
	pf, cs := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		CheckDestroy:      testCheckProtocolDestroy(cs),
		Steps: []resource.TestStep{
			{
				Config: testResourceProtocolsConfigBasic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_protection_protocol.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_protection_protocol.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_protection_protocol.test", "spec.0.description", "My Protocol"),
					resource.TestCheckResourceAttr("ec_protection_protocol.test", "spec.0.mitigation_name", "my-mitigation"),
					resource.TestCheckResourceAttr("ec_protection_protocol.test", "spec.0.protocol", "UDP"),
				),
			},
			{
				Config: testResourceProtocolsConfigWithTCP(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_protection_protocol.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_protection_protocol.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_protection_protocol.test", "spec.0.description", "My Protocol"),
					resource.TestCheckResourceAttr("ec_protection_protocol.test", "spec.0.mitigation_name", "my-mitigation"),
					resource.TestCheckResourceAttr("ec_protection_protocol.test", "spec.0.protocol", "TCP"),
				),
			},
			{
				ResourceName:      "ec_protection_protocol.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testResourceProtocolsConfigBasic(name string) string {
	return fmt.Sprintf(`resource "ec_protection_protocol" "test" {
  metadata {
    name = "%s"
  }
  spec {
    description = "My Protocol"
    mitigation_name = "my-mitigation"
    protocol = "UDP"
  }
}`, name)
}

func testResourceProtocolsConfigWithTCP(name string) string {
	return fmt.Sprintf(`resource "ec_protection_protocol" "test" {
  metadata {
    name = "%s"
  }
  spec {
    description = "My Protocol"
    mitigation_name = "my-mitigation"
    protocol = "TCP"
  }
}`, name)
}

func testCheckProtocolDestroy(cs clientset.Interface) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "ec_protection_protocol" {
				continue
			}

			name := rs.Primary.ID
			resp, err := cs.ProtectionV1().Protocols().Get(context.Background(), name, metav1.GetOptions{})
			if err == nil {
				if resp.Name == rs.Primary.ID {
					return fmt.Errorf("protocol still exists: %s", rs.Primary.ID)
				}
			}
		}
		return nil
	}
}
