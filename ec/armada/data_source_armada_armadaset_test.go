package armada_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
)

func TestArmadaDataSourceArmadaSets(t *testing.T) {
	name := "my-armadaset"
	pf, _ := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		Steps: []resource.TestStep{
			{
				Config: testArmadasDataSourceArmadaSetsConfigBasic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_armada_armadaset.test", "metadata.0.name", name),
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
				Config: testArmadasDataSourceArmadaSetsConfigBasic(name) +
					testArmadaDataSourceArmadaSetConfigRead(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ec_armada_armadaset.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("data.ec_armada_armadaset.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("data.ec_armada_armadaset.test", "spec.0.description", "My ArmadaSet"),
					resource.TestCheckResourceAttr("data.ec_armada_armadaset.test", "spec.0.armadas.0.region", "eu"),
					resource.TestCheckResourceAttr("data.ec_armada_armadaset.test", "spec.0.armadas.0.distribution.#", "1"),
					resource.TestCheckResourceAttr("data.ec_armada_armadaset.test", "spec.0.armadas.0.distribution.0.name", "baremetal"),
					resource.TestCheckResourceAttr("data.ec_armada_armadaset.test", "spec.0.armadas.0.distribution.0.min_replicas", "1"),
					resource.TestCheckResourceAttr("data.ec_armada_armadaset.test", "spec.0.armadas.0.distribution.0.max_replicas", "2"),
					resource.TestCheckResourceAttr("data.ec_armada_armadaset.test", "spec.0.armadas.0.distribution.0.buffer_size", "3"),
					resource.TestCheckResourceAttr("data.ec_armada_armadaset.test", "spec.0.template.0.metadata.0.labels.foo", "bar"),
					resource.TestCheckResourceAttr("data.ec_armada_armadaset.test", "spec.0.template.0.spec.0.containers.0.name", "my-ctr"),
					resource.TestCheckResourceAttr("data.ec_armada_armadaset.test", "spec.0.template.0.spec.0.containers.0.branch", "prod"),
					resource.TestCheckResourceAttr("data.ec_armada_armadaset.test", "spec.0.template.0.spec.0.containers.0.image", "test-xyz"),
				),
			},
		},
	})
}

func testArmadasDataSourceArmadaSetsConfigBasic(name string) string {
	return fmt.Sprintf(`resource "ec_armada_armadaset" "test" {
  metadata {
    name = "%s"
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
}
`, name)
}

func testArmadaDataSourceArmadaSetConfigRead() string {
	return `data "ec_armada_armadaset" "test" {
  metadata {
    name      = "${ec_armada_armadaset.test.metadata.0.name}"
  }
}
`
}
