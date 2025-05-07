package core_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	metav1 "github.com/gamefabric/gf-apicore/apis/meta/v1"
	"github.com/gamefabric/gf-core/pkg/apiclient/clientset"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
)

func TestResourceRegions(t *testing.T) {
	name := "my-region"
	env := "dflt"
	pf, cs := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		CheckDestroy:      testCheckArmadaRegionDestroy(cs),
		Steps: []resource.TestStep{
			{
				Config: testResourceRegionsConfigBasic(env, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_core_region.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_core_region.test", "metadata.0.environment", env),
					resource.TestCheckResourceAttr("ec_core_region.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_core_region.test", "spec.0.description", "My Region"),
					resource.TestCheckResourceAttr("ec_core_region.test", "spec.0.types.#", "1"),
					resource.TestCheckResourceAttr("ec_core_region.test", "spec.0.types.0.name", "my-type"),
					resource.TestCheckResourceAttr("ec_core_region.test", "spec.0.types.0.locations.#", "2"),
					resource.TestCheckResourceAttr("ec_core_region.test", "spec.0.types.0.locations.0", "test-loc-1"),
					resource.TestCheckResourceAttr("ec_core_region.test", "spec.0.types.0.locations.1", "test-loc-2"),
				),
			},
			{
				Config: testResourceRegionsConfigBasicWithEnv(env, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_core_region.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_core_region.test", "metadata.0.environment", env),
					resource.TestCheckResourceAttr("ec_core_region.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_core_region.test", "spec.0.description", "My Region"),
					resource.TestCheckResourceAttr("ec_core_region.test", "spec.0.types.#", "1"),
					resource.TestCheckResourceAttr("ec_core_region.test", "spec.0.types.0.name", "my-type"),
					resource.TestCheckResourceAttr("ec_core_region.test", "spec.0.types.0.locations.#", "2"),
					resource.TestCheckResourceAttr("ec_core_region.test", "spec.0.types.0.locations.0", "test-loc-1"),
					resource.TestCheckResourceAttr("ec_core_region.test", "spec.0.types.0.locations.1", "test-loc-2"),
					resource.TestCheckResourceAttr("ec_core_region.test", "spec.0.types.0.template.0.env.#", "2"),
					resource.TestCheckResourceAttr("ec_core_region.test", "spec.0.types.0.template.0.env.0.name", "foo"),
					resource.TestCheckResourceAttr("ec_core_region.test", "spec.0.types.0.template.0.env.0.value", "bar"),
					resource.TestCheckResourceAttr("ec_core_region.test", "spec.0.types.0.template.0.env.1.name", "baz"),
					resource.TestCheckResourceAttr("ec_core_region.test", "spec.0.types.0.template.0.env.1.value_from.0.config_file_key_ref.0.name", "bat"),
				),
			},
			{
				ResourceName:      "ec_core_region.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testResourceRegionsConfigBasic(env, name string) string {
	return fmt.Sprintf(`resource "ec_core_region" "test" {
  metadata {
    name = "%s"
    environment = "%s"
  }
  spec {
    description = "My Region"
    types {
      name = "my-type"
      locations = ["test-loc-1", "test-loc-2"]
    }
  }
}`, name, env)
}

func testResourceRegionsConfigBasicWithEnv(env, name string) string {
	return fmt.Sprintf(`resource "ec_core_region" "test" {
  metadata {
    name = "%s"
    environment = "%s"
  }
  spec {
    description = "My Region"
    types {
      name = "my-type"
      locations = ["test-loc-1", "test-loc-2"]
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
}`, name, env)
}

func testCheckArmadaRegionDestroy(cs clientset.Interface) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "ec_core_region" {
				continue
			}

			env, name, _ := strings.Cut(rs.Primary.ID, "/")
			resp, err := cs.CoreV1().Regions(env).Get(context.Background(), name, metav1.GetOptions{})
			if err == nil {
				if resp.Name == rs.Primary.ID {
					return fmt.Errorf("region still exists: %s", rs.Primary.ID)
				}
			}
		}
		return nil
	}
}
