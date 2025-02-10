package armada

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec"
	"github.com/nitrado/terraform-provider-ec/pkg/resource"
	metav1 "gitlab.com/nitrado/b2b/ec/apicore/apis/meta/v1"
)

// DataSourceArmada returns the data source resource for an Armada.
func DataSourceArmada() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to access information about an existing Armada.",
		ReadContext: dataSourceArmadaRead,
		Schema:      armadaSchema(),
	}
}

func dataSourceArmadaRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	inst, _ := d.Get("instance").(string)
	name := d.Get("metadata.0.name").(string)
	env := d.Get("metadata.0.environment").(string)
	d.SetId(ec.ScopedName(env, name))

	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	obj, err := clientSet.ArmadaV1().Armadas(env).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	data, err := ec.Converter().Flatten(obj, armadaSchema())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = resource.SetData(d, data); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
