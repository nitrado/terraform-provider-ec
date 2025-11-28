package core

import (
	"context"

	apierrors "github.com/gamefabric/gf-apicore/api/errors"
	metav1 "github.com/gamefabric/gf-apicore/apis/meta/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec"
	"github.com/nitrado/terraform-provider-ec/pkg/resource"
)

// DataSourceConfigFile returns the data source resource for a ConfigFile.
func DataSourceConfigFile() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to access information about an existing ConfigFile.",
		ReadContext: dataSourceConfigFileRead,
		Schema:      configFileSchemaDS(),
	}
}

func dataSourceConfigFileRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	inst, _ := d.Get("instance").(string)
	name := d.Get("metadata.0.name").(string)
	env := d.Get("metadata.0.environment").(string)
	d.SetId(ec.ScopedName(env, name))

	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	obj, err := clientSet.CoreV1().ConfigFiles(env).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return diag.Errorf("ConfigFile %q not found in environment %q", name, env)
		}
		return diag.FromErr(err)
	}

	data, err := ec.Converter().Flatten(obj, configFileSchemaDS())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = resource.SetData(d, data); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
