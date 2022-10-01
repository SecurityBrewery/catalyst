package auth

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/exp/slices"
	"golang.org/x/oauth2"
)

type Config struct {
	CookieSecret []byte

	SimpleAuthEnable bool
	APIKeyAuthEnable bool
	OIDCAuthEnable   bool

	InitialUser     string
	InitialPassword string
	InitialAPIKey   string

	OIDCIssuer       string
	AuthURL          string
	OAuth2           *oauth2.Config
	UserCreateConfig *UserCreateConfig

	provider *oidc.Provider
}

func (c *Config) Clone() *Config {
	return &Config{
		CookieSecret:     slices.Clone(c.CookieSecret),
		SimpleAuthEnable: c.SimpleAuthEnable,
		APIKeyAuthEnable: c.APIKeyAuthEnable,
		OIDCAuthEnable:   c.OIDCAuthEnable,
		OIDCIssuer:       c.OIDCIssuer,
		AuthURL:          c.AuthURL,
		InitialUser:      c.InitialUser,
		InitialPassword:  c.InitialPassword,
		InitialAPIKey:    c.InitialAPIKey,
		OAuth2: &oauth2.Config{
			ClientID:     c.OAuth2.ClientID,
			ClientSecret: c.OAuth2.ClientSecret,
			RedirectURL:  c.OAuth2.RedirectURL,
			Scopes:       slices.Clone(c.OAuth2.Scopes),
			Endpoint: oauth2.Endpoint{
				AuthURL:   c.OAuth2.Endpoint.AuthURL,
				TokenURL:  c.OAuth2.Endpoint.TokenURL,
				AuthStyle: c.OAuth2.Endpoint.AuthStyle,
			},
		},
		UserCreateConfig: &UserCreateConfig{
			AuthBlockNew:      c.UserCreateConfig.AuthBlockNew,
			AuthDefaultRoles:  slices.Clone(c.UserCreateConfig.AuthDefaultRoles),
			AuthAdminUsers:    slices.Clone(c.UserCreateConfig.AuthAdminUsers),
			OIDCClaimUsername: c.UserCreateConfig.OIDCClaimUsername,
			OIDCClaimEmail:    c.UserCreateConfig.OIDCClaimEmail,
			OIDCClaimName:     c.UserCreateConfig.OIDCClaimName,
			OIDCClaimGroups:   c.UserCreateConfig.OIDCClaimGroups,
		},
	}
}

type UserCreateConfig struct {
	AuthBlockNew     bool
	AuthDefaultRoles []string
	AuthAdminUsers   []string

	OIDCClaimUsername string
	OIDCClaimEmail    string
	OIDCClaimName     string
	OIDCClaimGroups   string
}

type Authenticator struct {
	config   *Config
	resolver LoginResolver
	jar      *Jar
}

func NewAuthenticator(ctx context.Context, config *Config, resolver LoginResolver) (*Authenticator, error) {
	a := &Authenticator{
		config:   config,
		resolver: resolver,
		jar:      NewJar(config.CookieSecret),
	}

	if config.InitialAPIKey != "" {
		if err := resolver.UserCreateIfNotExists(ctx, nil, config.InitialAPIKey); err != nil {
			return nil, err
		}
	}
	if config.InitialUser != "" && config.InitialPassword != "" {
		if err := resolver.UserCreateIfNotExists(ctx, &User{
			ID:      config.InitialUser,
			APIKey:  false,
			Blocked: false,
			Roles:   []string{AdminRole},
		}, config.InitialPassword); err != nil {
			return nil, err
		}
	}

	if config.OIDCAuthEnable {
		err := a.Load(ctx)
		if err != nil {
			return nil, err
		}
	}

	return a, nil
}

func (a *Authenticator) Verifier(ctx context.Context) (*oidc.IDTokenVerifier, error) {
	if a.config.provider == nil {
		if err := a.Load(ctx); err != nil {
			return nil, err
		}
	}

	config := &oidc.Config{ClientID: a.config.OAuth2.ClientID}
	if a.config.AuthURL != "" {
		config.SkipIssuerCheck = true
	}

	return a.config.provider.Verifier(config), nil
}

func (a *Authenticator) Load(ctx context.Context) error {
	provider, err := oidc.NewProvider(ctx, a.config.OIDCIssuer)
	if err == nil {
		a.config.provider = provider
		a.config.OAuth2.Endpoint = provider.Endpoint()
		if a.config.AuthURL != "" {
			a.config.OAuth2.Endpoint.AuthURL = a.config.AuthURL
		}
	}

	return err
}

func (a *Authenticator) Config() *Config {
	return a.config
}
