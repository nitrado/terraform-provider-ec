package core

import (
	"context"

	"github.com/gamefabric/gf-apicore/api/errors"
	metav1 "github.com/gamefabric/gf-apicore/apis/meta/v1"
	corev1 "github.com/gamefabric/gf-core/pkg/api/core/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec"
	"github.com/nitrado/terraform-provider-ec/pkg/resource"
)

// ResourceRegion returns the resource for a Region.
func ResourceRegion() *schema.Resource {
	return &schema.Resource{
		Description:   "A Region determines how Armadas are distributed across Sites.",
		ReadContext:   resourceRegionRead,
		CreateContext: resourceRegionCreate,
		UpdateContext: resourceRegionUpdate,
		DeleteContext: resourceRegionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: regionSchema(),
	}
}

func resourceRegionRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	inst, _ := d.Get("instance").(string)
	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	env, name := ec.SplitName(d.Id())

	obj, err := clientSet.CoreV1().Regions(env).Get(ctx, name, metav1.GetOptions{})
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

func resourceRegionCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	inst, _ := d.Get("instance").(string)
	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	obj := &corev1.Region{
		TypeMeta: metav1.TypeMeta{APIVersion: corev1.GroupVersion.String(), Kind: "Region"},
	}
	if err = ec.Converter().Expand(d.Get("metadata").([]any), &obj.ObjectMeta); err != nil {
		return diag.FromErr(err)
	}
	if err = ec.Converter().Expand(d.Get("spec").([]any), &obj.Spec); err != nil {
		return diag.FromErr(err)
	}

	out, err := clientSet.CoreV1().Regions(obj.Environment).Create(ctx, obj, metav1.CreateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ec.ScopedName(out.Environment, out.Name))
	return resourceRegionRead(ctx, d, m)
}

func resourceRegionUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	inst, _ := d.Get("instance").(string)
	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	obj := &corev1.Region{
		TypeMeta: metav1.TypeMeta{APIVersion: corev1.GroupVersion.String(), Kind: "Region"},
	}
	if err = ec.Converter().Expand(d.Get("metadata").([]any), &obj.ObjectMeta); err != nil {
		return diag.FromErr(err)
	}
	if err = ec.Converter().Expand(d.Get("spec").([]any), &obj.Spec); err != nil {
		return diag.FromErr(err)
	}

	out, err := clientSet.CoreV1().Regions(obj.Environment).Update(ctx, obj, metav1.UpdateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ec.ScopedName(out.Environment, out.Name))
	return resourceRegionRead(ctx, d, m)
}

func resourceRegionDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	inst, _ := d.Get("instance").(string)
	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	env, name := ec.SplitName(d.Id())

	if err = clientSet.CoreV1().Regions(env).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		switch {
		case errors.IsNotFound(err):
			// We will consider this a successful delete.
		default:
			return diag.FromErr(err)
		}
	}

	// Wait for the deletion to complete.
	if err = ec.WaitForDeletion(ctx, clientSet.CoreV1().Regions(env), name); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
