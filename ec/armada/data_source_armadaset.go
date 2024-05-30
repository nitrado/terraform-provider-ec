package armada

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec"
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
	name := d.Get("metadata.0.name").(string)
	env := d.Get("metadata.0.environment").(string)
	d.SetId(ec.ScopedName(env, name))

	return resourceArmadaSetRead(ctx, d, m)
}
