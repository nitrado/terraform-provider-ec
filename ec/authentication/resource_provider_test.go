package authentication_test

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

func TestResourceProvider(t *testing.T) {
	t.Parallel()

	resourceName := "ec_authentication_provider.test"
	name := "my-provider"
	pf, cs := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		CheckDestroy:      testCheckProviderDestroy(cs),
		Steps: []resource.TestStep{
			{
				Config: testResourceProviderConfigBasic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "metadata.0.name", name),
					resource.TestCheckResourceAttr(resourceName, "spec.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.display_name", "My Provider"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.oidc.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.oidc.0.issuer", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.oidc.0.client_id", "my-client-id"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.oidc.0.client_secret", "my-client-secret"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.oidc.0.redirect_uri", "https://example.com/callback"),
				),
			},
			{
				Config: testResourceProviderConfigBasicWithScopes(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "metadata.0.name", name),
					resource.TestCheckResourceAttr(resourceName, "spec.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.display_name", "My Provider"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.oidc.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.oidc.0.issuer", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.oidc.0.client_id", "my-client-id"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.oidc.0.client_secret", "my-client-secret"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.oidc.0.redirect_uri", "https://example.com/callback"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.oidc.0.scopes.#", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testResourceProviderConfigBasic(name string) string {
	return fmt.Sprintf(`resource "ec_authentication_provider" "test" {
  metadata {
    name = "%s"
  }
  spec {
    display_name = "My Provider"
    oidc {
      issuer = "https://example.com"
      client_id = "my-client-id"
      client_secret = "my-client-secret"
      redirect_uri = "https://example.com/callback"
    }
  }
}`, name)
}

func testResourceProviderConfigBasicWithScopes(name string) string {
	return fmt.Sprintf(`resource "ec_authentication_provider" "test" {
  metadata {
    name = "%s"
  }
  spec {
    display_name = "My Provider"
    oidc {
      issuer = "https://example.com"
      client_id = "my-client-id"
      client_secret = "my-client-secret"
      redirect_uri = "https://example.com/callback"
      scopes = ["openid", "profile", "email"]
    }
  }
}`, name)
}

func testCheckProviderDestroy(cs clientset.Interface) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "ec_authentication_provider" {
				continue
			}

			name := rs.Primary.ID
			resp, err := cs.AuthenticationV1Beta1().Providers().Get(context.Background(), name, metav1.GetOptions{})
			if err == nil {
				if resp.Name == rs.Primary.ID {
					return fmt.Errorf("provider still exists: %s", rs.Primary.ID)
				}
			}
		}
		return nil
	}
}
