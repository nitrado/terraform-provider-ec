package formation

// Code generated by schema-gen. DO NOT EDIT.

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec/meta"
)

func vesselSchema() map[string]*schema.Schema {
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
			Description: "Spec defines the desired vessel.",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"description": {
						Type:        schema.TypeString,
						Description: "Description is the optional description of the vessel.",
						Optional:    true,
					},
					"region": {
						Type:        schema.TypeString,
						Description: "Region defines the region the vessel is deployed to.",
						Required:    true,
					},
					"suspend": {
						Type:        schema.TypeList,
						Description: "Suspend specifies whether the vessel should create a game server or not.",
						Optional:    true,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"value": {
									Type:        schema.TypeBool,
									Description: "Suspend specifies whether the vessel should create a game server or not.",
									Required:    true,
								},
							},
						},
					},
					"template": {
						Type:        schema.TypeList,
						Description: "Template describes the game server that is created.",
						Required:    true,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"metadata": {
									Type:        schema.TypeList,
									Description: "Standard object's metadata.",
									Optional:    true,
									MaxItems:    1,
									Elem:        &schema.Resource{Schema: meta.Schema()},
								},
								"spec": {
									Type:        schema.TypeList,
									Description: "Spec defines the desired game server.",
									Optional:    true,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"containers": {
												Type:        schema.TypeList,
												Description: "Containers is a list of container belonging to the game server.",
												Required:    true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"args": {
															Type:        schema.TypeList,
															Description: "Args are arguments to the entrypoint.",
															Optional:    true,
															Elem:        &schema.Schema{Type: schema.TypeString},
														},
														"branch": {
															Type:        schema.TypeString,
															Description: "Branch is the name of the image branch.",
															Required:    true,
														},
														"command": {
															Type:        schema.TypeList,
															Description: "Command is the entrypoint array. This is not executed within a shell.",
															Optional:    true,
															Elem:        &schema.Schema{Type: schema.TypeString},
														},
														"config_files": {
															Type:        schema.TypeList,
															Description: "ConfigFiles is a list of configuration files to mount into the containers filesystem.",
															Optional:    true,
															Elem: &schema.Resource{
																Schema: map[string]*schema.Schema{
																	"mount_path": {
																		Type:        schema.TypeString,
																		Description: "MountPath is the path to mount the configuration file on.",
																		Required:    true,
																	},
																	"name": {
																		Type:        schema.TypeString,
																		Description: "Name is the name of the configuration file.",
																		Required:    true,
																	},
																},
															},
														},
														"env": {
															Type:        schema.TypeList,
															Description: "Env is a list of environment variables to set in the container.",
															Optional:    true,
															Elem:        &schema.Resource{Schema: envSchema()},
														},
														"image": {
															Type:        schema.TypeString,
															Description: "Image is a reference to the containerv1.Image to deploy in this container.",
															Required:    true,
														},
														"name": {
															Type:        schema.TypeString,
															Description: "Name is the name of the container.",
															Required:    true,
														},
														"ports": {
															Type:        schema.TypeList,
															Description: "Ports are the ports to expose from the container.",
															Optional:    true,
															Elem: &schema.Resource{
																Schema: map[string]*schema.Schema{
																	"container_port": {
																		Type:        schema.TypeInt,
																		Description: "ContainerPort is the port that is being opened on the specified container's process.",
																		Optional:    true,
																	},
																	"name": {
																		Type:        schema.TypeString,
																		Description: "Name is the name of the port.",
																		Required:    true,
																	},
																	"policy": {
																		Type:        schema.TypeString,
																		Description: "Policy defines the policy for how the HostPort is populated.",
																		Required:    true,
																	},
																	"protection_protocol": {
																		Type:        schema.TypeList,
																		Description: "ProtectionProtocol is the optional name of the protection protocol being used.",
																		Optional:    true,
																		MaxItems:    1,
																		Elem: &schema.Resource{
																			Schema: map[string]*schema.Schema{
																				"value": {
																					Type:        schema.TypeString,
																					Description: "ProtectionProtocol is the optional name of the protection protocol being used.",
																					Required:    true,
																				},
																			},
																		},
																	},
																	"protocol": {
																		Type:        schema.TypeString,
																		Description: "Protocol is the network protocol being used. Defaults to UDP. TCP and TCPUDP are other options.",
																		Optional:    true,
																	},
																},
															},
														},
														"resources": {
															Type:        schema.TypeList,
															Description: "Resources are the compute resources required by the container.",
															Optional:    true,
															MaxItems:    1,
															Elem: &schema.Resource{
																Schema: map[string]*schema.Schema{
																	"claims": {
																		Type:     schema.TypeList,
																		Optional: true,
																		Elem: &schema.Resource{
																			Schema: map[string]*schema.Schema{
																				"name": {
																					Type:     schema.TypeString,
																					Optional: true,
																				},
																				"request": {
																					Type:     schema.TypeString,
																					Optional: true,
																				},
																			},
																		},
																	},
																	"limits": {
																		Type:     schema.TypeMap,
																		Optional: true,
																		Elem:     &schema.Schema{Type: schema.TypeString},
																	},
																	"requests": {
																		Type:     schema.TypeMap,
																		Optional: true,
																		Elem:     &schema.Schema{Type: schema.TypeString},
																	},
																},
															},
														},
														"security_context": {
															Type:        schema.TypeList,
															Description: "SecurityContext defines the security options the container should be run with.",
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
																	"app_armor_profile": {
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
														"volume_mounts": {
															Type:        schema.TypeList,
															Description: "VolumeMounts are the volumes to mount into the container's filesystem.",
															Optional:    true,
															Elem: &schema.Resource{
																Schema: map[string]*schema.Schema{
																	"mount_path": {
																		Type:     schema.TypeString,
																		Optional: true,
																	},
																	"mount_propagation": {
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
																	"name": {
																		Type:     schema.TypeString,
																		Optional: true,
																	},
																	"read_only": {
																		Type:     schema.TypeBool,
																		Optional: true,
																	},
																	"recursive_read_only": {
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
																	"sub_path": {
																		Type:     schema.TypeString,
																		Optional: true,
																	},
																	"sub_path_expr": {
																		Type:     schema.TypeString,
																		Optional: true,
																	},
																},
															},
														},
													},
												},
											},
											"gateway_policies": {
												Type:        schema.TypeList,
												Description: "GatewayPolicies are the gateway policy names applied to the game servers.",
												Optional:    true,
												Elem:        &schema.Schema{Type: schema.TypeString},
											},
											"health": {
												Type:        schema.TypeList,
												Description: "Health is the health checking configuration for Agones game servers.",
												Optional:    true,
												MaxItems:    1,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"disabled": {
															Type:     schema.TypeBool,
															Optional: true,
														},
														"failure_threshold": {
															Type:     schema.TypeInt,
															Optional: true,
														},
														"initial_delay_seconds": {
															Type:     schema.TypeInt,
															Optional: true,
														},
														"period_seconds": {
															Type:     schema.TypeInt,
															Optional: true,
														},
													},
												},
											},
											"termination_grace_period_seconds": {
												Type:        schema.TypeList,
												Description: "TerminationGracePeriodSeconds is the optional duration in seconds the game servers need to terminate gracefully. Defaults to 30 seconds.",
												Optional:    true,
												MaxItems:    1,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"value": {
															Type:        schema.TypeInt,
															Description: "TerminationGracePeriodSeconds is the optional duration in seconds the game servers need to terminate gracefully. Defaults to 30 seconds.",
															Required:    true,
														},
													},
												},
											},
											"volumes": {
												Type:        schema.TypeList,
												Description: "Volumes are pod volumes.",
												Optional:    true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"empty_dir": {
															Type:        schema.TypeList,
															Description: "EmptyDir configures an empty dir volume.",
															Optional:    true,
															MaxItems:    1,
															Elem: &schema.Resource{
																Schema: map[string]*schema.Schema{
																	"medium": {
																		Type:        schema.TypeString,
																		Description: "Medium is the storage medium type.",
																		Optional:    true,
																	},
																	"size_limit": {
																		Type:        schema.TypeList,
																		Description: "SizeLimit is the maximum size of the volume.",
																		Optional:    true,
																		MaxItems:    1,
																		Elem: &schema.Resource{
																			Schema: map[string]*schema.Schema{
																				"value": {
																					Type:        schema.TypeString,
																					Description: "SizeLimit is the maximum size of the volume.",
																					Required:    true,
																				},
																			},
																		},
																	},
																},
															},
														},
														"name": {
															Type:        schema.TypeString,
															Description: "Name is the name of the volume mount.",
															Required:    true,
														},
														"persistent": {
															Type:        schema.TypeList,
															Description: "Persistent configures a persistent volume.",
															Optional:    true,
															MaxItems:    1,
															Elem: &schema.Resource{
																Schema: map[string]*schema.Schema{
																	"volume_name": {
																		Type:        schema.TypeString,
																		Description: "VolumeName is the name of the volume to store data in.",
																		Required:    true,
																	},
																},
															},
														},
														"type": {
															Type:        schema.TypeString,
															Description: "Type is the volume type.",
															Optional:    true,
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
					"termination_grace_periods": {
						Type:        schema.TypeList,
						Description: "TerminationGracePeriods are the optional durations that a game server has to terminate gracefully. If this value is nil, the default grace period for each situation will be used. These durations only apply when a game server is in use.",
						Optional:    true,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"maintenance": {
									Type:        schema.TypeList,
									Description: "Maintenance is the optional duration in seconds that a game server has to gracefully terminate when the site it is running is cordoned.",
									Optional:    true,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"value": {
												Type:        schema.TypeInt,
												Description: "Maintenance is the optional duration in seconds that a game server has to gracefully terminate when the site it is running is cordoned.",
												Required:    true,
											},
										},
									},
								},
								"spec_change": {
									Type:        schema.TypeList,
									Description: "SpecChange is the optional duration in seconds that a game server has to gracefully terminate when a spec change is detected.",
									Optional:    true,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"value": {
												Type:        schema.TypeInt,
												Description: "SpecChange is the optional duration in seconds that a game server has to gracefully terminate when a spec change is detected.",
												Required:    true,
											},
										},
									},
								},
								"user_initiated": {
									Type:        schema.TypeList,
									Description: "UserInitiated is the optional duration in seconds that a game server has to gracefully terminate when user initiates a restart or suspends a vessel.",
									Optional:    true,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"value": {
												Type:        schema.TypeInt,
												Description: "UserInitiated is the optional duration in seconds that a game server has to gracefully terminate when user initiates a restart or suspends a vessel.",
												Required:    true,
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
