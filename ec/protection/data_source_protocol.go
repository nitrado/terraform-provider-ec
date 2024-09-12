package protection

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceProtocol returns the data source resource for a Protocol.
func DataSourceProtocol() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to access information about an existing Protocol.",
		ReadContext: dataSourceProtocolRead,
		Schema:      protocolSchema(),
	}
}

func dataSourceProtocolRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	name := d.Get("metadata.0.name").(string)
	d.SetId(name)

	return resourceProtocolRead(ctx, d, m)
}
