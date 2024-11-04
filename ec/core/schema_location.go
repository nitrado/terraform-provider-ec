package core

// Code generated by schema-gen. DO NOT EDIT.

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec/meta"
)

func locationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance": {
			Type:        schema.TypeString,
			Description: "Name is an instance name configured in the provider.",
			Optional:    true,
		},
		"metadata": {
			Type:        schema.TypeList,
			Description: "Standard object's metadata.",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: meta.Schema()},
		},
		"spec": {
			Type:        schema.TypeList,
			Description: "Spec defines the desired location configuration.",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"region": {
						Type:        schema.TypeString,
						Description: "Region is the location region.",
						Optional:    true,
					},
					"sites": {
						Type:        schema.TypeList,
						Description: "Sites are the site names that makes up the location.",
						Optional:    true,
						Elem:        &schema.Schema{Type: schema.TypeString},
					},
					"type": {
						Type:        schema.TypeString,
						Description: "Type is the type of the datacenter, e.g. cloud or bare metal.",
						Optional:    true,
					},
				},
			},
		},
	}
}
