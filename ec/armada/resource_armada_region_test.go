package armada_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
	metav1 "gitlab.com/nitrado/b2b/ec/armada/pkg/api/apis/meta/v1"
	"gitlab.com/nitrado/b2b/ec/armada/pkg/apiclient/clientset"
)

func TestArmadaResourceRegions(t *testing.T) {
	name := "my-region"
	pf, cs := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		CheckDestroy:      testCheckArmadaRegionDestroy(cs),
		Steps: []resource.TestStep{
			{
				Config: testArmadasResourceRegionsConfigBasic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_armada_region.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_armada_region.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_armada_region.test", "spec.0.description", "My Region"),
					resource.TestCheckResourceAttr("ec_armada_region.test", "spec.0.types.#", "1"),
					resource.TestCheckResourceAttr("ec_armada_region.test", "spec.0.types.0.name", "my-type"),
					resource.TestCheckResourceAttr("ec_armada_region.test", "spec.0.types.0.sites.#", "2"),
					resource.TestCheckResourceAttr("ec_armada_region.test", "spec.0.types.0.sites.0", "test-site-1"),
					resource.TestCheckResourceAttr("ec_armada_region.test", "spec.0.types.0.sites.1", "test-site-2"),
				),
			},
			{
				Config: testArmadasResourceRegionsConfigBasicWithEnv(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_armada_region.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_armada_region.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_armada_region.test", "spec.0.description", "My Region"),
					resource.TestCheckResourceAttr("ec_armada_region.test", "spec.0.types.#", "1"),
					resource.TestCheckResourceAttr("ec_armada_region.test", "spec.0.types.0.name", "my-type"),
					resource.TestCheckResourceAttr("ec_armada_region.test", "spec.0.types.0.sites.#", "2"),
					resource.TestCheckResourceAttr("ec_armada_region.test", "spec.0.types.0.sites.0", "test-site-1"),
					resource.TestCheckResourceAttr("ec_armada_region.test", "spec.0.types.0.sites.1", "test-site-2"),
					resource.TestCheckResourceAttr("ec_armada_region.test", "spec.0.types.0.template.0.env.#", "2"),
					resource.TestCheckResourceAttr("ec_armada_region.test", "spec.0.types.0.template.0.env.0.name", "foo"),
					resource.TestCheckResourceAttr("ec_armada_region.test", "spec.0.types.0.template.0.env.0.value", "bar"),
					resource.TestCheckResourceAttr("ec_armada_region.test", "spec.0.types.0.template.0.env.1.name", "baz"),
					resource.TestCheckResourceAttr("ec_armada_region.test", "spec.0.types.0.template.0.env.1.value_from.0.config_file_key_ref.0.name", "bat"),
				),
			},
			{
				ResourceName:      "ec_armada_region.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testArmadasResourceRegionsConfigBasic(name string) string {
	return fmt.Sprintf(`resource "ec_armada_region" "test" {
  metadata {
    name = "%s"
  }
  spec {
    description = "My Region"
    types {
      name = "my-type"
      sites = ["test-site-1", "test-site-2"]
    }
  }
}`, name)
}

func testArmadasResourceRegionsConfigBasicWithEnv(name string) string {
	return fmt.Sprintf(`resource "ec_armada_region" "test" {
  metadata {
    name = "%s"
  }
  spec {
    description = "My Region"
    types {
      name = "my-type"
      sites = ["test-site-1", "test-site-2"]
      template {
         env {
           name = "foo"
           value = "bar"
          }
		  env {
            name = "baz"
            value_from {
              config_file_key_ref {
                name = "bat"
			  }
			}
		  }
		}
    }
  }
}`, name)
}

func testCheckArmadaRegionDestroy(cs clientset.Interface) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "ec_armada_region" {
				continue
			}

			name := rs.Primary.ID
			resp, err := cs.ArmadaV1().Regions().Get(context.Background(), name, metav1.GetOptions{})
			if err == nil {
				if resp.Name == rs.Primary.ID {
					return fmt.Errorf("region still exists: %s", rs.Primary.ID)
				}
			}
		}
		return nil
	}
}
