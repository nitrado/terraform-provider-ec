package provider

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/ec/armada"
	"github.com/nitrado/terraform-provider-ec/ec/container"
	"github.com/nitrado/terraform-provider-ec/ec/core"
	"gitlab.com/nitrado/b2b/ec/armada/pkg/apiclient/clientset"
	"gitlab.com/nitrado/b2b/ec/armada/pkg/apiclient/rest"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// Provider returns the enterprise console terraform provider.
func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARMADA_HOST", ""),
				Description: "The hostname (in form of URI) of Armada API.",
			},
			"token_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARMADA_TOKEN_ENDPOINT", ""),
				Description: "The URI to the token authentication endpoint.",
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARMADA_CLIENT_ID", ""),
				Description: "The oAuth2 client id to authenticate against.",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARMADA_CLIENT_SECRET", ""),
				Description: "The oAuth2 client secret to authenticate against.",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARMADA_USERNAME", ""),
				Description: "The user to authenticate with.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARMADA_PASSWORD", ""),
				Description: "The password to authenticate with.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ec_armada_site":         armada.DataSourceArmadaSite(),
			"ec_armada_site_v1":      armada.DataSourceArmadaSite(),
			"ec_armada_region":       armada.DataSourceArmadaRegion(),
			"ec_armada_region_v1":    armada.DataSourceArmadaRegion(),
			"ec_armada_armada":       armada.DataSourceArmadaArmada(),
			"ec_armada_armada_v1":    armada.DataSourceArmadaArmada(),
			"ec_armada_armadaset":    armada.DataSourceArmadaArmadaSet(),
			"ec_armada_armadaset_v1": armada.DataSourceArmadaArmadaSet(),
			"ec_container_branch":    container.DataSourceContainerBranch(),
			"ec_container_branch_v1": container.DataSourceContainerBranch(),
			"ec_core_environment":    core.DataSourceCoreEnvironment(),
			"ec_core_environment_v1": core.DataSourceCoreEnvironment(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"ec_armada_site":         armada.ResourceArmadaSite(),
			"ec_armada_site_v1":      armada.ResourceArmadaSite(),
			"ec_armada_region":       armada.ResourceArmadaRegion(),
			"ec_armada_region_v1":    armada.ResourceArmadaRegion(),
			"ec_armada_armada":       armada.ResourceArmadaArmada(),
			"ec_armada_armada_v1":    armada.ResourceArmadaArmada(),
			"ec_armada_armadaset":    armada.ResourceArmadaArmadaSet(),
			"ec_armada_armadaset_v1": armada.ResourceArmadaArmadaSet(),
			"ec_container_branch":    container.ResourceContainerBranch(),
			"ec_container_branch_v1": container.ResourceContainerBranch(),
			"ec_core_environment":    core.ResourceCoreEnvironment(),
			"ec_core_environment_v1": core.ResourceCoreEnvironment(),
		},
	}

	p.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		return providerConfigure(ctx, d, p.TerraformVersion)
	}

	return p
}

func providerConfigure(_ context.Context, d *schema.ResourceData, _ string) (any, diag.Diagnostics) {
	//nolint:contextcheck
	tok, err := resolveToken(d)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	if _, err = tok.Token(); err != nil {
		return nil, diag.FromErr(err)
	}

	cfg := createRESTConfig(d)
	cfg.BearerTokenSource = tok

	clientSet, err := clientset.New(cfg)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("could not to configure client: %w", err))
	}

	return clientSet, nil
}

func createRESTConfig(d *schema.ResourceData) rest.Config {
	var cfg rest.Config
	cfg.BaseURL = d.Get("host").(string)

	return cfg
}

func resolveToken(d *schema.ResourceData) (oauth2.TokenSource, error) {
	tokURL, err := resolveTokenURL(d)
	if err != nil {
		return nil, err
	}

	clientID := d.Get("client_id").(string)
	clientSecret, hasClientSecret := d.Get("client_secret").(string)
	user, hasUser := d.Get("username").(string)
	pass, hasPass := d.Get("password").(string)

	switch {
	case hasUser && hasPass:
		cfg := oauth2.Config{
			ClientID: clientID,
			Endpoint: oauth2.Endpoint{
				TokenURL:  tokURL,
				AuthStyle: oauth2.AuthStyleInHeader,
			},
		}
		tok, err := cfg.PasswordCredentialsToken(context.Background(), user, pass)
		if err != nil {
			return nil, err
		}
		return cfg.TokenSource(context.Background(), tok), nil
	case hasClientSecret:
		cfg := clientcredentials.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			TokenURL:     tokURL,
			AuthStyle:    oauth2.AuthStyleInHeader,
		}
		return cfg.TokenSource(context.Background()), nil
	default:
		return nil, fmt.Errorf("either client_secret or username and password must be set")
	}
}

func resolveTokenURL(d *schema.ResourceData) (string, error) {
	if tokenEndpoint, ok := d.Get("token_endpoint").(string); ok && tokenEndpoint != "" {
		return tokenEndpoint, nil
	}

	host := d.Get("host").(string)
	u, err := url.Parse(host)
	if err != nil {
		return "", fmt.Errorf("invalid host: %w", err)
	}
	u.Host = "auth-" + u.Host
	u.Path = "/auth/realms/enterprise-console/protocol/openid-connect/token"
	return u.String(), nil
}
