package storage_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	metav1 "github.com/gamefabric/gf-apicore/apis/meta/v1"
	"github.com/gamefabric/gf-core/pkg/apiclient/clientset"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
)

func TestResourceVolume(t *testing.T) {
	name := "my-volume"
	env := "dflt"
	pf, cs := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		CheckDestroy:      testCheckVolumeDestroy(cs),
		Steps: []resource.TestStep{
			{
				Config: testResourceVolumeConfigBasic(env, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_storage_volume.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_storage_volume.test", "metadata.0.environment", env),
					resource.TestCheckResourceAttr("ec_storage_volume.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_storage_volume.test", "spec.0.volume_store_name", "my-volume-store"),
					resource.TestCheckResourceAttr("ec_storage_volume.test", "spec.0.capacity", "1Gi"),
				),
			},
			{
				ResourceName:      "ec_storage_volume.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testResourceVolumeConfigBasic(env, name string) string {
	return fmt.Sprintf(`resource "ec_storage_volume" "test" {
  metadata {
    name = "%s"
    environment = "%s"
  }
  spec {
    volume_store_name = "my-volume-store"
    capacity = "1Gi"
  }
}`, name, env)
}

func testCheckVolumeDestroy(cs clientset.Interface) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "ec_storage_volume" {
				continue
			}

			env, name, _ := strings.Cut(rs.Primary.ID, "/")
			resp, err := cs.StorageV1Beta1().Volumes(env).Get(context.Background(), name, metav1.GetOptions{})
			if err == nil {
				if resp.Name == rs.Primary.ID {
					return fmt.Errorf("volume still exists: %s", rs.Primary.ID)
				}
			}
		}
		return nil
	}
}
