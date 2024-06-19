package armada_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
)

func TestDataSourceArmadas(t *testing.T) {
	name := "my-armada"
	env := "dflt"
	pf, _ := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceArmadasConfigBasic(name, env),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_armada_armada.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_armada_armada.test", "metadata.0.environment", env),
					resource.TestCheckResourceAttr("ec_armada_armada.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_armada_armada.test", "spec.0.description", "My Armada"),
					resource.TestCheckResourceAttr("ec_armada_armada.test", "spec.0.region", "eu"),
					resource.TestCheckResourceAttr("ec_armada_armada.test", "spec.0.distribution.#", "1"),
					resource.TestCheckResourceAttr("ec_armada_armada.test", "spec.0.distribution.0.name", "baremetal"),
					resource.TestCheckResourceAttr("ec_armada_armada.test", "spec.0.distribution.0.min_replicas", "1"),
					resource.TestCheckResourceAttr("ec_armada_armada.test", "spec.0.distribution.0.max_replicas", "2"),
					resource.TestCheckResourceAttr("ec_armada_armada.test", "spec.0.distribution.0.buffer_size", "3"),
					resource.TestCheckResourceAttr("ec_armada_armada.test", "spec.0.template.0.metadata.0.labels.foo", "bar"),
					resource.TestCheckResourceAttr("ec_armada_armada.test", "spec.0.template.0.spec.0.containers.0.name", "my-ctr"),
					resource.TestCheckResourceAttr("ec_armada_armada.test", "spec.0.template.0.spec.0.containers.0.branch", "prod"),
					resource.TestCheckResourceAttr("ec_armada_armada.test", "spec.0.template.0.spec.0.containers.0.image", "test-xyz"),
				),
			},
			{
				Config: testDataSourceArmadasConfigBasic(name, env) +
					testDataSourceArmadaConfigRead(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ec_armada_armada.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("data.ec_armada_armada.test", "metadata.0.environment", env),
					resource.TestCheckResourceAttr("data.ec_armada_armada.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("data.ec_armada_armada.test", "spec.0.description", "My Armada"),
					resource.TestCheckResourceAttr("data.ec_armada_armada.test", "spec.0.region", "eu"),
					resource.TestCheckResourceAttr("data.ec_armada_armada.test", "spec.0.distribution.#", "1"),
					resource.TestCheckResourceAttr("data.ec_armada_armada.test", "spec.0.distribution.0.name", "baremetal"),
					resource.TestCheckResourceAttr("data.ec_armada_armada.test", "spec.0.distribution.0.min_replicas", "1"),
					resource.TestCheckResourceAttr("data.ec_armada_armada.test", "spec.0.distribution.0.max_replicas", "2"),
					resource.TestCheckResourceAttr("data.ec_armada_armada.test", "spec.0.distribution.0.buffer_size", "3"),
					resource.TestCheckResourceAttr("data.ec_armada_armada.test", "spec.0.template.0.metadata.0.labels.foo", "bar"),
					resource.TestCheckResourceAttr("data.ec_armada_armada.test", "spec.0.template.0.spec.0.containers.0.name", "my-ctr"),
					resource.TestCheckResourceAttr("data.ec_armada_armada.test", "spec.0.template.0.spec.0.containers.0.branch", "prod"),
					resource.TestCheckResourceAttr("data.ec_armada_armada.test", "spec.0.template.0.spec.0.containers.0.image", "test-xyz"),
				),
			},
		},
	})
}

func testDataSourceArmadasConfigBasic(name, env string) string {
	return fmt.Sprintf(`resource "ec_armada_armada" "test" {
  metadata {
    name = "%s"
    environment = "%s"
  }
  spec {
    description = "My Armada"
    region = "eu"
    distribution {
      name = "baremetal"
      min_replicas = 1
      max_replicas = 2
      buffer_size = 3
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
`, name, env)
}

func testDataSourceArmadaConfigRead() string {
	return `data "ec_armada_armada" "test" {
  metadata {
    name = "${ec_armada_armada.test.metadata.0.name}"
	environment = "${ec_armada_armada.test.metadata.0.environment}"
  }
}
`
}
