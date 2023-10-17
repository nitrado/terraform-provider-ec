package armada_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
	metav1 "gitlab.com/nitrado/b2b/ec/armada/pkg/api/apis/meta/v1"
	"gitlab.com/nitrado/b2b/ec/armada/pkg/apiclient/clientset"
)

func TestArmadaResourceBranch(t *testing.T) {
	name := "my-branch"
	pf, cs := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		CheckDestroy:      testCheckArmadaBranchDestroy(cs),
		Steps: []resource.TestStep{
			{
				Config: testArmadasResourceBranchesConfigBasic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_armada_branch.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_armada_branch.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_armada_branch.test", "spec.0.description", "My Branch"),
				),
			},
			{
				ResourceName:      "ec_armada_branch.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testArmadasResourceBranchesConfigBasic(name string) string {
	return fmt.Sprintf(`resource "ec_armada_branch" "test" {
  metadata {
    name = "%s"
  }
  spec {
    description = "My Branch"
  }
}`, name)
}

func testCheckArmadaBranchDestroy(cs clientset.Interface) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "ec_armada_branch" {
				continue
			}

			name := rs.Primary.ID
			resp, err := cs.ArmadaV1().Regions().Get(context.Background(), name, metav1.GetOptions{})
			if err == nil {
				if resp.Name == rs.Primary.ID {
					return fmt.Errorf("region still exists: %s", rs.Primary.ID)
				}
			}
		}
		return nil
	}
}
