package auth

import (
	"context"

	"github.com/SecurityBrewery/catalyst/app/auth/oidc"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/mail"
)

type Config struct {
	AppSecret string `json:"appSecret" yaml:"appSecret"`
	URL       string `json:"url" yaml:"url"`
	Email     string `json:"email" yaml:"email"`
	Domain    string `json:"domain" yaml:"domain"`

	UserCreateConfig *UserCreateConfig `json:"userCreateConfig,omitempty" yaml:"userCreateConfig,omitempty"`

	OIDC *oidc.Config `json:"oidc,omitempty" yaml:"oidc,omitempty"`
}

type UserCreateConfig struct {
	Active bool `json:"active,omitempty" yaml:"active,omitempty"`
}

type Service struct {
	config  *Config
	queries *sqlc.Queries
	mailer  *mail.Mailer
	oidc    *oidc.Server
}

func New(ctx context.Context, queries *sqlc.Queries, mailer *mail.Mailer, config *Config) (*Service, error) {
	service := &Service{
		config:  config,
		queries: queries,
		mailer:  mailer,
	}

	if config.OIDC != nil && config.OIDC.OIDCAuth {
		oidc, err := oidc.New(ctx, queries, config.OIDC)
		if err != nil {
			return nil, err
		}

		service.oidc = oidc
	}

	return service, nil
}
