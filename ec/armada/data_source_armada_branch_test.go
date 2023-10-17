package armada_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
)

func TestArmadaDataSourceBranches(t *testing.T) {
	name := "my-branch"
	pf, _ := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		Steps: []resource.TestStep{
			{
				Config: testArmadasDataSourceBranchesConfigBasic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_armada_branch.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_armada_branch.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_armada_branch.test", "spec.0.description", "My Branch"),
				),
			},
			{
				Config: testArmadasDataSourceBranchesConfigBasic(name) +
					testArmadaDataSourceBranchConfigRead(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ec_armada_branch.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("data.ec_armada_branch.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("data.ec_armada_branch.test", "spec.0.description", "My Branch"),
				),
			},
		},
	})
}

func testArmadasDataSourceBranchesConfigBasic(name string) string {
	return fmt.Sprintf(`resource "ec_armada_branch" "test" {
  metadata {
    name = "%s"
  }
  spec {
    description = "My Branch"
  }
}
`, name)
}

func testArmadaDataSourceBranchConfigRead() string {
	return `data "ec_armada_branch" "test" {
  metadata {
    name      = "${ec_armada_branch.test.metadata.0.name}"
  }
}
`
}
