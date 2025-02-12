package protection

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec"
	"github.com/nitrado/terraform-provider-ec/pkg/resource"
	apierrors "gitlab.com/nitrado/b2b/ec/apicore/api/errors"
	metav1 "gitlab.com/nitrado/b2b/ec/apicore/apis/meta/v1"
)

// DataSourceMigration returns the data source resource for a Migration.
func DataSourceMigration() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to access information about an existing Migration.",
		ReadContext: dataSourceMigrationRead,
		Schema:      migrationSchema(),
	}
}

func dataSourceMigrationRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	inst, _ := d.Get("instance").(string)
	name, hasName := d.GetOk("metadata.0.name")
	if !hasName {
		return diag.FromErr(errors.New("metadata.0.name is required"))
	}
	d.SetId(name.(string))

	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	obj, err := clientSet.ProtectionV1Alpha1().Mitigations().Get(ctx, name.(string), metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return diag.Errorf("Migration %q not found", name)
		}
		return diag.FromErr(err)
	}

	data, err := ec.Converter().Flatten(obj, migrationSchema())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = resource.SetData(d, data); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
