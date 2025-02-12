package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
)

func TestDataSourceVolume(t *testing.T) {
	name := "my-volume"
	env := "dflt"
	pf, _ := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVolumeConfigBasic(name, env),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_storage_volume.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_storage_volume.test", "metadata.0.environment", env),
					resource.TestCheckResourceAttr("ec_storage_volume.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_storage_volume.test", "spec.0.volume_store_name", "my-volume-store"),
					resource.TestCheckResourceAttr("ec_storage_volume.test", "spec.0.capacity", "1Gi"),
				),
			},
			{
				Config: testDataSourceVolumeConfigBasic(name, env) +
					testDataSourceVolumeConfigRead(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ec_storage_volume.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("data.ec_storage_volume.test", "metadata.0.environment", env),
					resource.TestCheckResourceAttr("data.ec_storage_volume.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("data.ec_storage_volume.test", "spec.0.volume_store_name", "my-volume-store"),
					resource.TestCheckResourceAttr("data.ec_storage_volume.test", "spec.0.capacity", "1Gi"),
				),
			},
		},
	})
}

func testDataSourceVolumeConfigBasic(name, env string) string {
	return fmt.Sprintf(`resource "ec_storage_volume" "test" {
  metadata {
    name = "%s"
    environment = "%s"
  }
  spec {
    volume_store_name = "my-volume-store"
    capacity = "1Gi"
  }
}
`, name, env)
}

func testDataSourceVolumeConfigRead() string {
	return `data "ec_storage_volume" "test" {
  metadata {
    name      = "${ec_storage_volume.test.metadata.0.name}"
	environment = "${ec_storage_volume.test.metadata.0.environment}"
  }
}
`
}
