package formation

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec"
)

// DataSourceVessel returns the data source resource for a Vessel.
func DataSourceVessel() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to access information about an existing Vessel.",
		ReadContext: dataSourceVesselRead,
		Schema:      vesselSchema(),
	}
}

func dataSourceVesselRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	name := d.Get("metadata.0.name").(string)
	env := d.Get("metadata.0.environment").(string)
	d.SetId(ec.ScopedName(env, name))

	return resourceVesselRead(ctx, d, m)
}
