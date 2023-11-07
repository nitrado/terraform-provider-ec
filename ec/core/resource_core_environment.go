package core

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec"
	"github.com/nitrado/terraform-provider-ec/pkg/resource"
	"gitlab.com/nitrado/b2b/ec/apicore/api/errors"
	metav1 "gitlab.com/nitrado/b2b/ec/apicore/apis/meta/v1"
	corev1 "gitlab.com/nitrado/b2b/ec/armada/pkg/api/core/v1"
)

// ResourceCoreEnvironment returns the resource for an Environment.
func ResourceCoreEnvironment() *schema.Resource {
	return &schema.Resource{
		Description:   "An Environment provides a connection to deployment capacity.",
		ReadContext:   resourceCoreEnvironmentRead,
		CreateContext: resourceCoreEnvironmentCreate,
		UpdateContext: resourceCoreEnvironmentUpdate,
		DeleteContext: resourceCoreEnvironmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: environmentSchema(),
	}
}

func resourceCoreEnvironmentRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := ec.ResolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Id()

	obj, err := clientSet.CoreV1().Environments().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		switch {
		case errors.IsNotFound(err):
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	data, err := ec.Converter().Flatten(obj, environmentSchema())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = resource.SetData(d, data); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceCoreEnvironmentCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := ec.ResolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	obj := &corev1.Environment{}
	if err = ec.Converter().Expand(d.Get("metadata").([]any), &obj.ObjectMeta); err != nil {
		return diag.FromErr(err)
	}
	if err = ec.Converter().Expand(d.Get("spec").([]any), &obj.Spec); err != nil {
		return diag.FromErr(err)
	}

	out, err := clientSet.CoreV1().Environments().Create(ctx, obj, metav1.CreateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(out.Name)
	return resourceCoreEnvironmentRead(ctx, d, m)
}

func resourceCoreEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := ec.ResolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	obj := &corev1.Environment{}
	if err = ec.Converter().Expand(d.Get("metadata").([]any), &obj.ObjectMeta); err != nil {
		return diag.FromErr(err)
	}
	if err = ec.Converter().Expand(d.Get("spec").([]any), &obj.Spec); err != nil {
		return diag.FromErr(err)
	}

	out, err := clientSet.CoreV1().Environments().Update(ctx, obj, metav1.UpdateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(out.Name)
	return resourceCoreEnvironmentRead(ctx, d, m)
}

func resourceCoreEnvironmentDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := ec.ResolveClientSet(m)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Id()

	if err = clientSet.CoreV1().Environments().Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
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
