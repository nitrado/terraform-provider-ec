package core

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceCoreEnvironment returns the data source resource for an Environment.
func DataSourceCoreEnvironment() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to access information about an existing Environment.",
		ReadContext: dataSourceCoreEnvironmentRead,
		Schema:      environmentSchema(),
	}
}

func dataSourceCoreEnvironmentRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	name := d.Get("metadata.0.name").(string)
	d.SetId(name)

	return resourceCoreEnvironmentRead(ctx, d, m)
}
