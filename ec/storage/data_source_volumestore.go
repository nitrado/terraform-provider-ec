package storage

import (
	"context"

	apierrors "github.com/gamefabric/gf-apicore/api/errors"
	metav1 "github.com/gamefabric/gf-apicore/apis/meta/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec"
	"github.com/nitrado/terraform-provider-ec/pkg/resource"
)

// DataSourceVolumeStore returns the data source resource for a Volume Store.
func DataSourceVolumeStore() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to access information about an existing Volume Store.",
		ReadContext: dataSourceVolumeStoreRead,
		Schema:      volumeStoreSchema(),
	}
}

func dataSourceVolumeStoreRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	inst, _ := d.Get("instance").(string)
	name := d.Get("metadata.0.name").(string)
	d.SetId(name)

	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	obj, err := clientSet.StorageV1Beta1().VolumeStore().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return diag.Errorf("Volume Store %q not found", name)
		}
		return diag.FromErr(err)
	}

	data, err := ec.Converter().Flatten(obj, volumeStoreSchema())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = resource.SetData(d, data); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
