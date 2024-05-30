package formation

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec"
	"github.com/nitrado/terraform-provider-ec/pkg/resource"
	"gitlab.com/nitrado/b2b/ec/apicore/api/errors"
	metav1 "gitlab.com/nitrado/b2b/ec/apicore/apis/meta/v1"
	formationv1beta1 "gitlab.com/nitrado/b2b/ec/core/pkg/api/formation/v1beta1"
)

// ResourceFormation returns the resource for a Formation.
func ResourceFormation() *schema.Resource {
	return &schema.Resource{
		Description:   "A Formation is a configuration for multiple vessels.",
		ReadContext:   resourceFormationRead,
		CreateContext: resourceFormationCreate,
		UpdateContext: resourceFormationUpdate,
		DeleteContext: resourceFormationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: formationSchema(),
	}
}

func resourceFormationRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	inst, _ := d.Get("instance").(string)
	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	env, name := ec.SplitName(d.Id())

	obj, err := clientSet.FormationV1Beta1().Formations(env).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		switch {
		case errors.IsNotFound(err):
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	data, err := ec.Converter().Flatten(obj, formationSchema())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = resource.SetData(d, data); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceFormationCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	inst, _ := d.Get("instance").(string)
	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	obj := &formationv1beta1.Formation{}
	if err = ec.Converter().Expand(d.Get("metadata").([]any), &obj.ObjectMeta); err != nil {
		return diag.FromErr(err)
	}
	if err = ec.Converter().Expand(d.Get("spec").([]any), &obj.Spec); err != nil {
		return diag.FromErr(err)
	}

	out, err := clientSet.FormationV1Beta1().Formations(obj.Environment).Create(ctx, obj, metav1.CreateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ec.ScopedName(out.Environment, out.Name))
	return resourceFormationRead(ctx, d, m)
}

func resourceFormationUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	inst, _ := d.Get("instance").(string)
	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	obj := &formationv1beta1.Formation{}
	if err = ec.Converter().Expand(d.Get("metadata").([]any), &obj.ObjectMeta); err != nil {
		return diag.FromErr(err)
	}
	if err = ec.Converter().Expand(d.Get("spec").([]any), &obj.Spec); err != nil {
		return diag.FromErr(err)
	}

	out, err := clientSet.FormationV1Beta1().Formations(obj.Environment).Update(ctx, obj, metav1.UpdateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ec.ScopedName(out.Environment, out.Name))
	return resourceFormationRead(ctx, d, m)
}

func resourceFormationDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	inst, _ := d.Get("instance").(string)
	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	env, name := ec.SplitName(d.Id())

	if err = clientSet.FormationV1Beta1().Formations(env).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
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
