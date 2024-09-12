package protection

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceGatewayPolicy returns the data source resource for a Gateway Policy.
func DataSourceGatewayPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to access information about an existing Gateway Policy.",
		ReadContext: dataSourceGatewayPolicyRead,
		Schema:      gatewayPolicySchema(),
	}
}

func dataSourceGatewayPolicyRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	name := d.Get("metadata.0.name").(string)
	d.SetId(name)

	return resourceGatewayPolicyRead(ctx, d, m)
}
