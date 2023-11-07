package container_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
	metav1 "gitlab.com/nitrado/b2b/ec/apicore/apis/meta/v1"
	"gitlab.com/nitrado/b2b/ec/armada/pkg/apiclient/clientset"
)

func TestContainerResourceBranch(t *testing.T) {
	name := "my-branch"
	pf, cs := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		CheckDestroy:      testCheckContainerBranchDestroy(cs),
		Steps: []resource.TestStep{
			{
				Config: testContainerResourceBranchesConfigBasic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_container_branch.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_container_branch.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_container_branch.test", "spec.0.description", "My Branch"),
				),
			},
			{
				ResourceName:      "ec_container_branch.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testContainerResourceBranchesConfigBasic(name string) string {
	return fmt.Sprintf(`resource "ec_container_branch" "test" {
  metadata {
    name = "%s"
  }
  spec {
    description = "My Branch"
  }
}`, name)
}

func testCheckContainerBranchDestroy(cs clientset.Interface) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "ec_container_branch" {
				continue
			}

			name := rs.Primary.ID
			resp, err := cs.ContainerV1().Branches().Get(context.Background(), name, metav1.GetOptions{})
			if err == nil {
				if resp.Name == rs.Primary.ID {
					return fmt.Errorf("region still exists: %s", rs.Primary.ID)
				}
			}
		}
		return nil
	}
}
