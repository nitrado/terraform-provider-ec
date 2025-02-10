package container

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec"
	"github.com/nitrado/terraform-provider-ec/pkg/resource"
	apierrors "gitlab.com/nitrado/b2b/ec/apicore/api/errors"
	metav1 "gitlab.com/nitrado/b2b/ec/apicore/apis/meta/v1"
)

// DataSourceBranch returns the data source resource for a Branch.
func DataSourceBranch() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to access information about an existing Branch.",
		ReadContext: dataSourceBranchRead,
		Schema:      branchSchema(),
	}
}

func dataSourceBranchRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	inst, _ := d.Get("instance").(string)
	name := d.Get("metadata.0.name").(string)
	d.SetId(name)

	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	obj, err := clientSet.ContainerV1().Branches().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return diag.Errorf("Branch %q not found", name)
		}
		return diag.FromErr(err)
	}

	data, err := ec.Converter().Flatten(obj, branchSchema())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = resource.SetData(d, data); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
