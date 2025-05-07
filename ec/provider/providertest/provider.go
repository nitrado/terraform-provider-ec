package providertest

import (
	"context"
	"testing"

	"github.com/gamefabric/gf-apicore/runtime"
	"github.com/gamefabric/gf-core/pkg/apiclient/clientset"
	"github.com/gamefabric/gf-core/pkg/apiclient/fake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec"
	"github.com/nitrado/terraform-provider-ec/ec/provider"
	"github.com/stretchr/testify/require"
)

// SetupProviderFactories returns a configured test terraform provider.
func SetupProviderFactories(t *testing.T, objs ...runtime.Object) (map[string]func() (*schema.Provider, error), clientset.Interface) {
	t.Helper()

	cs, err := fake.New(objs...)
	require.NoError(t, err)

	pf := map[string]func() (*schema.Provider, error){
		//nolint:unparam // Implementing the interface.
		"ec": func() (*schema.Provider, error) {
			p := provider.Provider()
			p.ConfigureContextFunc = func(context.Context, *schema.ResourceData) (any, diag.Diagnostics) {
				return ec.NewProviderContext(cs, map[string]clientset.Interface{
					"test": cs,
				}), nil
			}
			return p, nil
		},
	}

	return pf, cs
}
