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

func TestResourceVolumeStoreRetentionPolicy(t *testing.T) {
	name := "my-volume"
	env := "dflt"
	pf, cs := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		CheckDestroy:      testCheckVolumeStoreRetentionPolicyDestroy(cs),
		Steps: []resource.TestStep{
			{
				Config: testResourceVolumeStoreRetentionPolicyConfigBasic(env, name),
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
				ResourceName:      "ec_storage_volumestoreretentionpolicy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testResourceVolumeStoreRetentionPolicyConfigBasic(env, name string) string {
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
}`, name, env)
}

func testCheckVolumeStoreRetentionPolicyDestroy(cs clientset.Interface) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "ec_storage_volumestoreretentionpolicy" {
				continue
			}

			env, name, _ := strings.Cut(rs.Primary.ID, "/")
			resp, err := cs.StorageV1Beta1().VolumeStoreRetentionPolicies(env).Get(context.Background(), name, metav1.GetOptions{})
			if err == nil {
				if resp.Name == rs.Primary.ID {
					return fmt.Errorf("volume store retention policy still exists: %s", rs.Primary.ID)
				}
			}
		}
		return nil
	}
}
