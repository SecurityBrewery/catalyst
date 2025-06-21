package auth

import (
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/mail"
)

type Service struct {
	queries *sqlc.Queries
	mailer  *mail.Mailer
}

func New(queries *sqlc.Queries, mailer *mail.Mailer) *Service {
	return &Service{
		queries: queries,
		mailer:  mailer,
	}
}
