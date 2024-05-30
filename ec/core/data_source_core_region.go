package core

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec"
)

// DataSourceRegion returns the data source resource for a Region.
func DataSourceRegion() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to access information about an existing Region.",
		ReadContext: dataSourceRegionRead,
		Schema:      regionSchema(),
	}
}

func dataSourceRegionRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	name := d.Get("metadata.0.name").(string)
	env := d.Get("metadata.0.environment").(string)
	d.SetId(ec.ScopedName(env, name))

	return resourceRegionRead(ctx, d, m)
}
