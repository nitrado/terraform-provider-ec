package core

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEnvironment returns the data source resource for an Environment.
func DataSourceEnvironment() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to access information about an existing Environment.",
		ReadContext: dataSourceEnvironmentRead,
		Schema:      environmentSchema(),
	}
}

func dataSourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	name := d.Get("metadata.0.name").(string)
	d.SetId(name)

	return resourceEnvironmentRead(ctx, d, m)
}
