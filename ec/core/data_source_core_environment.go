package core

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec"
	"github.com/nitrado/terraform-provider-ec/pkg/resource"
	metav1 "gitlab.com/nitrado/b2b/ec/apicore/apis/meta/v1"
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
	inst, _ := d.Get("instance").(string)
	d.SetId(name)

	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	obj, err := clientSet.CoreV1().Environments().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	data, err := ec.Converter().Flatten(obj, environmentSchema())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = resource.SetData(d, data); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
