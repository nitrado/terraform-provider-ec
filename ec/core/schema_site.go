package core

// Code generated by schema-gen. DO NOT EDIT.

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec/meta"
)

func siteSchema() map[string]*schema.Schema {
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
			Description: "Spec defines the desired site configuration.",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"cordoned": {
						Type:        schema.TypeBool,
						Description: "Cordoned determines if a site can have fleets scheduled.",
						Optional:    true,
					},
					"cpu_ratio": {
						Type:        schema.TypeFloat,
						Description: "CPURatio is the applied ratio for any subordinate game server CPU request or limit.  This facilitates the optimal utilization of various CPU generations for a game. The default is 1.0.",
						Optional:    true,
					},
					"credentials": {
						Type:        schema.TypeList,
						Description: "Credentials are the credentials used to access the site.",
						Required:    true,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"certificate": {
									Type:        schema.TypeString,
									Description: "Certificate is the CA certificate of the cluster.",
									Required:    true,
								},
								"endpoint": {
									Type:        schema.TypeString,
									Description: "Endpoint is the address of the kubernetes cluster.",
									Required:    true,
								},
								"namespace": {
									Type:        schema.TypeString,
									Description: "Namespace is the cluster namespace assigned to the site.",
									Required:    true,
								},
								"token": {
									Type:        schema.TypeString,
									Description: "Token is the authorization token.",
									Required:    true,
								},
							},
						},
					},
					"description": {
						Type:        schema.TypeString,
						Description: "Description is the optional description of the site.",
						Optional:    true,
					},
					"resources": {
						Type:        schema.TypeList,
						Description: "Resources defines the resource limit of the site.",
						Required:    true,
						MaxItems:    1,
						Elem:        &schema.Resource{Schema: resourcesSchema()},
					},
					"template": {
						Type:        schema.TypeList,
						Description: "Template is the optional configuration to apply to all fleets on this site.",
						Optional:    true,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"affinity": {
									Type:        schema.TypeList,
									Description: "Affinity is a group of affinity scheduling rules.",
									Optional:    true,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"node_affinity": {
												Type:     schema.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"preferred_during_scheduling_ignored_during_execution": {
															Type:     schema.TypeList,
															Optional: true,
															Elem: &schema.Resource{
																Schema: map[string]*schema.Schema{
																	"preference": {
																		Type:     schema.TypeList,
																		Optional: true,
																		MaxItems: 1,
																		Elem: &schema.Resource{
																			Schema: map[string]*schema.Schema{
																				"match_expressions": {
																					Type:     schema.TypeList,
																					Optional: true,
																					Elem: &schema.Resource{
																						Schema: map[string]*schema.Schema{
																							"key": {
																								Type:     schema.TypeString,
																								Optional: true,
																							},
																							"operator": {
																								Type:     schema.TypeString,
																								Optional: true,
																							},
																							"values": {
																								Type:     schema.TypeList,
																								Optional: true,
																								Elem:     &schema.Schema{Type: schema.TypeString},
																							},
																						},
																					},
																				},
																				"match_fields": {
																					Type:     schema.TypeList,
																					Optional: true,
																					Elem: &schema.Resource{
																						Schema: map[string]*schema.Schema{
																							"key": {
																								Type:     schema.TypeString,
																								Optional: true,
																							},
																							"operator": {
																								Type:     schema.TypeString,
																								Optional: true,
																							},
																							"values": {
																								Type:     schema.TypeList,
																								Optional: true,
																								Elem:     &schema.Schema{Type: schema.TypeString},
																							},
																						},
																					},
																				},
																			},
																		},
																	},
																	"weight": {
																		Type:     schema.TypeInt,
																		Optional: true,
																	},
																},
															},
														},
														"required_during_scheduling_ignored_during_execution": {
															Type:     schema.TypeList,
															Optional: true,
															MaxItems: 1,
															Elem: &schema.Resource{
																Schema: map[string]*schema.Schema{
																	"node_selector_terms": {
																		Type:     schema.TypeList,
																		Optional: true,
																		Elem: &schema.Resource{
																			Schema: map[string]*schema.Schema{
																				"match_expressions": {
																					Type:     schema.TypeList,
																					Optional: true,
																					Elem: &schema.Resource{
																						Schema: map[string]*schema.Schema{
																							"key": {
																								Type:     schema.TypeString,
																								Optional: true,
																							},
																							"operator": {
																								Type:     schema.TypeString,
																								Optional: true,
																							},
																							"values": {
																								Type:     schema.TypeList,
																								Optional: true,
																								Elem:     &schema.Schema{Type: schema.TypeString},
																							},
																						},
																					},
																				},
																				"match_fields": {
																					Type:     schema.TypeList,
																					Optional: true,
																					Elem: &schema.Resource{
																						Schema: map[string]*schema.Schema{
																							"key": {
																								Type:     schema.TypeString,
																								Optional: true,
																							},
																							"operator": {
																								Type:     schema.TypeString,
																								Optional: true,
																							},
																							"values": {
																								Type:     schema.TypeList,
																								Optional: true,
																								Elem:     &schema.Schema{Type: schema.TypeString},
																							},
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
											"pod_affinity": {
												Type:     schema.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"preferred_during_scheduling_ignored_during_execution": {
															Type:     schema.TypeList,
															Optional: true,
															Elem: &schema.Resource{
																Schema: map[string]*schema.Schema{
																	"pod_affinity_term": {
																		Type:     schema.TypeList,
																		Optional: true,
																		MaxItems: 1,
																		Elem: &schema.Resource{
																			Schema: map[string]*schema.Schema{
																				"label_selector": {
																					Type:     schema.TypeList,
																					Optional: true,
																					MaxItems: 1,
																					Elem: &schema.Resource{
																						Schema: map[string]*schema.Schema{
																							"match_expressions": {
																								Type:     schema.TypeList,
																								Optional: true,
																								Elem: &schema.Resource{
																									Schema: map[string]*schema.Schema{
																										"key": {
																											Type:     schema.TypeString,
																											Optional: true,
																										},
																										"operator": {
																											Type:     schema.TypeString,
																											Optional: true,
																										},
																										"values": {
																											Type:     schema.TypeList,
																											Optional: true,
																											Elem:     &schema.Schema{Type: schema.TypeString},
																										},
																									},
																								},
																							},
																							"match_labels": {
																								Type:     schema.TypeMap,
																								Optional: true,
																								Elem:     &schema.Schema{Type: schema.TypeString},
																							},
																						},
																					},
																				},
																				"namespace_selector": {
																					Type:     schema.TypeList,
																					Optional: true,
																					MaxItems: 1,
																					Elem: &schema.Resource{
																						Schema: map[string]*schema.Schema{
																							"match_expressions": {
																								Type:     schema.TypeList,
																								Optional: true,
																								Elem: &schema.Resource{
																									Schema: map[string]*schema.Schema{
																										"key": {
																											Type:     schema.TypeString,
																											Optional: true,
																										},
																										"operator": {
																											Type:     schema.TypeString,
																											Optional: true,
																										},
																										"values": {
																											Type:     schema.TypeList,
																											Optional: true,
																											Elem:     &schema.Schema{Type: schema.TypeString},
																										},
																									},
																								},
																							},
																							"match_labels": {
																								Type:     schema.TypeMap,
																								Optional: true,
																								Elem:     &schema.Schema{Type: schema.TypeString},
																							},
																						},
																					},
																				},
																				"namespaces": {
																					Type:     schema.TypeList,
																					Optional: true,
																					Elem:     &schema.Schema{Type: schema.TypeString},
																				},
																				"topology_key": {
																					Type:     schema.TypeString,
																					Optional: true,
																				},
																			},
																		},
																	},
																	"weight": {
																		Type:     schema.TypeInt,
																		Optional: true,
																	},
																},
															},
														},
														"required_during_scheduling_ignored_during_execution": {
															Type:     schema.TypeList,
															Optional: true,
															Elem: &schema.Resource{
																Schema: map[string]*schema.Schema{
																	"label_selector": {
																		Type:     schema.TypeList,
																		Optional: true,
																		MaxItems: 1,
																		Elem: &schema.Resource{
																			Schema: map[string]*schema.Schema{
																				"match_expressions": {
																					Type:     schema.TypeList,
																					Optional: true,
																					Elem: &schema.Resource{
																						Schema: map[string]*schema.Schema{
																							"key": {
																								Type:     schema.TypeString,
																								Optional: true,
																							},
																							"operator": {
																								Type:     schema.TypeString,
																								Optional: true,
																							},
																							"values": {
																								Type:     schema.TypeList,
																								Optional: true,
																								Elem:     &schema.Schema{Type: schema.TypeString},
																							},
																						},
																					},
																				},
																				"match_labels": {
																					Type:     schema.TypeMap,
																					Optional: true,
																					Elem:     &schema.Schema{Type: schema.TypeString},
																				},
																			},
																		},
																	},
																	"namespace_selector": {
																		Type:     schema.TypeList,
																		Optional: true,
																		MaxItems: 1,
																		Elem: &schema.Resource{
																			Schema: map[string]*schema.Schema{
																				"match_expressions": {
																					Type:     schema.TypeList,
																					Optional: true,
																					Elem: &schema.Resource{
																						Schema: map[string]*schema.Schema{
																							"key": {
																								Type:     schema.TypeString,
																								Optional: true,
																							},
																							"operator": {
																								Type:     schema.TypeString,
																								Optional: true,
																							},
																							"values": {
																								Type:     schema.TypeList,
																								Optional: true,
																								Elem:     &schema.Schema{Type: schema.TypeString},
																							},
																						},
																					},
																				},
																				"match_labels": {
																					Type:     schema.TypeMap,
																					Optional: true,
																					Elem:     &schema.Schema{Type: schema.TypeString},
																				},
																			},
																		},
																	},
																	"namespaces": {
																		Type:     schema.TypeList,
																		Optional: true,
																		Elem:     &schema.Schema{Type: schema.TypeString},
																	},
																	"topology_key": {
																		Type:     schema.TypeString,
																		Optional: true,
																	},
																},
															},
														},
													},
												},
											},
											"pod_anti_affinity": {
												Type:     schema.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"preferred_during_scheduling_ignored_during_execution": {
															Type:     schema.TypeList,
															Optional: true,
															Elem: &schema.Resource{
																Schema: map[string]*schema.Schema{
																	"pod_affinity_term": {
																		Type:     schema.TypeList,
																		Optional: true,
																		MaxItems: 1,
																		Elem: &schema.Resource{
																			Schema: map[string]*schema.Schema{
																				"label_selector": {
																					Type:     schema.TypeList,
																					Optional: true,
																					MaxItems: 1,
																					Elem: &schema.Resource{
																						Schema: map[string]*schema.Schema{
																							"match_expressions": {
																								Type:     schema.TypeList,
																								Optional: true,
																								Elem: &schema.Resource{
																									Schema: map[string]*schema.Schema{
																										"key": {
																											Type:     schema.TypeString,
																											Optional: true,
																										},
																										"operator": {
																											Type:     schema.TypeString,
																											Optional: true,
																										},
																										"values": {
																											Type:     schema.TypeList,
																											Optional: true,
																											Elem:     &schema.Schema{Type: schema.TypeString},
																										},
																									},
																								},
																							},
																							"match_labels": {
																								Type:     schema.TypeMap,
																								Optional: true,
																								Elem:     &schema.Schema{Type: schema.TypeString},
																							},
																						},
																					},
																				},
																				"namespace_selector": {
																					Type:     schema.TypeList,
																					Optional: true,
																					MaxItems: 1,
																					Elem: &schema.Resource{
																						Schema: map[string]*schema.Schema{
																							"match_expressions": {
																								Type:     schema.TypeList,
																								Optional: true,
																								Elem: &schema.Resource{
																									Schema: map[string]*schema.Schema{
																										"key": {
																											Type:     schema.TypeString,
																											Optional: true,
																										},
																										"operator": {
																											Type:     schema.TypeString,
																											Optional: true,
																										},
																										"values": {
																											Type:     schema.TypeList,
																											Optional: true,
																											Elem:     &schema.Schema{Type: schema.TypeString},
																										},
																									},
																								},
																							},
																							"match_labels": {
																								Type:     schema.TypeMap,
																								Optional: true,
																								Elem:     &schema.Schema{Type: schema.TypeString},
																							},
																						},
																					},
																				},
																				"namespaces": {
																					Type:     schema.TypeList,
																					Optional: true,
																					Elem:     &schema.Schema{Type: schema.TypeString},
																				},
																				"topology_key": {
																					Type:     schema.TypeString,
																					Optional: true,
																				},
																			},
																		},
																	},
																	"weight": {
																		Type:     schema.TypeInt,
																		Optional: true,
																	},
																},
															},
														},
														"required_during_scheduling_ignored_during_execution": {
															Type:     schema.TypeList,
															Optional: true,
															Elem: &schema.Resource{
																Schema: map[string]*schema.Schema{
																	"label_selector": {
																		Type:     schema.TypeList,
																		Optional: true,
																		MaxItems: 1,
																		Elem: &schema.Resource{
																			Schema: map[string]*schema.Schema{
																				"match_expressions": {
																					Type:     schema.TypeList,
																					Optional: true,
																					Elem: &schema.Resource{
																						Schema: map[string]*schema.Schema{
																							"key": {
																								Type:     schema.TypeString,
																								Optional: true,
																							},
																							"operator": {
																								Type:     schema.TypeString,
																								Optional: true,
																							},
																							"values": {
																								Type:     schema.TypeList,
																								Optional: true,
																								Elem:     &schema.Schema{Type: schema.TypeString},
																							},
																						},
																					},
																				},
																				"match_labels": {
																					Type:     schema.TypeMap,
																					Optional: true,
																					Elem:     &schema.Schema{Type: schema.TypeString},
																				},
																			},
																		},
																	},
																	"namespace_selector": {
																		Type:     schema.TypeList,
																		Optional: true,
																		MaxItems: 1,
																		Elem: &schema.Resource{
																			Schema: map[string]*schema.Schema{
																				"match_expressions": {
																					Type:     schema.TypeList,
																					Optional: true,
																					Elem: &schema.Resource{
																						Schema: map[string]*schema.Schema{
																							"key": {
																								Type:     schema.TypeString,
																								Optional: true,
																							},
																							"operator": {
																								Type:     schema.TypeString,
																								Optional: true,
																							},
																							"values": {
																								Type:     schema.TypeList,
																								Optional: true,
																								Elem:     &schema.Schema{Type: schema.TypeString},
																							},
																						},
																					},
																				},
																				"match_labels": {
																					Type:     schema.TypeMap,
																					Optional: true,
																					Elem:     &schema.Schema{Type: schema.TypeString},
																				},
																			},
																		},
																	},
																	"namespaces": {
																		Type:     schema.TypeList,
																		Optional: true,
																		Elem:     &schema.Schema{Type: schema.TypeString},
																	},
																	"topology_key": {
																		Type:     schema.TypeString,
																		Optional: true,
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
								"env": {
									Type:        schema.TypeList,
									Description: "Env is a list of environment variables to set on all containers on this site.",
									Optional:    true,
									Elem:        &schema.Resource{Schema: envSchema()},
								},
								"image_pull_secrets": {
									Type:        schema.TypeList,
									Description: "ImagePullSecrets points to secrets with authorization tokens that store docker credentials to access a registry.",
									Optional:    true,
									Elem:        &schema.Schema{Type: schema.TypeString},
								},
								"security_context": {
									Type:        schema.TypeList,
									Description: "SecurityContext defines the security options the container should be run with. This security context overrides the user security context if a top level property is set.",
									Optional:    true,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"allow_privilege_escalation": {
												Type:     schema.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"value": {
															Type:     schema.TypeBool,
															Required: true,
														},
													},
												},
											},
											"capabilities": {
												Type:     schema.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"add": {
															Type:     schema.TypeList,
															Optional: true,
															Elem:     &schema.Schema{Type: schema.TypeString},
														},
														"drop": {
															Type:     schema.TypeList,
															Optional: true,
															Elem:     &schema.Schema{Type: schema.TypeString},
														},
													},
												},
											},
											"privileged": {
												Type:     schema.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"value": {
															Type:     schema.TypeBool,
															Required: true,
														},
													},
												},
											},
											"proc_mount": {
												Type:     schema.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"value": {
															Type:     schema.TypeString,
															Required: true,
														},
													},
												},
											},
											"read_only_root_filesystem": {
												Type:     schema.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"value": {
															Type:     schema.TypeBool,
															Required: true,
														},
													},
												},
											},
											"run_as_group": {
												Type:     schema.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"value": {
															Type:     schema.TypeInt,
															Required: true,
														},
													},
												},
											},
											"run_as_non_root": {
												Type:     schema.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"value": {
															Type:     schema.TypeBool,
															Required: true,
														},
													},
												},
											},
											"run_as_user": {
												Type:     schema.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"value": {
															Type:     schema.TypeInt,
															Required: true,
														},
													},
												},
											},
											"se_linux_options": {
												Type:     schema.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"level": {
															Type:     schema.TypeString,
															Optional: true,
														},
														"role": {
															Type:     schema.TypeString,
															Optional: true,
														},
														"type": {
															Type:     schema.TypeString,
															Optional: true,
														},
														"user": {
															Type:     schema.TypeString,
															Optional: true,
														},
													},
												},
											},
											"seccomp_profile": {
												Type:     schema.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"localhost_profile": {
															Type:     schema.TypeList,
															Optional: true,
															MaxItems: 1,
															Elem: &schema.Resource{
																Schema: map[string]*schema.Schema{
																	"value": {
																		Type:     schema.TypeString,
																		Required: true,
																	},
																},
															},
														},
														"type": {
															Type:     schema.TypeString,
															Optional: true,
														},
													},
												},
											},
											"windows_options": {
												Type:     schema.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"gmsa_credential_spec": {
															Type:     schema.TypeList,
															Optional: true,
															MaxItems: 1,
															Elem: &schema.Resource{
																Schema: map[string]*schema.Schema{
																	"value": {
																		Type:     schema.TypeString,
																		Required: true,
																	},
																},
															},
														},
														"gmsa_credential_spec_name": {
															Type:     schema.TypeList,
															Optional: true,
															MaxItems: 1,
															Elem: &schema.Resource{
																Schema: map[string]*schema.Schema{
																	"value": {
																		Type:     schema.TypeString,
																		Required: true,
																	},
																},
															},
														},
														"host_process": {
															Type:     schema.TypeList,
															Optional: true,
															MaxItems: 1,
															Elem: &schema.Resource{
																Schema: map[string]*schema.Schema{
																	"value": {
																		Type:     schema.TypeBool,
																		Required: true,
																	},
																},
															},
														},
														"run_as_user_name": {
															Type:     schema.TypeList,
															Optional: true,
															MaxItems: 1,
															Elem: &schema.Resource{
																Schema: map[string]*schema.Schema{
																	"value": {
																		Type:     schema.TypeString,
																		Required: true,
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
								"tolerations": {
									Type:        schema.TypeList,
									Description: "Tolerations is a set of pod tolerations.",
									Optional:    true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"effect": {
												Type:     schema.TypeString,
												Optional: true,
											},
											"key": {
												Type:     schema.TypeString,
												Optional: true,
											},
											"operator": {
												Type:     schema.TypeString,
												Optional: true,
											},
											"toleration_seconds": {
												Type:     schema.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"value": {
															Type:     schema.TypeInt,
															Required: true,
														},
													},
												},
											},
											"value": {
												Type:     schema.TypeString,
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
