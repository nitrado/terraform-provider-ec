package container

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceBranch returns the data source resource for a Branch.
func DataSourceBranch() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to access information about an existing Branch.",
		ReadContext: dataSourceBranchRead,
		Schema:      branchSchema(),
	}
}

func dataSourceBranchRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	name := d.Get("metadata.0.name").(string)
	d.SetId(name)

	return resourceBranchRead(ctx, d, m)
}
