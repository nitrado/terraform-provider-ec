package core

import (
	"context"

	apierrors "github.com/gamefabric/gf-apicore/api/errors"
	metav1 "github.com/gamefabric/gf-apicore/apis/meta/v1"
	corev1 "github.com/gamefabric/gf-core/pkg/api/core/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec"
	"github.com/nitrado/terraform-provider-ec/pkg/resource"
)

// DataSourceLocation returns the data source resource for a Location.
func DataSourceLocation() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to access information about an existing Location.",
		ReadContext: dataSourceLocationRead,
		Schema:      locationSchema(),
	}
}

func dataSourceLocationRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	inst, _ := d.Get("instance").(string)
	name := d.Get("metadata.0.name").(string)

	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	var obj *corev1.Location
	obj, err = clientSet.CoreV1().Locations().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return diag.Errorf("Location %q not found", name)
		}
		return diag.FromErr(err)
	}

	d.SetId(name)

	data, err := ec.Converter().Flatten(obj, locationSchema())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = resource.SetData(d, data); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
