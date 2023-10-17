package armada

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceArmadaSite returns the data source resource for a Site.
func DataSourceArmadaSite() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceArmadaSiteRead,
		Schema:      siteSchema(),
	}
}

func dataSourceArmadaSiteRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	name := d.Get("metadata.0.name").(string)
	d.SetId(name)

	return resourceArmadaSiteRead(ctx, d, m)
}
