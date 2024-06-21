package container

import (
	"context"
	"errors"
	"slices"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec"
	"github.com/nitrado/terraform-provider-ec/pkg/resource"
	apierrors "gitlab.com/nitrado/b2b/ec/apicore/api/errors"
	metav1 "gitlab.com/nitrado/b2b/ec/apicore/apis/meta/v1"
	containerv1 "gitlab.com/nitrado/b2b/ec/core/pkg/api/container/v1"
)

// DataSourceImage returns the data source resource for an image.
func DataSourceImage() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to access information about an existing Image.",
		ReadContext: dataSourceImageRead,
		Schema:      imageSchema(),
	}
}

func dataSourceImageRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	inst, _ := d.Get("instance").(string)
	clientSet, err := ec.ResolveClientSet(m, inst)
	if err != nil {
		return diag.FromErr(err)
	}

	branch := d.Get("metadata.0.branch").(string)
	name, hasName := d.GetOk("metadata.0.object_meta.0.name")
	image, hasImage := d.GetOk("spec.0.image")
	tag, hasTag := d.GetOk("spec.0.tag")

	var obj *containerv1.Image
	switch {
	case hasName:
		obj, err = clientSet.ContainerV1().Images(branch).Get(ctx, name.(string), metav1.GetOptions{})
		if err != nil {
			switch {
			case apierrors.IsNotFound(err):
				d.SetId("")
				return nil
			default:
				return diag.FromErr(err)
			}
		}
	case hasImage:
		fieldSelector := map[string]string{
			"spec.image": image.(string),
		}
		if hasTag {
			fieldSelector["spec.tag"] = tag.(string)
		}
		list, err := clientSet.ContainerV1().Images(branch).List(ctx, metav1.ListOptions{FieldSelector: fieldSelector})
		if err != nil {
			return diag.FromErr(err)
		}
		if len(list.Items) == 0 {
			d.SetId("")
			return nil
		}
		latestImg := slices.MaxFunc(list.Items, func(a, b containerv1.Image) int {
			switch {
			case a.CreatedTimestamp.Before(b.CreatedTimestamp):
				return -1
			case a.CreatedTimestamp.Equal(b.CreatedTimestamp):
				return 0
			}
			return 1
		})
		obj = &latestImg
	default:
		return diag.FromErr(errors.New("either metadata.0.name or spec.0.image is required"))
	}

	d.SetId(obj.Branch + "/" + obj.Name)

	data, err := ec.Converter().Flatten(obj, imageSchema())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = resource.SetData(d, data); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
