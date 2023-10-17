package armada

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/pkg/resource"
	armadav1 "gitlab.com/nitrado/b2b/ec/armada/pkg/api/apis/armada/v1"
	metav1 "gitlab.com/nitrado/b2b/ec/armada/pkg/api/apis/meta/v1"
	"gitlab.com/nitrado/b2b/ec/armada/pkg/api/errors"
)

// ResourceArmadaArmadaSet returns the resource for an ArmadaSet.
func ResourceArmadaArmadaSet() *schema.Resource {
	return &schema.Resource{
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
	clientSet, err := resolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Id()

	obj, err := clientSet.ArmadaV1().ArmadaSets().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		switch {
		case errors.IsNotFound(err):
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	data, err := converter().Flatten(obj, armadaSetSchema())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = resource.SetData(d, data); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceArmadaArmadaSetCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := resolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	obj := &armadav1.ArmadaSet{
		TypeMeta: metav1.TypeMeta{APIVersion: armadav1.GroupVersion.String(), Kind: "ArmadaSet"},
	}
	if err = converter().Expand(d.Get("metadata").([]any), &obj.ObjectMeta); err != nil {
		return diag.FromErr(err)
	}
	if err = converter().Expand(d.Get("spec").([]any), &obj.Spec); err != nil {
		return diag.FromErr(err)
	}

	out, err := clientSet.ArmadaV1().ArmadaSets().Create(ctx, obj, metav1.CreateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(out.Name)
	return resourceArmadaArmadaSetRead(ctx, d, m)
}

func resourceArmadaArmadaSetUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := resolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	obj := &armadav1.ArmadaSet{
		TypeMeta: metav1.TypeMeta{APIVersion: armadav1.GroupVersion.String(), Kind: "ArmadaSet"},
	}
	if err = converter().Expand(d.Get("metadata").([]any), &obj.ObjectMeta); err != nil {
		return diag.FromErr(err)
	}
	if err = converter().Expand(d.Get("spec").([]any), &obj.Spec); err != nil {
		return diag.FromErr(err)
	}

	out, err := clientSet.ArmadaV1().ArmadaSets().Update(ctx, obj, metav1.UpdateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(out.Name)
	return resourceArmadaArmadaSetRead(ctx, d, m)
}

func resourceArmadaArmadaSetDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := resolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Id()

	if err = clientSet.ArmadaV1().ArmadaSets().Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
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
