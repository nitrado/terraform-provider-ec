package core

// Code generated by schema-gen. DO NOT EDIT.

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func envSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Name is the name of the environment variable.",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "Value is the value of the environment variable.",
			Optional:    true,
		},
		"value_from": {
			Type:        schema.TypeList,
			Description: "ValueFrom is the source for the environment variable's value.",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"config_file_key_ref": {
						Type:        schema.TypeList,
						Description: "ConfigFileKeyRef select the configuration file.",
						Optional:    true,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:        schema.TypeString,
									Description: "Name is the name of the configuration file.",
									Required:    true,
								},
							},
						},
					},
					"field_ref": {
						Type:        schema.TypeList,
						Description: "FieldRef selects the field of the pod. Supports metadata.name, metadata.namespace, `metadata.labels['<KEY>']`, `metadata.annotations['<KEY>']`, metadata.armadaName, metadata.regionName, metadata.regionTypeName, metadata.siteName, metadata.imageBranch, metadata.imageName, metadata.imageTag, spec.nodeName, spec.serviceAccountName, status.hostIP, status.podIP, status.podIPs.",
						Optional:    true,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"api_version": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"field_path": {
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
					},
				},
			},
		},
	}
}
