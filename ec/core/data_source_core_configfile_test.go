package core_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
)

func TestDataSourceConfigFiles(t *testing.T) {
	name := "my-region"
	env := "dflt"
	pf, _ := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceConfigFilesConfigBasic(name, env),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_core_configfile.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_core_configfile.test", "metadata.0.environment", env),
					resource.TestCheckResourceAttr("ec_core_configfile.test", "description", "My Config File"),
					resource.TestCheckResourceAttr("ec_core_configfile.test", "data", "some data"),
				),
			},
			{
				Config: testDataSourceConfigFilesConfigBasic(name, env) +
					testDataSourceConfigFileConfigRead(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ec_core_configfile.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("data.ec_core_configfile.test", "metadata.0.environment", env),
					resource.TestCheckResourceAttr("data.ec_core_configfile.test", "description", "My Config File"),
					resource.TestCheckResourceAttr("data.ec_core_configfile.test", "data", "some data"),
				),
			},
		},
	})
}

func testDataSourceConfigFilesConfigBasic(name, env string) string {
	return fmt.Sprintf(`resource "ec_core_configfile" "test" {
  metadata {
    name = "%s"
    environment = "%s"
  }
  description = "My Config File"
  data = "some data"
}
`, name, env)
}

func testDataSourceConfigFileConfigRead() string {
	return `data "ec_core_configfile" "test" {
  metadata {
    name = "${ec_core_configfile.test.metadata.0.name}"
    environment = "${ec_core_configfile.test.metadata.0.environment}"
  }
}
`
}
