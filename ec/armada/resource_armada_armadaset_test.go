package armada_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
	metav1 "gitlab.com/nitrado/b2b/ec/apicore/apis/meta/v1"
	"gitlab.com/nitrado/b2b/ec/armada/pkg/apiclient/clientset"
)

func TestArmadaResourceArmadaSets(t *testing.T) {
	name := "my-armadaset"
	env := "dflt"
	pf, cs := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		CheckDestroy:      testCheckArmadaArmadaSetsDestroy(cs),
		Steps: []resource.TestStep{
			{
				Config: testArmadasResourceArmadaSetsConfigBasic(env, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "metadata.0.environment", env),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.0.description", "My ArmadaSet"),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.0.armadas.0.region", "eu"),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.0.armadas.0.distribution.#", "1"),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.0.armadas.0.distribution.0.name", "baremetal"),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.0.armadas.0.distribution.0.min_replicas", "1"),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.0.armadas.0.distribution.0.max_replicas", "2"),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.0.armadas.0.distribution.0.buffer_size", "3"),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.0.template.0.metadata.0.labels.foo", "bar"),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.0.template.0.spec.0.containers.0.name", "my-ctr"),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.0.template.0.spec.0.containers.0.branch", "prod"),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.0.template.0.spec.0.containers.0.image", "test-xyz"),
				),
			},
			{
				Config: testArmadasResourceArmadaSetsConfigBasicWithEnv(env, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "metadata.0.environment", env),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.0.description", "My ArmadaSet"),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.0.armadas.0.region", "eu"),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.0.armadas.0.distribution.#", "1"),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.0.armadas.0.distribution.0.name", "baremetal"),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.0.armadas.0.distribution.0.min_replicas", "1"),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.0.armadas.0.distribution.0.max_replicas", "2"),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.0.armadas.0.distribution.0.buffer_size", "3"),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.0.template.0.metadata.0.labels.foo", "bar"),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.0.template.0.spec.0.containers.0.name", "my-ctr"),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.0.template.0.spec.0.containers.0.branch", "prod"),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.0.template.0.spec.0.containers.0.image", "test-xyz"),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.0.template.0.spec.0.containers.0.env.#", "2"),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.0.template.0.spec.0.containers.0.env.0.name", "foo"),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.0.template.0.spec.0.containers.0.env.0.value", "bar"),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.0.template.0.spec.0.containers.0.env.1.name", "baz"),
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "spec.0.template.0.spec.0.containers.0.env.1.value_from.0.config_file_key_ref.0.name", "bat"),
				),
			},
			{
				ResourceName:      "ec_armada_armadaset.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testArmadasResourceArmadaSetsConfigBasic(env, name string) string {
	return fmt.Sprintf(`resource "ec_armada_armadaset" "test" {
  metadata {
    name = "%s"
    environment = "%s"
  }
  spec {
    description = "My ArmadaSet"
    armadas {
      name = "eu-armada"
      region = "eu"
      distribution {
        name = "baremetal"
        min_replicas = 1
        max_replicas = 2
        buffer_size = 3
		}
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

func testArmadasResourceArmadaSetsConfigBasicWithEnv(env, name string) string {
	return fmt.Sprintf(`resource "ec_armada_armadaset" "test" {
  metadata {
    name = "%s"
    environment = "%s"
  }
  spec {
    description = "My ArmadaSet"
    armadas {
      name = "eu-armada"
      region = "eu"
      distribution {
        name = "baremetal"
        min_replicas = 1
        max_replicas = 2
        buffer_size = 3
		}
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

func testCheckArmadaArmadaSetsDestroy(cs clientset.Interface) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "ec_armada_armadaset" {
				continue
			}

			env, name, _ := strings.Cut(rs.Primary.ID, "/")
			resp, err := cs.ArmadaV1().ArmadaSets(env).Get(context.Background(), name, metav1.GetOptions{})
			if err == nil {
				if resp.Name == rs.Primary.ID {
					return fmt.Errorf("armada set still exists: %s", rs.Primary.ID)
				}
			}
		}
		return nil
	}
}
