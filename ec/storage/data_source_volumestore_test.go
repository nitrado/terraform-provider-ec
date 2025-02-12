package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
)

func TestDataSourceVolumeStore(t *testing.T) {
	name := "my-volume-store"
	pf, _ := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVolumeStoreConfigBasic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_storage_volumestore.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_storage_volumestore.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_storage_volumestore.test", "spec.0.destination", "gcs:///bucket"),
					resource.TestCheckResourceAttr("ec_storage_volumestore.test", "spec.0.region", "eu"),
					resource.TestCheckResourceAttr("ec_storage_volumestore.test", "spec.0.max_volume_size", "1Gi"),
				),
			},
			{
				Config: testDataSourceVolumeStoreConfigBasic(name) +
					testDataSourceVolumeStoreConfigRead(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ec_storage_volumestore.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("data.ec_storage_volumestore.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("data.ec_storage_volumestore.test", "spec.0.destination", "gcs:///bucket"),
					resource.TestCheckResourceAttr("data.ec_storage_volumestore.test", "spec.0.region", "eu"),
					resource.TestCheckResourceAttr("data.ec_storage_volumestore.test", "spec.0.max_volume_size", "1Gi"),
				),
			},
		},
	})
}

func testDataSourceVolumeStoreConfigBasic(name string) string {
	return fmt.Sprintf(`resource "ec_storage_volumestore" "test" {
  metadata {
    name = "%s"
  }
  spec {
    destination = "gcs:///bucket"
    region = "eu"
    max_volume_size = "1Gi"
  }
}
`, name)
}

func testDataSourceVolumeStoreConfigRead() string {
	return `data "ec_storage_volumestore" "test" {
  metadata {
    name      = "${ec_storage_volumestore.test.metadata.0.name}"
  }
}
`
}
