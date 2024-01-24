package core

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceSite returns the data source resource for a Site.
func DataSourceSite() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to access information about an existing Site.",
		ReadContext: dataSourceSiteRead,
		Schema:      siteSchema(),
	}
}

func dataSourceSiteRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	name := d.Get("metadata.0.name").(string)
	d.SetId(name)

	return resourceSiteRead(ctx, d, m)
}
