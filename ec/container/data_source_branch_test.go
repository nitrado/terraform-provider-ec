package container_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
)

func TestDataSourceBranches(t *testing.T) {
	name := "my-branch"
	pf, _ := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceBranchesConfigBasic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_container_branch.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_container_branch.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_container_branch.test", "spec.0.description", "My Branch"),
				),
			},
			{
				Config: testDataSourceBranchesConfigBasic(name) +
					testDataSourceBranchConfigRead(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ec_container_branch.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("data.ec_container_branch.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("data.ec_container_branch.test", "spec.0.description", "My Branch"),
				),
			},
		},
	})
}

func testDataSourceBranchesConfigBasic(name string) string {
	return fmt.Sprintf(`resource "ec_container_branch" "test" {
  metadata {
    name = "%s"
  }
  spec {
    description = "My Branch"
  }
}
`, name)
}

func testDataSourceBranchConfigRead() string {
	return `data "ec_container_branch" "test" {
  metadata {
    name      = "${ec_container_branch.test.metadata.0.name}"
  }
}
`
}
