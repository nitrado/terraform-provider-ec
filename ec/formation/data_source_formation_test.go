package formation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
)

func TestDataSourceFormations(t *testing.T) {
	name := "my-armadaset"
	env := "dflt"
	pf, _ := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceFormationConfigBasic(name, env),
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
				Config: testDataSourceFormationConfigBasic(name, env) +
					testDataSourceFormationConfigRead(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ec_formation_formation.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_formation_formation.test", "metadata.0.environment", env),
					resource.TestCheckResourceAttr("data.ec_formation_formation.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("data.ec_formation_formation.test", "spec.0.description", "My Formation"),
					resource.TestCheckResourceAttr("data.ec_formation_formation.test", "spec.0.vessels.0.region", "eu"),
					resource.TestCheckResourceAttr("data.ec_formation_formation.test", "spec.0.template.0.metadata.0.labels.foo", "bar"),
					resource.TestCheckResourceAttr("data.ec_formation_formation.test", "spec.0.template.0.spec.0.containers.0.name", "my-ctr"),
					resource.TestCheckResourceAttr("data.ec_formation_formation.test", "spec.0.template.0.spec.0.containers.0.branch", "prod"),
					resource.TestCheckResourceAttr("data.ec_formation_formation.test", "spec.0.template.0.spec.0.containers.0.image", "test-xyz"),
				),
			},
		},
	})
}

func testDataSourceFormationConfigBasic(name, env string) string {
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
}
`, name, env)
}

func testDataSourceFormationConfigRead() string {
	return `data "ec_formation_formation" "test" {
  metadata {
    name = "${ec_formation_formation.test.metadata.0.name}"
    environment = "${ec_formation_formation.test.metadata.0.environment}"
  }
}
`
}
