package armada

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/pkg/resource"
	containerv1 "gitlab.com/nitrado/b2b/ec/armada/pkg/api/apis/container/v1"
	metav1 "gitlab.com/nitrado/b2b/ec/armada/pkg/api/apis/meta/v1"
	"gitlab.com/nitrado/b2b/ec/armada/pkg/api/errors"
)

// ResourceArmadaBranch returns the resource for a Branch.
func ResourceArmadaBranch() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceArmadaBranchRead,
		CreateContext: resourceArmadaBranchCreate,
		UpdateContext: resourceArmadaBranchUpdate,
		DeleteContext: resourceArmadaBranchDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: branchSchema(),
	}
}

func resourceArmadaBranchRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := resolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Id()

	obj, err := clientSet.ContainerV1().Branches().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		switch {
		case errors.IsNotFound(err):
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	data, err := converter().Flatten(obj, branchSchema())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = resource.SetData(d, data); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceArmadaBranchCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := resolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	obj := &containerv1.Branch{
		TypeMeta: metav1.TypeMeta{APIVersion: containerv1.GroupVersion.String(), Kind: "Branch"},
	}
	if err = converter().Expand(d.Get("metadata").([]any), &obj.ObjectMeta); err != nil {
		return diag.FromErr(err)
	}
	if err = converter().Expand(d.Get("spec").([]any), &obj.Spec); err != nil {
		return diag.FromErr(err)
	}

	out, err := clientSet.ContainerV1().Branches().Create(ctx, obj, metav1.CreateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(out.Name)
	return resourceArmadaBranchRead(ctx, d, m)
}

func resourceArmadaBranchUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := resolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	obj := &containerv1.Branch{
		TypeMeta: metav1.TypeMeta{APIVersion: containerv1.GroupVersion.String(), Kind: "Branch"},
	}
	if err = converter().Expand(d.Get("metadata").([]any), &obj.ObjectMeta); err != nil {
		return diag.FromErr(err)
	}
	if err = converter().Expand(d.Get("spec").([]any), &obj.Spec); err != nil {
		return diag.FromErr(err)
	}

	out, err := clientSet.ContainerV1().Branches().Update(ctx, obj, metav1.UpdateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(out.Name)
	return resourceArmadaBranchRead(ctx, d, m)
}

func resourceArmadaBranchDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := resolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Id()

	if err = clientSet.ContainerV1().Branches().Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		switch {
		case errors.IsNotFound(err):
			// We will consider this a successful delete.
		default:
			return diag.FromErr(err)
		}
	}

	d.SetId("")
	return nil
}
