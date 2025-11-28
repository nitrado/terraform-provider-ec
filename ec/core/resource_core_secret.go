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

// ResourceSecret returns the resource for a Secret.
func ResourceSecret() *schema.Resource {
	return &schema.Resource{
		Description:   "Secret data",
		ReadContext:   resourceSecretRead,
		CreateContext: resourceSecretCreate,
		UpdateContext: resourceSecretUpdate,
		DeleteContext: resourceSecretDelete,
		Schema:        secretSchema(),
	}
}

func resourceSecretRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	inst, _ := d.Get("instance").(string)
	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	env, name := ec.SplitName(d.Id())

	obj, err := clientSet.CoreV1().Secrets(env).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		switch {
		case errors.IsNotFound(err):
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	data, err := ec.Converter().Flatten(obj, secretSchema())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = resource.SetDataKeys(d, data, "metadata", "description"); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceSecretCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	inst, _ := d.Get("instance").(string)
	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	obj := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{APIVersion: corev1.GroupVersion.String(), Kind: "Secret"},
	}
	if err = ec.Converter().Expand(d.Get("metadata").([]any), &obj.ObjectMeta); err != nil {
		return diag.FromErr(err)
	}
	if err = ec.Converter().Expand(d.Get("data").(map[string]any), &obj.Data); err != nil {
		return diag.FromErr(err)
	}
	if v, ok := d.GetOk("description"); ok {
		obj.Description = v.(string)
	}

	out, err := clientSet.CoreV1().Secrets(obj.Environment).Create(ctx, obj, metav1.CreateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ec.ScopedName(out.Environment, out.Name))
	return resourceSecretRead(ctx, d, m)
}

func resourceSecretUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	inst, _ := d.Get("instance").(string)
	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	obj := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{APIVersion: corev1.GroupVersion.String(), Kind: "Secret"},
	}
	if err = ec.Converter().Expand(d.Get("metadata").([]any), &obj.ObjectMeta); err != nil {
		return diag.FromErr(err)
	}
	if err = ec.Converter().Expand(d.Get("data").(map[string]any), &obj.Data); err != nil {
		return diag.FromErr(err)
	}
	if v, ok := d.GetOk("description"); ok {
		obj.Description = v.(string)
	}

	out, err := clientSet.CoreV1().Secrets(obj.Environment).Update(ctx, obj, metav1.UpdateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ec.ScopedName(out.Environment, out.Name))
	return resourceSecretRead(ctx, d, m)
}

func resourceSecretDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	inst, _ := d.Get("instance").(string)
	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	env, name := ec.SplitName(d.Id())

	if err = clientSet.CoreV1().Secrets(env).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		switch {
		case errors.IsNotFound(err):
			// We will consider this a successful delete.
		default:
			return diag.FromErr(err)
		}
	}

	// Wait for the deletion to complete.
	if err = ec.WaitForDeletion(ctx, clientSet.CoreV1().Secrets(env), name); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
