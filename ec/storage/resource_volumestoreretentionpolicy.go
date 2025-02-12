package storage

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec"
	"github.com/nitrado/terraform-provider-ec/pkg/resource"
	"gitlab.com/nitrado/b2b/ec/apicore/api/errors"
	metav1 "gitlab.com/nitrado/b2b/ec/apicore/apis/meta/v1"
	storagev1beta1 "gitlab.com/nitrado/b2b/ec/core/pkg/api/storage/v1beta1"
)

// ResourceVolumeStoreRetentionPolicy returns the resource for a Volume Store Retention Policy.
func ResourceVolumeStoreRetentionPolicy() *schema.Resource {
	return &schema.Resource{
		Description:   "A Volume Store Retention Policy configures how Volume Snapshots are retained.",
		ReadContext:   resourceVolumeStoreRetentionPolicyRead,
		CreateContext: resourceVolumeStoreRetentionPolicyCreate,
		UpdateContext: resourceVolumeStoreRetentionPolicyUpdate,
		DeleteContext: resourceVolumeStoreRetentionPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: volumeStoreRetentionPolicySchema(),
	}
}

func resourceVolumeStoreRetentionPolicyRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	inst, _ := d.Get("instance").(string)
	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	env, name := ec.SplitName(d.Id())

	obj, err := clientSet.StorageV1Beta1().VolumeStoreRetentionPolicies(env).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		switch {
		case errors.IsNotFound(err):
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	data, err := ec.Converter().Flatten(obj, volumeStoreRetentionPolicySchema())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = resource.SetData(d, data); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceVolumeStoreRetentionPolicyCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	inst, _ := d.Get("instance").(string)
	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	obj := &storagev1beta1.VolumeStoreRetentionPolicy{}
	if err = ec.Converter().Expand(d.Get("metadata").([]any), &obj.ObjectMeta); err != nil {
		return diag.FromErr(err)
	}
	if err = ec.Converter().Expand(d.Get("spec").([]any), &obj.Spec); err != nil {
		return diag.FromErr(err)
	}

	out, err := clientSet.StorageV1Beta1().VolumeStoreRetentionPolicies(obj.Environment).Create(ctx, obj, metav1.CreateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ec.ScopedName(out.Environment, out.Name))
	return resourceVolumeStoreRetentionPolicyRead(ctx, d, m)
}

func resourceVolumeStoreRetentionPolicyUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	inst, _ := d.Get("instance").(string)
	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	obj := &storagev1beta1.VolumeStoreRetentionPolicy{}
	if err = ec.Converter().Expand(d.Get("metadata").([]any), &obj.ObjectMeta); err != nil {
		return diag.FromErr(err)
	}
	if err = ec.Converter().Expand(d.Get("spec").([]any), &obj.Spec); err != nil {
		return diag.FromErr(err)
	}

	out, err := clientSet.StorageV1Beta1().VolumeStoreRetentionPolicies(obj.Environment).Update(ctx, obj, metav1.UpdateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ec.ScopedName(out.Environment, out.Name))
	return resourceVolumeStoreRetentionPolicyRead(ctx, d, m)
}

func resourceVolumeStoreRetentionPolicyDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	inst, _ := d.Get("instance").(string)
	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	env, name := ec.SplitName(d.Id())

	if err = clientSet.StorageV1Beta1().VolumeStoreRetentionPolicies(env).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		switch {
		case errors.IsNotFound(err):
			// We will consider this a successful delete.
		default:
			return diag.FromErr(err)
		}
	}

	// Wait for the deletion to complete.
	if err = ec.WaitForDeletion(ctx, clientSet.StorageV1Beta1().VolumeStoreRetentionPolicies(env), name); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
