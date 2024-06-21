package container

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec/meta"
)

// imageSchema is manually created, due to specific requirements.
func imageSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance": {
			Type:        schema.TypeString,
			Description: "Name is an instance name configured in the provider.",
			Optional:    true,
		},
		"metadata": {
			Type:        schema.TypeList,
			Description: "Image object metadata.",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"branch": {
						Type:        schema.TypeString,
						Description: "Branch defines the branch within which each name must be unique.",
						Required:    true,
					},
					"object_meta": {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem:     &schema.Resource{Schema: meta.Schema()},
					},
				},
			},
		},
		"spec": {
			Type:        schema.TypeList,
			Description: "Spec defines the desired image.",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"digest": {
						Type:        schema.TypeString,
						Description: "Digest is the hash of the image.",
						Computed:    true,
					},
					"image": {
						Type:        schema.TypeString,
						Description: "Image is the image name.",
						Required:    true,
					},
					"registry": {
						Type:        schema.TypeString,
						Description: "Registry is the registry that contains the image.",
						Computed:    true,
					},
					"tag": {
						Type:        schema.TypeString,
						Description: "Tag is the image tag.",
						Optional:    true,
						Computed:    true,
					},
				},
			},
		},
	}
}
