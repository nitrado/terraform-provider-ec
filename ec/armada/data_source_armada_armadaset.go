package armada

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceArmadaArmadaSet returns the data source resource for an ArmadaSet.
func DataSourceArmadaArmadaSet() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to access information about an existing ArmadaSet.",
		ReadContext: dataSourceArmadaArmadaSetRead,
		Schema:      armadaSetSchema(),
	}
}

func dataSourceArmadaArmadaSetRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	name := d.Get("metadata.0.name").(string)
	d.SetId(name)

	return resourceArmadaArmadaSetRead(ctx, d, m)
}
