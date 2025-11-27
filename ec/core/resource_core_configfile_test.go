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

func TestResourceConfigFiles(t *testing.T) {
	name := "my-config-file"
	env := "dflt"
	pf, cs := providertest.SetupProviderFactories(t)

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: pf,
		CheckDestroy:      testCheckConfigFileDestroy(cs),
		Steps: []resource.TestStep{
			{
				Config: testResourceConfigFileConfigBasic(env, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_core_configfile.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_core_configfile.test", "metadata.0.environment", env),
					resource.TestCheckResourceAttr("ec_core_configfile.test", "data", "some data"),
				),
			},
			{
				Config: testResourceConfigFilesConfigBasicWithDescription(env, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ec_core_configfile.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("ec_core_configfile.test", "metadata.0.environment", env),
					resource.TestCheckResourceAttr("ec_core_configfile.test", "description", "My Config File"),
					resource.TestCheckResourceAttr("ec_core_configfile.test", "data", "some data"),
				),
			},
			{
				ResourceName:      "ec_core_configfile.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testResourceConfigFileConfigBasic(env, name string) string {
	return fmt.Sprintf(`resource "ec_core_configfile" "test" {
  metadata {
    name = "%s"
    environment = "%s"
  }
  data = "some data"
}`, name, env)
}

func testResourceConfigFilesConfigBasicWithDescription(env, name string) string {
	return fmt.Sprintf(`resource "ec_core_configfile" "test" {
  metadata {
    name = "%s"
    environment = "%s"
  }
  description = "My Config File"
  data = "some data"
}`, name, env)
}

func testCheckConfigFileDestroy(cs clientset.Interface) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "ec_core_configfile" {
				continue
			}

			env, name, _ := strings.Cut(rs.Primary.ID, "/")
			resp, err := cs.CoreV1().ConfigFiles(env).Get(context.Background(), name, metav1.GetOptions{})
			if err == nil {
				if resp.Name == rs.Primary.ID {
					return fmt.Errorf("config file still exists: %s", rs.Primary.ID)
				}
			}
		}
		return nil
	}
}
