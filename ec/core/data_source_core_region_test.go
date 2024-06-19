package core_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
)

func TestDataSourceRegions(t *testing.T) {
	name := "my-region"
	env := "dflt"
	pf, _ := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRegionsConfigBasic(name, env),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_core_region.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_core_region.test", "metadata.0.environment", env),
					resource.TestCheckResourceAttr("ec_core_region.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_core_region.test", "spec.0.description", "My Region"),
					resource.TestCheckResourceAttr("ec_core_region.test", "spec.0.types.#", "1"),
					resource.TestCheckResourceAttr("ec_core_region.test", "spec.0.types.0.name", "my-type"),
					resource.TestCheckResourceAttr("ec_core_region.test", "spec.0.types.0.sites.#", "2"),
					resource.TestCheckResourceAttr("ec_core_region.test", "spec.0.types.0.sites.0", "test-site-1"),
					resource.TestCheckResourceAttr("ec_core_region.test", "spec.0.types.0.sites.1", "test-site-2"),
				),
			},
			{
				Config: testDataSourceRegionsConfigBasic(name, env) +
					testDataSourceRegionConfigRead(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ec_core_region.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("data.ec_core_region.test", "metadata.0.environment", env),
					resource.TestCheckResourceAttr("data.ec_core_region.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("data.ec_core_region.test", "spec.0.description", "My Region"),
					resource.TestCheckResourceAttr("data.ec_core_region.test", "spec.0.types.#", "1"),
					resource.TestCheckResourceAttr("data.ec_core_region.test", "spec.0.types.0.name", "my-type"),
					resource.TestCheckResourceAttr("data.ec_core_region.test", "spec.0.types.0.sites.#", "2"),
					resource.TestCheckResourceAttr("data.ec_core_region.test", "spec.0.types.0.sites.0", "test-site-1"),
					resource.TestCheckResourceAttr("data.ec_core_region.test", "spec.0.types.0.sites.1", "test-site-2"),
				),
			},
		},
	})
}

func testDataSourceRegionsConfigBasic(name, env string) string {
	return fmt.Sprintf(`resource "ec_core_region" "test" {
  metadata {
    name = "%s"
    environment = "%s"
  }
  spec {
    description = "My Region"
    types {
      name = "my-type"
      sites = ["test-site-1", "test-site-2"]
    }
  }
}
`, name, env)
}

func testDataSourceRegionConfigRead() string {
	return `data "ec_core_region" "test" {
  metadata {
    name = "${ec_core_region.test.metadata.0.name}"
    environment = "${ec_core_region.test.metadata.0.environment}"
  }
}
`
}
