package armada

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec"
	"github.com/nitrado/terraform-provider-ec/pkg/resource"
	"gitlab.com/nitrado/b2b/ec/apicore/api/errors"
	metav1 "gitlab.com/nitrado/b2b/ec/apicore/apis/meta/v1"
	armadav1 "gitlab.com/nitrado/b2b/ec/armada/pkg/api/armada/v1"
)

// ResourceArmadaRegion returns the resource for a Region.
func ResourceArmadaRegion() *schema.Resource {
	return &schema.Resource{
		Description:   "A Region determines how Armadas are distributed across Sites.",
		ReadContext:   resourceArmadaRegionRead,
		CreateContext: resourceArmadaRegionCreate,
		UpdateContext: resourceArmadaRegionUpdate,
		DeleteContext: resourceArmadaRegionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: regionSchema(),
	}
}

func resourceArmadaRegionRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := ec.ResolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	env, name := ec.SplitName(d.Id())

	obj, err := clientSet.ArmadaV1().Regions(env).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		switch {
		case errors.IsNotFound(err):
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	data, err := ec.Converter().Flatten(obj, regionSchema())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = resource.SetData(d, data); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceArmadaRegionCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := ec.ResolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	obj := &armadav1.Region{
		TypeMeta: metav1.TypeMeta{APIVersion: armadav1.GroupVersion.String(), Kind: "Region"},
	}
	if err = ec.Converter().Expand(d.Get("metadata").([]any), &obj.ObjectMeta); err != nil {
		return diag.FromErr(err)
	}
	if err = ec.Converter().Expand(d.Get("spec").([]any), &obj.Spec); err != nil {
		return diag.FromErr(err)
	}

	out, err := clientSet.ArmadaV1().Regions(obj.Environment).Create(ctx, obj, metav1.CreateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ec.ScopedName(out.Environment, out.Name))
	return resourceArmadaRegionRead(ctx, d, m)
}

func resourceArmadaRegionUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := ec.ResolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	obj := &armadav1.Region{
		TypeMeta: metav1.TypeMeta{APIVersion: armadav1.GroupVersion.String(), Kind: "Region"},
	}
	if err = ec.Converter().Expand(d.Get("metadata").([]any), &obj.ObjectMeta); err != nil {
		return diag.FromErr(err)
	}
	if err = ec.Converter().Expand(d.Get("spec").([]any), &obj.Spec); err != nil {
		return diag.FromErr(err)
	}

	out, err := clientSet.ArmadaV1().Regions(obj.Environment).Update(ctx, obj, metav1.UpdateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ec.ScopedName(out.Environment, out.Name))
	return resourceArmadaRegionRead(ctx, d, m)
}

func resourceArmadaRegionDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := ec.ResolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	env, name := ec.SplitName(d.Id())

	if err = clientSet.ArmadaV1().Regions(env).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
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
