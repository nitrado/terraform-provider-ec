package armada_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
)

func TestArmadaDataSourceRegions(t *testing.T) {
	name := "my-region"
	pf, _ := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		Steps: []resource.TestStep{
			{
				Config: testArmadasDataSourceRegionsConfigBasic(name),
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
				Config: testArmadasDataSourceRegionsConfigBasic(name) +
					testArmadaDataSourceRegionConfigRead(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ec_armada_region.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("data.ec_armada_region.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("data.ec_armada_region.test", "spec.0.description", "My Region"),
					resource.TestCheckResourceAttr("data.ec_armada_region.test", "spec.0.types.#", "1"),
					resource.TestCheckResourceAttr("data.ec_armada_region.test", "spec.0.types.0.name", "my-type"),
					resource.TestCheckResourceAttr("data.ec_armada_region.test", "spec.0.types.0.sites.#", "2"),
					resource.TestCheckResourceAttr("data.ec_armada_region.test", "spec.0.types.0.sites.0", "test-site-1"),
					resource.TestCheckResourceAttr("data.ec_armada_region.test", "spec.0.types.0.sites.1", "test-site-2"),
				),
			},
		},
	})
}

func testArmadasDataSourceRegionsConfigBasic(name string) string {
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
}
`, name)
}

func testArmadaDataSourceRegionConfigRead() string {
	return `data "ec_armada_region" "test" {
  metadata {
    name      = "${ec_armada_region.test.metadata.0.name}"
  }
}
`
}
