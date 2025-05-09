package core_test

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

func TestResourceSites(t *testing.T) {
	name := "my-site"
	pf, cs := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		CheckDestroy:      testCheckSiteDestroy(cs),
		Steps: []resource.TestStep{
			{
				Config: testResourceSitesConfigBasic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_core_site.test", "instance", "test"),
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
				),
			},
			{
				Config: testResourceSitesConfigBasic2(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_core_site.test", "instance", "test"),
					resource.TestCheckResourceAttr("ec_core_site.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_core_site.test", "metadata.0.labels.foo", "bar"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.description", "My Other Site"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.credentials.0.endpoint", "endpoint"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.credentials.0.certificate", "cert"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.credentials.0.namespace", "ns"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.credentials.0.token", "tok"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.resources.0.cpu", "250m"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.resources.0.memory", "1Gi"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.resources.0.pods", "100"),
				),
			},
			{
				Config: testResourceSitesConfigBasicWithEnv(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_core_site.test", "instance", "test"),
					resource.TestCheckResourceAttr("ec_core_site.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.description", "My Other Site"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.credentials.0.endpoint", "endpoint"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.credentials.0.certificate", "cert"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.credentials.0.namespace", "ns"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.credentials.0.token", "tok"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.resources.0.cpu", "250m"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.resources.0.memory", "1Gi"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.resources.0.pods", "100"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.template.0.env.#", "2"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.template.0.env.0.name", "foo"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.template.0.env.0.value", "bar"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.template.0.env.1.name", "baz"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.template.0.env.1.value_from.0.config_file_key_ref.0.name", "bat"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.template.0.security_context.0.allow_privilege_escalation.#", "1"),
					resource.TestCheckResourceAttr("ec_core_site.test", "spec.0.template.0.security_context.0.allow_privilege_escalation.0.value", "false"),
				),
			},
			{
				ResourceName:            "ec_core_site.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance"},
			},
		},
	})
}

func testResourceSitesConfigBasic(name string) string {
	return fmt.Sprintf(`resource "ec_core_site" "test" {
  instance = "test"
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
  }
}`, name)
}

func testResourceSitesConfigBasic2(name string) string {
	return fmt.Sprintf(`resource "ec_core_site" "test" {
  instance = "test"
  metadata {
    name = "%s"
    labels = {
      "foo" = "bar"
    }
  }
  spec {
    description = "My Other Site"
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
  }
}`, name)
}

func testResourceSitesConfigBasicWithEnv(name string) string {
	return fmt.Sprintf(`resource "ec_core_site" "test" {
  instance = "test"
  metadata {
    name = "%s"
  }
  spec {
    description = "My Other Site"
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
	  env {
	    name = "baz"
	    value_from {
		  config_file_key_ref {
		    name = "bat"
		  }
	    }
	  }
      security_context {
        allow_privilege_escalation {
		  value = false
        }
      }
    }
  }
}`, name)
}

func testCheckSiteDestroy(cs clientset.Interface) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "ec_core_site" {
				continue
			}

			name := rs.Primary.ID
			resp, err := cs.CoreV1().Sites().Get(context.Background(), name, metav1.GetOptions{})
			if err == nil {
				if resp.Name == rs.Primary.ID {
					return fmt.Errorf("site still exists: %s", rs.Primary.ID)
				}
			}
		}
		return nil
	}
}
