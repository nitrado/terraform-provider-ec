package formation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
)

func TestDataSourceVessels(t *testing.T) {
	name := "my-vessel"
	env := "dflt"
	pf, _ := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVesselsConfigBasic(name, env),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_formation_vessel.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_formation_vessel.test", "metadata.0.environment", env),
					resource.TestCheckResourceAttr("ec_formation_vessel.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_formation_vessel.test", "spec.0.description", "My Vessel"),
					resource.TestCheckResourceAttr("ec_formation_vessel.test", "spec.0.region", "eu"),
					resource.TestCheckResourceAttr("ec_formation_vessel.test", "spec.0.template.0.metadata.0.labels.foo", "bar"),
					resource.TestCheckResourceAttr("ec_formation_vessel.test", "spec.0.template.0.spec.0.containers.0.name", "my-ctr"),
					resource.TestCheckResourceAttr("ec_formation_vessel.test", "spec.0.template.0.spec.0.containers.0.branch", "prod"),
					resource.TestCheckResourceAttr("ec_formation_vessel.test", "spec.0.template.0.spec.0.containers.0.image", "test-xyz"),
				),
			},
			{
				Config: testDataSourceVesselsConfigBasic(name, env) +
					testDataSourceVesselConfigRead(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ec_formation_vessel.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("data.ec_formation_vessel.test", "metadata.0.environment", env),
					resource.TestCheckResourceAttr("data.ec_formation_vessel.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("data.ec_formation_vessel.test", "spec.0.description", "My Vessel"),
					resource.TestCheckResourceAttr("data.ec_formation_vessel.test", "spec.0.region", "eu"),
					resource.TestCheckResourceAttr("data.ec_formation_vessel.test", "spec.0.template.0.metadata.0.labels.foo", "bar"),
					resource.TestCheckResourceAttr("data.ec_formation_vessel.test", "spec.0.template.0.spec.0.containers.0.name", "my-ctr"),
					resource.TestCheckResourceAttr("data.ec_formation_vessel.test", "spec.0.template.0.spec.0.containers.0.branch", "prod"),
					resource.TestCheckResourceAttr("data.ec_formation_vessel.test", "spec.0.template.0.spec.0.containers.0.image", "test-xyz"),
				),
			},
		},
	})
}

func testDataSourceVesselsConfigBasic(name, env string) string {
	return fmt.Sprintf(`resource "ec_formation_vessel" "test" {
  metadata {
    name = "%s"
    environment = "%s"
  }
  spec {
    description = "My Vessel"
    region = "eu"
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

func testDataSourceVesselConfigRead() string {
	return `data "ec_formation_vessel" "test" {
  metadata {
    name = "${ec_formation_vessel.test.metadata.0.name}"
    environment = "${ec_formation_vessel.test.metadata.0.environment}"
  }
}
`
}
