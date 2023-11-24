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

// ResourceArmadaArmadaSet returns the resource for an ArmadaSet.
func ResourceArmadaArmadaSet() *schema.Resource {
	return &schema.Resource{
		Description:   "An ArmadaSet manages Armadas across multiple Regions, while sharing a common specification.",
		ReadContext:   resourceArmadaArmadaSetRead,
		CreateContext: resourceArmadaArmadaSetCreate,
		UpdateContext: resourceArmadaArmadaSetUpdate,
		DeleteContext: resourceArmadaArmadaSetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: armadaSetSchema(),
	}
}

func resourceArmadaArmadaSetRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := ec.ResolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	env, name := ec.SplitName(d.Id())

	obj, err := clientSet.ArmadaV1().ArmadaSets(env).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		switch {
		case errors.IsNotFound(err):
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	data, err := ec.Converter().Flatten(obj, armadaSetSchema())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = resource.SetData(d, data); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceArmadaArmadaSetCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := ec.ResolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	obj := &armadav1.ArmadaSet{
		TypeMeta: metav1.TypeMeta{APIVersion: armadav1.GroupVersion.String(), Kind: "ArmadaSet"},
	}
	if err = ec.Converter().Expand(d.Get("metadata").([]any), &obj.ObjectMeta); err != nil {
		return diag.FromErr(err)
	}
	if err = ec.Converter().Expand(d.Get("spec").([]any), &obj.Spec); err != nil {
		return diag.FromErr(err)
	}

	out, err := clientSet.ArmadaV1().ArmadaSets(obj.Environment).Create(ctx, obj, metav1.CreateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ec.ScopedName(out.Environment, out.Name))
	return resourceArmadaArmadaSetRead(ctx, d, m)
}

func resourceArmadaArmadaSetUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := ec.ResolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	obj := &armadav1.ArmadaSet{
		TypeMeta: metav1.TypeMeta{APIVersion: armadav1.GroupVersion.String(), Kind: "ArmadaSet"},
	}
	if err = ec.Converter().Expand(d.Get("metadata").([]any), &obj.ObjectMeta); err != nil {
		return diag.FromErr(err)
	}
	if err = ec.Converter().Expand(d.Get("spec").([]any), &obj.Spec); err != nil {
		return diag.FromErr(err)
	}

	out, err := clientSet.ArmadaV1().ArmadaSets(obj.Environment).Update(ctx, obj, metav1.UpdateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ec.ScopedName(out.Environment, out.Name))
	return resourceArmadaArmadaSetRead(ctx, d, m)
}

func resourceArmadaArmadaSetDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := ec.ResolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	env, name := ec.SplitName(d.Id())

	if err = clientSet.ArmadaV1().ArmadaSets(env).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
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
