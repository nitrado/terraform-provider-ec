package protection_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
	metav1 "gitlab.com/nitrado/b2b/ec/apicore/apis/meta/v1"
	protectionv1alpha1 "gitlab.com/nitrado/b2b/ec/core/pkg/api/protection/v1alpha1"
)

func TestDataSourceMigrations(t *testing.T) {
	migration1 := &protectionv1alpha1.Mitigation{
		ObjectMeta: metav1.ObjectMeta{Name: "my-migration-1"},
		Spec: protectionv1alpha1.MitigationSpec{
			DisplayName: "my migration 1",
		},
	}
	migration2 := &protectionv1alpha1.Mitigation{
		ObjectMeta: metav1.ObjectMeta{Name: "my-migration-2"},
		Spec: protectionv1alpha1.MitigationSpec{
			DisplayName: "my migration 2",
		},
	}

	pf, _ := providertest.SetupProviderFactories(t, migration1, migration2)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceImageNameConfigRead(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ec_protection_migration.by_name", "metadata.0.name", "my-migration-1"),
					resource.TestCheckResourceAttr("data.ec_protection_migration.by_name", "spec.#", "1"),
					resource.TestCheckResourceAttr("data.ec_protection_migration.by_name", "spec.0.display_name", "my migration 1"),
				),
			},
		},
	})
}

func testDataSourceImageNameConfigRead() string {
	return `data "ec_protection_migration" "by_name" {
  metadata {
    name   = "my-migration-1"
  }
}
`
}