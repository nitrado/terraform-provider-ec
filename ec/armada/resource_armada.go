package armada

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec"
	"github.com/nitrado/terraform-provider-ec/pkg/resource"
	"gitlab.com/nitrado/b2b/ec/apicore/api/errors"
	metav1 "gitlab.com/nitrado/b2b/ec/apicore/apis/meta/v1"
	armadav1 "gitlab.com/nitrado/b2b/ec/core/pkg/api/armada/v1"
)

// ResourceArmada returns the resource for an Armada.
func ResourceArmada() *schema.Resource {
	return &schema.Resource{
		Description:   "An Armada distributes a specified number of Game Servers across a Region.",
		ReadContext:   resourceArmadaRead,
		CreateContext: resourceArmadaCreate,
		UpdateContext: resourceArmadaUpdate,
		DeleteContext: resourceArmadaDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: armadaSchema(),
	}
}

func resourceArmadaRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := ec.ResolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	env, name := ec.SplitName(d.Id())

	obj, err := clientSet.ArmadaV1().Armadas(env).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		switch {
		case errors.IsNotFound(err):
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	data, err := ec.Converter().Flatten(obj, armadaSchema())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = resource.SetData(d, data); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceArmadaCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := ec.ResolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	obj := &armadav1.Armada{
		TypeMeta: metav1.TypeMeta{APIVersion: armadav1.GroupVersion.String(), Kind: "Armada"},
	}
	if err = ec.Converter().Expand(d.Get("metadata").([]any), &obj.ObjectMeta); err != nil {
		return diag.FromErr(err)
	}
	if err = ec.Converter().Expand(d.Get("spec").([]any), &obj.Spec); err != nil {
		return diag.FromErr(err)
	}

	out, err := clientSet.ArmadaV1().Armadas(obj.Environment).Create(ctx, obj, metav1.CreateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ec.ScopedName(out.Environment, out.Name))
	return resourceArmadaRead(ctx, d, m)
}

func resourceArmadaUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := ec.ResolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	obj := &armadav1.Armada{
		TypeMeta: metav1.TypeMeta{APIVersion: armadav1.GroupVersion.String(), Kind: "Armada"},
	}
	if err = ec.Converter().Expand(d.Get("metadata").([]any), &obj.ObjectMeta); err != nil {
		return diag.FromErr(err)
	}
	if err = ec.Converter().Expand(d.Get("spec").([]any), &obj.Spec); err != nil {
		return diag.FromErr(err)
	}

	out, err := clientSet.ArmadaV1().Armadas(obj.Environment).Update(ctx, obj, metav1.UpdateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ec.ScopedName(out.Environment, out.Name))
	return resourceArmadaRead(ctx, d, m)
}

func resourceArmadaDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := ec.ResolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	env, name := ec.SplitName(d.Id())

	if err = clientSet.ArmadaV1().Armadas(env).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
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
