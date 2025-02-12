package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
)

func TestDataSourceVolumeStoreRetentionPolicy(t *testing.T) {
	name := "my-volume-store-retention-policy"
	env := "dflt"
	pf, _ := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVolumeStoreRetentionPolicyConfigBasic(name, env),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_storage_volumestoreretentionpolicy.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_storage_volumestoreretentionpolicy.test", "metadata.0.environment", env),
					resource.TestCheckResourceAttr("ec_storage_volumestoreretentionpolicy.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_storage_volumestoreretentionpolicy.test", "spec.0.volume_store_name", "my-volume-store"),
					resource.TestCheckResourceAttr("ec_storage_volumestoreretentionpolicy.test", "spec.0.snapshots.#", "1"),
					resource.TestCheckResourceAttr("ec_storage_volumestoreretentionpolicy.test", "spec.0.snapshots.0.offline_count", "2"),
				),
			},
			{
				Config: testDataSourceVolumeStoreRetentionPolicyConfigBasic(name, env) +
					testDataSourceVolumeStoreRetentionPolicyConfigRead(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ec_storage_volumestoreretentionpolicy.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("data.ec_storage_volumestoreretentionpolicy.test", "metadata.0.environment", env),
					resource.TestCheckResourceAttr("data.ec_storage_volumestoreretentionpolicy.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("data.ec_storage_volumestoreretentionpolicy.test", "spec.0.volume_store_name", "my-volume-store"),
					resource.TestCheckResourceAttr("data.ec_storage_volumestoreretentionpolicy.test", "spec.0.snapshots.#", "1"),
					resource.TestCheckResourceAttr("data.ec_storage_volumestoreretentionpolicy.test", "spec.0.snapshots.0.offline_count", "2"),
				),
			},
		},
	})
}

func testDataSourceVolumeStoreRetentionPolicyConfigBasic(name, env string) string {
	return fmt.Sprintf(`resource "ec_storage_volumestoreretentionpolicy" "test" {
  metadata {
    name = "%s"
    environment = "%s"
  }
  spec {
    volume_store_name = "my-volume-store"
    snapshots {
      offline_count = 2
    }
  }
}
`, name, env)
}

func testDataSourceVolumeStoreRetentionPolicyConfigRead() string {
	return `data "ec_storage_volumestoreretentionpolicy" "test" {
  metadata {
    name      = "${ec_storage_volumestoreretentionpolicy.test.metadata.0.name}"
	environment = "${ec_storage_volumestoreretentionpolicy.test.metadata.0.environment}"
  }
}
`
}
