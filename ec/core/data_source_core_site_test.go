package core_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/nitrado/terraform-provider-ec/ec/provider/providertest"
)

func TestDataSourceSites(t *testing.T) {
	name := "my-site"
	pf, _ := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSitesConfigBasic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_core_site.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.description", "My Site"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.credentials.0.endpoint", "endpoint"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.credentials.0.certificate", "cert"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.credentials.0.namespace", "ns"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.credentials.0.token", "tok"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.resources.0.cpu", "250m"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.resources.0.memory", "1Gi"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.resources.0.pods", "100"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.template.0.image_pull_secrets.0", "test-secret"),
				),
			},
			{
				Config: testDataSourceSitesConfigBasic(name) +
					testDataSourceSiteConfigRead(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ec_core_site.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("data.ec_core_site.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("data.ec_core_site.test", "spec.0.description", "My Site"),
					resource.TestCheckResourceAttr("data.ec_core_site.test", "spec.0.credentials.0.endpoint", "endpoint"),
					resource.TestCheckResourceAttr("data.ec_core_site.test", "spec.0.credentials.0.certificate", "cert"),
					resource.TestCheckResourceAttr("data.ec_core_site.test", "spec.0.credentials.0.namespace", "ns"),
					resource.TestCheckResourceAttr("data.ec_core_site.test", "spec.0.credentials.0.token", "tok"),
					resource.TestCheckResourceAttr("data.ec_core_site.test", "spec.0.resources.0.cpu", "250m"),
					resource.TestCheckResourceAttr("data.ec_core_site.test", "spec.0.resources.0.memory", "1Gi"),
					resource.TestCheckResourceAttr("data.ec_core_site.test", "spec.0.resources.0.pods", "100"),
					resource.TestCheckResourceAttr("data.ec_core_site.test", "spec.0.template.0.image_pull_secrets.0", "test-secret"),
					resource.TestCheckResourceAttr("data.ec_core_site.test", "spec.0.template.0.env.0.name", "foo"),
					resource.TestCheckResourceAttr("data.ec_core_site.test", "spec.0.template.0.env.0.value", "bar"),
				),
			},
		},
	})
}

func testDataSourceSitesConfigBasic(name string) string {
	return fmt.Sprintf(`resource "ec_core_site" "test" {
  metadata {
    name = "%s"
  }
  spec {
    description = "My Site"
    credentials {
      endpoint    = "endpoint"
      certificate = "cert"
      namespace   = "ns"
      token       = "tok"
    }
    resources {
      cpu = "250m"
      memory = "1Gi"
      pods = 100
	}
	template {
      env {
        name = "foo"
        value = "bar"
      }
	  image_pull_secrets = [ "test-secret" ]
    }
  }
}
`, name)
}

func testDataSourceSiteConfigRead() string {
	return `data "ec_core_site" "test" {
  metadata {
    name      = "${ec_core_site.test.metadata.0.name}"
  }
}
`
}
