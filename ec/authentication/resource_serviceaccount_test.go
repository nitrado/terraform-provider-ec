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

func TestResourceServiceAccount(t *testing.T) {
	t.Parallel()

	resourceName := "ec_authentication_serviceaccount.test"
	name := "my-service-account"
	pf, cs := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		CheckDestroy:      testCheckServiceAccountDestroy(cs),
		Steps: []resource.TestStep{
			{
				Config: testResourceServiceAccountConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "metadata.0.name", name),
					resource.TestCheckResourceAttr(resourceName, "spec.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.username", "user"),
					resource.TestCheckResourceAttr(resourceName, "spec.0.email", "user@example.com"),
					resource.TestCheckResourceAttrSet(resourceName, "password"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testResourceServiceAccountConfig(name string) string {
	return fmt.Sprintf(`resource "ec_authentication_serviceaccount" "test" {
  metadata {
    name = "%s"
  }
  spec {
    username = "user"
    email = "user@example.com"
  }
}`, name)
}

func testCheckServiceAccountDestroy(cs clientset.Interface) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "ec_authentication_serviceaccount" {
				continue
			}

			name := rs.Primary.ID
			resp, err := cs.AuthenticationV1Beta1().ServiceAccounts().Get(context.Background(), name, metav1.GetOptions{})
			if err == nil {
				if resp.Name == rs.Primary.ID {
					return fmt.Errorf("service account still exists: %s", rs.Primary.ID)
				}
			}
		}
		return nil
	}
}
