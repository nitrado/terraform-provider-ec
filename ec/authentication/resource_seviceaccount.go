package authentication

import (
	"context"
	"fmt"

	"github.com/gamefabric/gf-apicore/api/errors"
	metav1 "github.com/gamefabric/gf-apicore/apis/meta/v1"
	authenticationv1beta1 "github.com/gamefabric/gf-core/pkg/api/authentication/v1beta1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec"
	"github.com/nitrado/terraform-provider-ec/pkg/resource"
)

// ResourceServiceAccount returns the resource for a ServiceAccount.
func ResourceServiceAccount() *schema.Resource {
	return &schema.Resource{
		Description:   "A Service Account allows access to the system for machine users.",
		ReadContext:   resourceServiceAccountRead,
		CreateContext: resourceServiceAccountCreate,
		DeleteContext: resourceServiceAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: serviceAccountSchema(),
	}
}

func resourceServiceAccountRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	inst, _ := d.Get("instance").(string)
	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Id()

	obj, err := clientSet.AuthenticationV1Beta1().ServiceAccounts().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		switch {
		case errors.IsNotFound(err):
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	data, err := ec.Converter().Flatten(obj, serviceAccountSchema())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = resource.SetDataKeys(d, data, "metadata", "spec"); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func extractDataKey(data any, key string) any {
	if a, ok := data.([]any); ok && len(a) > 0 {
		data = a[0]
	}
	m, ok := data.(map[string]any)
	if !ok {
		return nil
	}
	v, ok := m[key]
	if !ok {
		return nil
	}
	return v
}

func resourceServiceAccountCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	inst, _ := d.Get("instance").(string)
	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	obj := &authenticationv1beta1.ServiceAccount{}
	if err = ec.Converter().Expand(d.Get("metadata").([]any), &obj.ObjectMeta); err != nil {
		return diag.FromErr(err)
	}
	if err = ec.Converter().Expand(d.Get("spec").([]any), &obj.Spec); err != nil {
		return diag.FromErr(err)
	}

	out, err := clientSet.AuthenticationV1Beta1().ServiceAccounts().Create(ctx, obj, metav1.CreateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(out.Name)
	// Data only available after creation.
	if err = d.Set("password", out.Status.Password); err != nil {
		return diag.FromErr(fmt.Errorf("could not set password: %w", err))
	}
	return resourceServiceAccountRead(ctx, d, m)
}

func resourceServiceAccountDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	inst, _ := d.Get("instance").(string)
	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Id()

	if err = clientSet.AuthenticationV1Beta1().ServiceAccounts().Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		switch {
		case errors.IsNotFound(err):
			// We will consider this a successful delete.
		default:
			return diag.FromErr(err)
		}
	}

	// Wait for the deletion to complete.
	if err = ec.WaitForDeletion(ctx, clientSet.AuthenticationV1Beta1().ServiceAccounts(), name); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
