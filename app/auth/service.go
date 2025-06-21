package auth

import (
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/mail"
)

type Config struct {
	AuthToken  string
	ResetToken string
	URL        string
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
