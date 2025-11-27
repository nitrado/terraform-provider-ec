package core

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec/meta"
)

// Keep in sync with schema_conigfile.go.
func configFileSchemaDS() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance": {
			Type:        schema.TypeString,
			Description: "Name is an instance name configured in the provider.",
			Optional:    true,
		},
		"data": {
			Type:        schema.TypeString,
			Description: "Data is the content of the configuration file.",
			Optional:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "Description is the optional description of the config file.",
			Optional:    true,
		},
		"metadata": {
			Type:        schema.TypeList,
			Description: "Standard object's metadata.",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: meta.Schema()},
		},
	}
}
