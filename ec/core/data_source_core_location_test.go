package core_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
	metav1 "gitlab.com/nitrado/b2b/ec/apicore/apis/meta/v1"
	corev1 "gitlab.com/nitrado/b2b/ec/core/pkg/api/core/v1"
)

func TestDataSourceLocations(t *testing.T) {
	loc := &corev1.Location{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-location",
		},
		Spec: corev1.LocationSpec{
			Sites: []string{"site1", "site2"},
		},
	}

	pf, _ := providertest.SetupProviderFactories(t, loc)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLocationConfigRead(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ec_core_location.test", "metadata.0.name", "my-location"),
					resource.TestCheckResourceAttr("data.ec_core_location.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("data.ec_core_location.test", "spec.0.sites.#", "2"),
					resource.TestCheckResourceAttr("data.ec_core_location.test", "spec.0.sites.0", "site1"),
					resource.TestCheckResourceAttr("data.ec_core_location.test", "spec.0.sites.1", "site2"),
				),
			},
		},
	})
}

func testDataSourceLocationConfigRead() string {
	return `data "ec_core_location" "test" {
  metadata {
    name      = "my-location"
  }
}
`
}
