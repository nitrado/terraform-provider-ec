package core_test

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

func TestResourceSecrets(t *testing.T) {
	name := "my-secret"
	env := "dflt"
	pf, cs := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		CheckDestroy:      testCheckSecretDestroy(cs),
		Steps: []resource.TestStep{
			{
				Config: testResourceSecretConfigBasic(env, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_core_secret.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_core_secret.test", "metadata.0.environment", env),
					resource.TestCheckResourceAttr("ec_core_secret.test", "data.%", "1"),
					resource.TestCheckResourceAttr("ec_core_secret.test", "data.key", "value"),
				),
			},
			{
				Config: testResourceSecretsConfigBasicWithDescription(env, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_core_secret.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_core_secret.test", "metadata.0.environment", env),
					resource.TestCheckResourceAttr("ec_core_secret.test", "description", "My Secret"),

					resource.TestCheckResourceAttr("ec_core_secret.test", "data.%", "1"),
					resource.TestCheckResourceAttr("ec_core_secret.test", "data.key", "value"),
				),
			},
		},
	})
}

func testResourceSecretConfigBasic(env, name string) string {
	return fmt.Sprintf(`resource "ec_core_secret" "test" {
  metadata {
    name = "%s"
    environment = "%s"
  }
  data = {
    key = "value"
  }
}`, name, env)
}

func testResourceSecretsConfigBasicWithDescription(env, name string) string {
	return fmt.Sprintf(`resource "ec_core_secret" "test" {
  metadata {
    name = "%s"
    environment = "%s"
  }
  description = "My Secret"
  data = {
    key = "value"
  }
}`, name, env)
}

func testCheckSecretDestroy(cs clientset.Interface) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "ec_core_secret" {
				continue
			}

			env, name, _ := strings.Cut(rs.Primary.ID, "/")
			resp, err := cs.CoreV1().Secrets(env).Get(context.Background(), name, metav1.GetOptions{})
			if err == nil {
				if resp.Name == rs.Primary.ID {
					return fmt.Errorf("secret still exists: %s", rs.Primary.ID)
				}
			}
		}
		return nil
	}
}
