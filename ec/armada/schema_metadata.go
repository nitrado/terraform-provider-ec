package armada

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func metadataSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"annotations": {
			Type:        schema.TypeMap,
			Description: "An unstructured map of keys and values stored on an object.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"labels": {
			Type:        schema.TypeMap,
			Description: "A map of keys and values that can be used to organize and categorize objects.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"name": {
			Type:             schema.TypeString,
			Description:      "The unique object name within its scope.",
			Optional:         true,
			ForceNew:         true,
			Computed:         true,
			ValidateDiagFunc: validateName,
		},
		"revision": {
			Type:        schema.TypeString,
			Description: "An opaque resource revision.",
			Computed:    true,
		},
		"uid": {
			Type:        schema.TypeString,
			Description: "A unique identifier for each an object.",
			Computed:    true,
		},
	}
}
