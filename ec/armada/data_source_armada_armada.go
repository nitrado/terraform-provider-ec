package armada

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceArmadaArmada returns the data source resource for an Armada.
func DataSourceArmadaArmada() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceArmadaArmadaRead,
		Schema:      armadaSchema(),
	}
}

func dataSourceArmadaArmadaRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	name := d.Get("metadata.0.name").(string)
	d.SetId(name)

	return resourceArmadaArmadaRead(ctx, d, m)
}
