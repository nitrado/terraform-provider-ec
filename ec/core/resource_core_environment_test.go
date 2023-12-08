package core_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
	metav1 "gitlab.com/nitrado/b2b/ec/apicore/apis/meta/v1"
	"gitlab.com/nitrado/b2b/ec/armada/pkg/apiclient/clientset"
)

func TestResourceEnvironments(t *testing.T) {
	name := "dflt"
	pf, cs := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		CheckDestroy:      testCheckEnvironmentDestroy(cs),
		Steps: []resource.TestStep{
			{
				Config: testResourceEnvironmentConfigBasic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_core_environment.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_core_environment.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_core_environment.test", "spec.0.display_name", "My Env"),
					resource.TestCheckResourceAttr("ec_core_environment.test", "spec.0.description", "My Env Description"),
				),
			},
			{
				ResourceName:      "ec_core_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testResourceEnvironmentConfigBasic(name string) string {
	return fmt.Sprintf(`resource "ec_core_environment" "test" {
  metadata {
    name = "%s"
  }
  spec {
    display_name = "My Env"
    description = "My Env Description"
  }
}`, name)
}

func testCheckEnvironmentDestroy(cs clientset.Interface) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "ec_core_environment" {
				continue
			}

			name := rs.Primary.ID
			resp, err := cs.CoreV1().Environments().Get(context.Background(), name, metav1.GetOptions{})
			if err == nil {
				if resp.Name == rs.Primary.ID {
					return fmt.Errorf("Environment still exists: %s", rs.Primary.ID)
				}
			}
		}
		return nil
	}
}
