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

// ResourceArmadaSite returns the resource for a Site.
func ResourceArmadaSite() *schema.Resource {
	return &schema.Resource{
		Description:   "A Site provides a connection to deployment capacity.",
		ReadContext:   resourceArmadaSiteRead,
		CreateContext: resourceArmadaSiteCreate,
		UpdateContext: resourceArmadaSiteUpdate,
		DeleteContext: resourceArmadaSiteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: siteSchema(),
	}
}

func resourceArmadaSiteRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := ec.ResolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Id()

	obj, err := clientSet.ArmadaV1().Sites().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		switch {
		case errors.IsNotFound(err):
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	data, err := ec.Converter().Flatten(obj, siteSchema())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = resource.SetData(d, data); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceArmadaSiteCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := ec.ResolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	obj := &armadav1.Site{
		TypeMeta: metav1.TypeMeta{APIVersion: armadav1.GroupVersion.String(), Kind: "Site"},
	}
	if err = ec.Converter().Expand(d.Get("metadata").([]any), &obj.ObjectMeta); err != nil {
		return diag.FromErr(err)
	}
	if err = ec.Converter().Expand(d.Get("spec").([]any), &obj.Spec); err != nil {
		return diag.FromErr(err)
	}

	out, err := clientSet.ArmadaV1().Sites().Create(ctx, obj, metav1.CreateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(out.Name)
	return resourceArmadaSiteRead(ctx, d, m)
}

func resourceArmadaSiteUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := ec.ResolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	obj := &armadav1.Site{
		TypeMeta: metav1.TypeMeta{APIVersion: armadav1.GroupVersion.String(), Kind: "Site"},
	}
	if err = ec.Converter().Expand(d.Get("metadata").([]any), &obj.ObjectMeta); err != nil {
		return diag.FromErr(err)
	}
	if err = ec.Converter().Expand(d.Get("spec").([]any), &obj.Spec); err != nil {
		return diag.FromErr(err)
	}

	out, err := clientSet.ArmadaV1().Sites().Update(ctx, obj, metav1.UpdateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(out.Name)
	return resourceArmadaSiteRead(ctx, d, m)
}

func resourceArmadaSiteDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := ec.ResolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Id()

	if err = clientSet.ArmadaV1().Sites().Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
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
