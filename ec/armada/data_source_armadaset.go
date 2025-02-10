package armada

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec"
	"github.com/nitrado/terraform-provider-ec/pkg/resource"
	apierrors "gitlab.com/nitrado/b2b/ec/apicore/api/errors"
	metav1 "gitlab.com/nitrado/b2b/ec/apicore/apis/meta/v1"
)

// DataSourceArmadaSet returns the data source resource for an ArmadaSet.
func DataSourceArmadaSet() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to access information about an existing ArmadaSet.",
		ReadContext: dataSourceArmadaSetRead,
		Schema:      armadaSetSchema(),
	}
}

func dataSourceArmadaSetRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	inst, _ := d.Get("instance").(string)
	name := d.Get("metadata.0.name").(string)
	env := d.Get("metadata.0.environment").(string)
	d.SetId(ec.ScopedName(env, name))

	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	obj, err := clientSet.ArmadaV1().ArmadaSets(env).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return diag.Errorf("ArmadaSet %q not found in environment %q", name, env)
		}
		return diag.FromErr(err)
	}

	data, err := ec.Converter().Flatten(obj, armadaSetSchema())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = resource.SetData(d, data); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
