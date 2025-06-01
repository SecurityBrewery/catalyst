package auth

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/mail"
)

type Config struct {
	AppSecret string `json:"appSecret" yaml:"appSecret"`
	URL       string `json:"url" yaml:"url"`
	Email     string `json:"email" yaml:"email"`
	Domain    string `json:"domain" yaml:"domain"`

	OIDCAuth     bool   `json:"oidcAuth,omitempty" yaml:"oidcAuth,omitempty"`
	OIDCIssuer   string `json:"oidcIssuer,omitempty" yaml:"oidcIssuer,omitempty"`
	ClientID     string `json:"clientID,omitempty" yaml:"clientID,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty" yaml:"clientSecret,omitempty"`
	RedirectURL  string `json:"redirectURL,omitempty" yaml:"redirectURL,omitempty"`
	AuthURL      string `json:"authURL,omitempty" yaml:"authURL,omitempty"`

	UserCreateConfig *UserCreateConfig `json:"userCreateConfig,omitempty" yaml:"userCreateConfig,omitempty"`
}

type UserCreateConfig struct {
	Active bool `json:"active,omitempty" yaml:"active,omitempty"`

	OIDCClaimUsername string `json:"oidcClaimUsername,omitempty" yaml:"oidcClaimUsername,omitempty"`
	OIDCClaimEmail    string `json:"oidcClaimEmail,omitempty" yaml:"oidcClaimEmail,omitempty"`
	OIDCClaimName     string `json:"oidcClaimName,omitempty" yaml:"oidcClaimName,omitempty"`
}

type Service struct {
	config       *Config
	queries      *sqlc.Queries
	oauth2Config oauth2.Config
	mailer       *mail.Mailer
	verifier     *oidc.IDTokenVerifier
}

func New(ctx context.Context, queries *sqlc.Queries, mailer *mail.Mailer, config *Config) (*Service, error) {
	service := &Service{
		config:  config,
		queries: queries,
		mailer:  mailer,
	}

	if config.OIDCAuth {
		provider, err := oidc.NewProvider(ctx, config.OIDCIssuer)
		if err != nil {
			return nil, err
		}

		oauth2Config := oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			RedirectURL:  config.RedirectURL,

			// Discovery returns the OAuth2 endpoints.
			Endpoint: provider.Endpoint(),

			// "openid" is a required scope for OpenID Connect flows.
			Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
		}

		if config.AuthURL != "" {
			oauth2Config.Endpoint.AuthURL = config.AuthURL
		}

		service.oauth2Config = oauth2Config
		service.verifier = provider.Verifier(&oidc.Config{ClientID: config.ClientID, SkipIssuerCheck: true})
	}

	return service, nil
}
