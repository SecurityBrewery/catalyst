package auth

import (
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/mail"
)

type Config struct {
	AppSecret string `json:"appSecret" yaml:"appSecret"`
	URL       string `json:"url" yaml:"url"`
	Email     string `json:"email" yaml:"email"`
}

type Service struct {
	config  *Config
	queries *sqlc.Queries
	mailer  *mail.Mailer
}

func New(queries *sqlc.Queries, mailer *mail.Mailer, config *Config) *Service {
	return &Service{
		config:  config,
		queries: queries,
		mailer:  mailer,
	}
}
