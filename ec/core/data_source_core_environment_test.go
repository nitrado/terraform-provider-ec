package core_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
)

func TestCoreDataSourceEnvironments(t *testing.T) {
	name := "dflt"
	pf, _ := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		Steps: []resource.TestStep{
			{
				Config: testCoreRDataSourceEnvironmentConfigBasic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_core_environment.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_core_environment.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_core_environment.test", "spec.0.display_name", "My Env"),
					resource.TestCheckResourceAttr("ec_core_environment.test", "spec.0.description", "My Env Description"),
				),
			},
			{
				Config: testCoreRDataSourceEnvironmentConfigBasic(name) +
					testArmadaDataSourceSiteConfigRead(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ec_core_environment.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("data.ec_core_environment.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("data.ec_core_environment.test", "spec.0.display_name", "My Env"),
					resource.TestCheckResourceAttr("data.ec_core_environment.test", "spec.0.description", "My Env Description"),
				),
			},
		},
	})
}

func testCoreRDataSourceEnvironmentConfigBasic(name string) string {
	return fmt.Sprintf(`resource "ec_core_environment" "test" {
  metadata {
    name = "%s"
  }
  spec {
    display_name = "My Env"
    description = "My Env Description"
  }
}
`, name)
}

func testArmadaDataSourceSiteConfigRead() string {
	return `data "ec_core_environment" "test" {
  metadata {
    name      = "${ec_core_environment.test.metadata.0.name}"
  }
}
`
}
