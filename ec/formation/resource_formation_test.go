package formation_test

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

func TestResourceFormations(t *testing.T) {
	name := "my-formation"
	env := "dflt"
	pf, cs := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		CheckDestroy:      testCheckFormationsDestroy(cs),
		Steps: []resource.TestStep{
			{
				Config: testResourceFormationsConfigBasic(env, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_formation_formation.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_formation_formation.test", "metadata.0.environment", env),
					resource.TestCheckResourceAttr("ec_formation_formation.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_formation_formation.test", "spec.0.description", "My Formation"),
					resource.TestCheckResourceAttr("ec_formation_formation.test", "spec.0.vessels.0.region", "eu"),
					resource.TestCheckResourceAttr("ec_formation_formation.test", "spec.0.template.0.metadata.0.labels.foo", "bar"),
					resource.TestCheckResourceAttr("ec_formation_formation.test", "spec.0.template.0.spec.0.containers.0.name", "my-ctr"),
					resource.TestCheckResourceAttr("ec_formation_formation.test", "spec.0.template.0.spec.0.containers.0.branch", "prod"),
					resource.TestCheckResourceAttr("ec_formation_formation.test", "spec.0.template.0.spec.0.containers.0.image", "test-xyz"),
				),
			},
			{
				Config: testResourceFormationsConfigBasicWithEnv(env, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_formation_formation.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_formation_formation.test", "metadata.0.environment", env),
					resource.TestCheckResourceAttr("ec_formation_formation.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_formation_formation.test", "spec.0.description", "My Formation"),
					resource.TestCheckResourceAttr("ec_formation_formation.test", "spec.0.vessels.0.region", "eu"),
					resource.TestCheckResourceAttr("ec_formation_formation.test", "spec.0.template.0.metadata.0.labels.foo", "bar"),
					resource.TestCheckResourceAttr("ec_formation_formation.test", "spec.0.template.0.spec.0.containers.0.name", "my-ctr"),
					resource.TestCheckResourceAttr("ec_formation_formation.test", "spec.0.template.0.spec.0.containers.0.branch", "prod"),
					resource.TestCheckResourceAttr("ec_formation_formation.test", "spec.0.template.0.spec.0.containers.0.image", "test-xyz"),
					resource.TestCheckResourceAttr("ec_formation_formation.test", "spec.0.template.0.spec.0.containers.0.env.#", "2"),
					resource.TestCheckResourceAttr("ec_formation_formation.test", "spec.0.template.0.spec.0.containers.0.env.0.name", "foo"),
					resource.TestCheckResourceAttr("ec_formation_formation.test", "spec.0.template.0.spec.0.containers.0.env.0.value", "bar"),
					resource.TestCheckResourceAttr("ec_formation_formation.test", "spec.0.template.0.spec.0.containers.0.env.1.name", "baz"),
					resource.TestCheckResourceAttr("ec_formation_formation.test", "spec.0.template.0.spec.0.containers.0.env.1.value_from.0.config_file_key_ref.0.name", "bat"),
				),
			},
			{
				ResourceName:      "ec_formation_formation.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testResourceFormationsConfigBasic(env, name string) string {
	return fmt.Sprintf(`resource "ec_formation_formation" "test" {
  metadata {
    name = "%s"
    environment = "%s"
  }
  spec {
    description = "My Formation"
    vessels {
      name = "eu-vessel"
      region = "eu"
    }
    template {
      metadata {
        labels = {
          "foo" = "bar"
        }
      }
      spec {
        containers {
          name = "my-ctr"
          branch = "prod"
          image = "test-xyz"
        }
      }
    }
  }
}`, name, env)
}

func testResourceFormationsConfigBasicWithEnv(env, name string) string {
	return fmt.Sprintf(`resource "ec_formation_formation" "test" {
  metadata {
    name = "%s"
    environment = "%s"
  }
  spec {
    description = "My Formation"
    vessels {
      name = "eu-vessel"
      region = "eu"
    }
    template {
      metadata {
        labels = {
          "foo" = "bar"
        }
      }
      spec {
        containers {
          name = "my-ctr"
          branch = "prod"
          image = "test-xyz"
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
  }
}`, name, env)
}

func testCheckFormationsDestroy(cs clientset.Interface) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "ec_formation_formation" {
				continue
			}

			env, name, _ := strings.Cut(rs.Primary.ID, "/")
			resp, err := cs.FormationV1().Formations(env).Get(context.Background(), name, metav1.GetOptions{})
			if err == nil {
				if resp.Name == rs.Primary.ID {
					return fmt.Errorf("formation still exists: %s", rs.Primary.ID)
				}
			}
		}
		return nil
	}
}
