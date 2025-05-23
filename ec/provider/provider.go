package provider

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"

	"github.com/gamefabric/gf-core/pkg/apiclient/clientset"
	"github.com/gamefabric/gf-core/pkg/apiclient/rest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec"
	"github.com/nitrado/terraform-provider-ec/ec/armada"
	"github.com/nitrado/terraform-provider-ec/ec/authentication"
	"github.com/nitrado/terraform-provider-ec/ec/container"
	"github.com/nitrado/terraform-provider-ec/ec/core"
	"github.com/nitrado/terraform-provider-ec/ec/formation"
	"github.com/nitrado/terraform-provider-ec/ec/protection"
	"github.com/nitrado/terraform-provider-ec/ec/storage"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// Provider returns the enterprise console terraform provider.
func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("EC_HOST", ""),
				Description: "The hostname (in form of URI) of the Enterprise Console API.",
			},
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("EC_TOKEN", ""),
				Description: "The oAuth token to authenticate with.",
			},
			"token_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("EC_TOKEN_ENDPOINT", ""),
				Description: "The URI to the token authentication endpoint.",
			},
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("EC_CLIENT_ID", ""),
				Description: "The oAuth2 client id to authenticate against.",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("EC_CLIENT_SECRET", ""),
				Description: "The oAuth2 client secret to authenticate against.",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("EC_USERNAME", ""),
				Description: "The user to authenticate with.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("EC_PASSWORD", ""),
				Description: "The password to authenticate with.",
			},
			"instances": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Named Enterprise Console instances.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The instance name.",
						},
						"host": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The hostname (in form of URI) of the Enterprise Console API.",
						},
						"token": {
							Type:        schema.TypeString,
							Optional:    true,
							DefaultFunc: schema.EnvDefaultFunc("EC_TOKEN", ""),
							Description: "The oAuth token to authenticate with.",
						},
						"token_endpoint": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The URI to the token authentication endpoint.",
						},
						"client_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The oAuth2 client id to authenticate against.",
						},
						"client_secret": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The oAuth2 client secret to authenticate against.",
						},
						"username": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The user to authenticate with.",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The password to authenticate with.",
						},
					},
				},
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ec_armada_armada":                              armada.DataSourceArmada(),
			"ec_armada_armada_v1":                           armada.DataSourceArmada(),
			"ec_armada_armadaset":                           armada.DataSourceArmadaSet(),
			"ec_armada_armadaset_v1":                        armada.DataSourceArmadaSet(),
			"ec_container_branch":                           container.DataSourceBranch(),
			"ec_container_branch_v1":                        container.DataSourceBranch(),
			"ec_container_image":                            container.DataSourceImage(),
			"ec_container_image_v1":                         container.DataSourceImage(),
			"ec_core_environment":                           core.DataSourceEnvironment(),
			"ec_core_environment_v1":                        core.DataSourceEnvironment(),
			"ec_core_location":                              core.DataSourceLocation(),
			"ec_core_location_v1":                           core.DataSourceLocation(),
			"ec_core_region":                                core.DataSourceRegion(),
			"ec_core_region_v1":                             core.DataSourceRegion(),
			"ec_core_site":                                  core.DataSourceSite(),
			"ec_core_site_v1":                               core.DataSourceSite(),
			"ec_formation_vessel":                           formation.DataSourceVessel(),
			"ec_formation_vessel_v1beta1":                   formation.DataSourceVessel(),
			"ec_formation_formation":                        formation.DataSourceFormation(),
			"ec_formation_formation_v1beta1":                formation.DataSourceFormation(),
			"ec_protection_gatewaypolicy":                   protection.DataSourceGatewayPolicy(),
			"ec_protection_gatewaypolicy_v1alpha1":          protection.DataSourceGatewayPolicy(),
			"ec_protection_migration":                       protection.DataSourceMigration(),
			"ec_protection_migration_v1alpha1":              protection.DataSourceMigration(),
			"ec_protection_protocol":                        protection.DataSourceProtocol(),
			"ec_protection_protocol_v1alpha1":               protection.DataSourceProtocol(),
			"ec_storage_volume":                             storage.DataSourceVolume(),
			"ec_storage_volume_v1beta1":                     storage.DataSourceVolume(),
			"ec_storage_volumestore":                        storage.DataSourceVolumeStore(),
			"ec_storage_volumestore_v1beta1":                storage.DataSourceVolumeStore(),
			"ec_storage_volumestoreretentionpolicy":         storage.DataSourceVolumeStoreRetentionPolicy(),
			"ec_storage_volumestoreretentionpolicy_v1beta1": storage.DataSourceVolumeStoreRetentionPolicy(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"ec_armada_armada":                              armada.ResourceArmada(),
			"ec_armada_armada_v1":                           armada.ResourceArmada(),
			"ec_armada_armadaset":                           armada.ResourceArmadaSet(),
			"ec_armada_armadaset_v1":                        armada.ResourceArmadaSet(),
			"ec_authentication_provider":                    authentication.ResourceProvider(),
			"ec_authentication_provider_v1beta1":            authentication.ResourceProvider(),
			"ec_authentication_serviceaccount":              authentication.ResourceServiceAccount(),
			"ec_authentication_serviceaccount_v1beta1":      authentication.ResourceServiceAccount(),
			"ec_container_branch":                           container.ResourceBranch(),
			"ec_container_branch_v1":                        container.ResourceBranch(),
			"ec_core_environment":                           core.ResourceEnvironment(),
			"ec_core_environment_v1":                        core.ResourceEnvironment(),
			"ec_core_site":                                  core.ResourceSite(),
			"ec_core_site_v1":                               core.ResourceSite(),
			"ec_core_region":                                core.ResourceRegion(),
			"ec_core_region_v1":                             core.ResourceRegion(),
			"ec_formation_vessel":                           formation.ResourceVessel(),
			"ec_formation_vessel_v1beta1":                   formation.ResourceVessel(),
			"ec_formation_formation":                        formation.ResourceFormation(),
			"ec_formation_formation_v1beta1":                formation.ResourceFormation(),
			"ec_protection_gatewaypolicy":                   protection.ResourceGatewayPolicy(),
			"ec_protection_gatewaypolicy_v1alpha1":          protection.ResourceGatewayPolicy(),
			"ec_protection_protocol":                        protection.ResourceProtocol(),
			"ec_protection_protocol_v1alpha1":               protection.ResourceProtocol(),
			"ec_storage_volume":                             storage.ResourceVolume(),
			"ec_storage_volume_v1beta1":                     storage.ResourceVolume(),
			"ec_storage_volumestore":                        storage.ResourceVolumeStore(),
			"ec_storage_volumestore_v1beta1":                storage.ResourceVolumeStore(),
			"ec_storage_volumestoreretentionpolicy":         storage.ResourceVolumeStoreRetentionPolicy(),
			"ec_storage_volumestoreretentionpolicy_v1beta1": storage.ResourceVolumeStoreRetentionPolicy(),
		},
	}

	p.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		return providerConfigure(ctx, d, p.TerraformVersion)
	}

	return p
}

func providerConfigure(_ context.Context, d *schema.ResourceData, _ string) (any, diag.Diagnostics) {
	var defaultClientSet clientset.Interface
	if v, ok := d.Get("host").(string); ok && v != "" {
		var err error
		defaultClientSet, err = createClientSet("", collectConnData(d)) //nolint:contextcheck
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}

	var instances map[string]clientset.Interface
	if insts, ok := d.Get("instances").([]any); ok {
		instances = make(map[string]clientset.Interface, len(insts))
		for _, v := range insts {
			inst := v.(map[string]any)

			name := inst["name"].(string)
			cs, err := createClientSet(name, inst) //nolint:contextcheck
			if err != nil {
				return nil, diag.FromErr(err)
			}
			instances[name] = cs
		}
	}

	if defaultClientSet == nil && len(instances) == 0 {
		return nil, diag.FromErr(errors.New("at least one instance or default connection details must be provided"))
	}
	return ec.NewProviderContext(defaultClientSet, instances), nil
}

func collectConnData(d *schema.ResourceData) map[string]any {
	return map[string]any{
		"host":           d.Get("host"),
		"token_endpoint": d.Get("token_endpoint"),
		"client_id":      d.Get("client_id"),
		"client_secret":  d.Get("client_secret"),
		"username":       d.Get("username"),
		"password":       d.Get("password"),
		"token":          d.Get("token"),
	}
}

func createClientSet(name string, m map[string]any) (clientset.Interface, error) {
	var forInstance string
	if name != "" {
		forInstance = `for instance "` + name + `"`
	}

	tok, err := resolveToken(m)
	if err != nil {
		return nil, fmt.Errorf("retrieveing token %s: %w", forInstance, err)
	}

	cfg := createRESTConfig(m, tok)
	clientSet, err := clientset.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("could not to configure client %s: %w", forInstance, err)
	}

	return clientSet, nil
}

func createRESTConfig(m map[string]any, tok oauth2.TokenSource) rest.Config {
	var cfg rest.Config
	cfg.BaseURL = m["host"].(string)
	cfg.BearerTokenSource = tok

	return cfg
}

func resolveToken(m map[string]any) (oauth2.TokenSource, error) {
	if token, ok := m["token"].(string); ok && token != "" {
		return oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: token,
			TokenType:   "Bearer",
		}), nil
	}

	tokURL, err := resolveTokenURL(m)
	if err != nil {
		return nil, err
	}

	clientID, ok := m["client_id"].(string)
	if !ok {
		return nil, errors.New("client id is required")
	}

	clientSecret, hasClientSecret := m["client_secret"].(string)
	user, hasUser := m["username"].(string)
	pass, hasPass := m["password"].(string)

	switch {
	case hasUser && hasPass:
		return newLazyTokenSource(func() (oauth2.TokenSource, error) {
			cfg := oauth2.Config{
				ClientID: clientID,
				Scopes:   []string{"openid", "email", "profile", "offline_access"},
				Endpoint: oauth2.Endpoint{
					AuthStyle: oauth2.AuthStyleInHeader,
					TokenURL:  tokURL,
				},
			}
			tok, err := cfg.PasswordCredentialsToken(context.Background(), user, pass)
			if err != nil {
				return nil, err
			}
			return cfg.TokenSource(context.Background(), tok), nil
		}), nil
	case hasClientSecret:
		return newLazyTokenSource(func() (oauth2.TokenSource, error) {
			cfg := clientcredentials.Config{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				TokenURL:     tokURL,
				Scopes:       []string{"openid", "email", "profile", "offline_access"},
				AuthStyle:    oauth2.AuthStyleInHeader,
			}
			return cfg.TokenSource(context.Background()), nil
		}), nil
	default:
		return nil, errors.New("either client_secret or username and password must be set")
	}
}

func resolveTokenURL(m map[string]any) (string, error) {
	if tokenEndpoint, ok := m["token_endpoint"].(string); ok && tokenEndpoint != "" {
		return tokenEndpoint, nil
	}

	host := m["host"].(string)
	u, err := url.Parse(host)
	if err != nil {
		return "", fmt.Errorf("invalid host: %w", err)
	}
	u.Path = "/auth/token"
	return u.String(), nil
}

type lazyTokenSource struct {
	mu    sync.Mutex
	ts    oauth2.TokenSource
	newFn func() (oauth2.TokenSource, error)
}

func newLazyTokenSource(newFn func() (oauth2.TokenSource, error)) *lazyTokenSource {
	return &lazyTokenSource{newFn: newFn}
}

func (s *lazyTokenSource) Token() (*oauth2.Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ts == nil {
		var err error
		s.ts, err = s.newFn()
		if err != nil {
			return nil, err
		}
	}
	return s.ts.Token()
}
