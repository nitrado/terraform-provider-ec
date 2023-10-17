package armada

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceArmadaBranch returns the data source resource for a Branch.
func DataSourceArmadaBranch() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to access information about an existing Branch.",
		ReadContext: dataSourceArmadaBranchRead,
		Schema:      branchSchema(),
	}
}

func dataSourceArmadaBranchRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	name := d.Get("metadata.0.name").(string)
	d.SetId(name)

	return resourceArmadaBranchRead(ctx, d, m)
}
