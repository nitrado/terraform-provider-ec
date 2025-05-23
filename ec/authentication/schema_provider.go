package authentication

// Code generated by schema-gen. DO NOT EDIT.

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec/meta"
)

func providerSchema() map[string]*schema.Schema {
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
			Description: "Spec defines the desired state of the authentication provider.",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"display_name": {
						Type:        schema.TypeString,
						Description: "DisplayName is the display name of the provider.",
						Optional:    true,
					},
					"oidc": {
						Type:        schema.TypeList,
						Description: "OIDC is the OIDC provider configuration.",
						Optional:    true,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"acr_values": {
									Type:        schema.TypeList,
									Description: "ACRValues (Authentication Context Class Reference Values) that specifies the Authentication Context Class Values within the Authentication Request that the Authorization Server is being requested to use for processing requests from this Client, with the values appearing in order of preference.",
									Optional:    true,
									Elem:        &schema.Schema{Type: schema.TypeString},
								},
								"allowed_groups": {
									Type:        schema.TypeList,
									Description: "AllowedGroups is a list of groups that are allowed to authenticate with this provider.",
									Optional:    true,
									Elem:        &schema.Schema{Type: schema.TypeString},
								},
								"basic_auth_unsupported": {
									Type:        schema.TypeList,
									Description: "BasicAuthUnsupported causes client_secret to be passed as POST parameters instead of basic auth. This is specifically \"NOT RECOMMENDED\" by the OAuth2 RFC, but some providers require it.  https://tools.ietf.org/html/rfc6749#section-2.3.1",
									Optional:    true,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"value": {
												Type:        schema.TypeBool,
												Description: "BasicAuthUnsupported causes client_secret to be passed as POST parameters instead of basic auth. This is specifically \"NOT RECOMMENDED\" by the OAuth2 RFC, but some providers require it.  https://tools.ietf.org/html/rfc6749#section-2.3.1",
												Required:    true,
											},
										},
									},
								},
								"claim_mapping": {
									Type:        schema.TypeList,
									Description: "ClaimMapping contains all claim mapping overrides.",
									Optional:    true,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"email": {
												Type:        schema.TypeString,
												Description: "Configurable key which contains the email claims. Defaults to \"email\".",
												Optional:    true,
											},
											"groups": {
												Type:        schema.TypeString,
												Description: "Configurable key which contains the groups claims Defaults to \"groups\".",
												Optional:    true,
											},
											"preferred_username": {
												Type:        schema.TypeString,
												Description: "Configurable key which contains the preferred username claims Defaults to \"preferred_username\".",
												Optional:    true,
											},
										},
									},
								},
								"claim_modifications": {
									Type:        schema.TypeList,
									Description: "ClaimMutations contains all claim mutations options.",
									Optional:    true,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"filter_group_claims": {
												Type:        schema.TypeList,
												Description: "FilterGroupClaims is a regex filter used to keep only the matching groups.",
												Optional:    true,
												MaxItems:    1,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"groups_filter": {
															Type:        schema.TypeString,
															Description: "GroupsFilter is the regex filter used to keep only the matching groups.",
															Optional:    true,
														},
													},
												},
											},
											"new_group_from_claims": {
												Type:        schema.TypeList,
												Description: "NewGroupFromClaims specifies how claims can be joined to create groups.",
												Optional:    true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"claims": {
															Type:        schema.TypeList,
															Description: "Claims is a list of claim to join together",
															Optional:    true,
															Elem:        &schema.Schema{Type: schema.TypeString},
														},
														"clear_delimiter": {
															Type:        schema.TypeBool,
															Description: "ClearDelimiter indicates if the Delimiter string should be removed from claim values.",
															Optional:    true,
														},
														"delimiter": {
															Type:        schema.TypeString,
															Description: "Delimiter is the string used to separate the claims.",
															Optional:    true,
														},
														"prefix": {
															Type:        schema.TypeString,
															Description: "Prefix is a string to place before the first claim.",
															Optional:    true,
														},
													},
												},
											},
										},
									},
								},
								"client_id": {
									Type:        schema.TypeString,
									Description: "ClientID is the client ID of the OIDC provider.",
									Required:    true,
								},
								"client_secret": {
									Type:        schema.TypeString,
									Description: "ClientSecret is the client secret of the OIDC provider.",
									Required:    true,
								},
								"get_user_info": {
									Type:        schema.TypeBool,
									Description: "GetUserInfo uses the userinfo endpoint to get additional claims for the token. This is especially useful where upstreams return \"thin\" id tokens.",
									Optional:    true,
								},
								"hosted_domains": {
									Type:        schema.TypeList,
									Description: "HostedDomains was an optional list of whitelisted domains when using the OIDC provider with Google. Only users from a whitelisted domain were allowed to log in. Support for this option was removed from the OIDC provider. Consider switching to the Google provider which supports this option.  Deprecated: will be removed in future releases.",
									Optional:    true,
									Elem:        &schema.Schema{Type: schema.TypeString},
								},
								"insecure_enable_groups": {
									Type:        schema.TypeBool,
									Description: "InsecureEnableGroups enables groups claims.",
									Optional:    true,
								},
								"insecure_skip_email_verified": {
									Type:        schema.TypeBool,
									Description: "InsecureSkipEmailVerified overrides the value of email_verified to true in the returned claims.",
									Optional:    true,
								},
								"insecure_skip_verify": {
									Type:        schema.TypeBool,
									Description: "InsecureSkipVerify disabled certificate verification. Use with caution.",
									Optional:    true,
								},
								"issuer": {
									Type:        schema.TypeString,
									Description: "Issuer is the issuer URL of the OIDC provider.",
									Required:    true,
								},
								"override_claim_mapping": {
									Type:        schema.TypeBool,
									Description: "OverrideClaimMapping will be used to override the options defined in claimMappings. i.e. if there are 'email' and `preferred_email` claims available, by default Dex will always use the `email` claim independent of the ClaimMapping.EmailKey. This setting allows you to override the default behavior of Dex and enforce the mappings defined in `claimMapping`.",
									Optional:    true,
								},
								"prompt_type": {
									Type:        schema.TypeList,
									Description: "PromptType will be used for the prompt parameter. When offline_access scope is used this defaults to prompt=consent.",
									Optional:    true,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"value": {
												Type:        schema.TypeString,
												Description: "PromptType will be used for the prompt parameter. When offline_access scope is used this defaults to prompt=consent.",
												Required:    true,
											},
										},
									},
								},
								"provider_discovery_overrides": {
									Type:        schema.TypeList,
									Description: "The section to override options discovered automatically from the providers' discovery URL (.well-known/openid-configuration).",
									Optional:    true,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"auth_url": {
												Type:        schema.TypeString,
												Description: "AuthURL provides a way to user overwrite the Auth URL from the .well-known/openid-configuration authorization_endpoint",
												Optional:    true,
											},
											"jwks_url": {
												Type:        schema.TypeString,
												Description: "JWKSURL provides a way to user overwrite the JWKS URL from the .well-known/openid-configuration jwks_uri",
												Optional:    true,
											},
											"token_url": {
												Type:        schema.TypeString,
												Description: "TokenURL provides a way to user overwrite the Token URL from the .well-known/openid-configuration token_endpoint",
												Optional:    true,
											},
										},
									},
								},
								"redirect_uri": {
									Type:        schema.TypeString,
									Description: "RedirectURI is the redirect URI of the OIDC provider.",
									Required:    true,
								},
								"root_c_as": {
									Type:        schema.TypeList,
									Description: "RootCAs are root certificates for SSL validation.",
									Optional:    true,
									Elem:        &schema.Schema{Type: schema.TypeString},
								},
								"scopes": {
									Type:        schema.TypeList,
									Description: "Scopes are the scopes requested from the OIDC provider. Defaults to \"profile\" and \"email\".",
									Optional:    true,
									Elem:        &schema.Schema{Type: schema.TypeString},
								},
								"user_id_key": {
									Type:        schema.TypeString,
									Description: "UserIDKey is the key used to identify the user in the claims.",
									Optional:    true,
								},
								"user_name_key": {
									Type:        schema.TypeString,
									Description: "UserNameKey is the key used to identify the username in the claims.",
									Optional:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}
