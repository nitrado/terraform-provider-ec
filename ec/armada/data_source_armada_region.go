package armada

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceArmadaRegion returns the data source resource for a Region.
func DataSourceArmadaRegion() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to access information about an existing Region.",
		ReadContext: dataSourceArmadaRegionRead,
		Schema:      regionSchema(),
	}
}

func dataSourceArmadaRegionRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	name := d.Get("metadata.0.name").(string)
	d.SetId(name)

	return resourceArmadaRegionRead(ctx, d, m)
}
