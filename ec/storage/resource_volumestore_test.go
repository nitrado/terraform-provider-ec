package storage_test

import (
	"context"
	"fmt"
	"testing"

	metav1 "github.com/gamefabric/gf-apicore/apis/meta/v1"
	"github.com/gamefabric/gf-core/pkg/apiclient/clientset"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
)

func TestResourceVolumeStore(t *testing.T) {
	name := "my-volume-store"
	pf, cs := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		CheckDestroy:      testCheckVolumeStoreDestroy(cs),
		Steps: []resource.TestStep{
			{
				Config: testResourceVolumeStoreConfigBasic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_storage_volumestore.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_storage_volumestore.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_storage_volumestore.test", "spec.0.destination", "gcs:///bucket"),
					resource.TestCheckResourceAttr("ec_storage_volumestore.test", "spec.0.region", "eu"),
					resource.TestCheckResourceAttr("ec_storage_volumestore.test", "spec.0.max_volume_size", "1Gi"),
				),
			},
			{
				ResourceName:      "ec_storage_volumestore.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testResourceVolumeStoreConfigBasic(name string) string {
	return fmt.Sprintf(`resource "ec_storage_volumestore" "test" {
  metadata {
    name = "%s"
  }
  spec {
    destination = "gcs:///bucket"
    region = "eu"
    max_volume_size = "1Gi"
  }
}`, name)
}

func testCheckVolumeStoreDestroy(cs clientset.Interface) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "ec_storage_volumestore" {
				continue
			}

			name := rs.Primary.ID
			resp, err := cs.StorageV1Beta1().VolumeStore().Get(context.Background(), name, metav1.GetOptions{})
			if err == nil {
				if resp.Name == rs.Primary.ID {
					return fmt.Errorf("volume store still exists: %s", rs.Primary.ID)
				}
			}
		}
		return nil
	}
}
