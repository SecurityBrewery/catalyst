package app

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/go-chi/chi/v5"

	"github.com/SecurityBrewery/catalyst/app/auth"
	"github.com/SecurityBrewery/catalyst/app/auth/oidc"
	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/hook"
	"github.com/SecurityBrewery/catalyst/app/mail"
	"github.com/SecurityBrewery/catalyst/app/reaction/schedule"
	"github.com/SecurityBrewery/catalyst/app/service"
)

type App struct {
	Queries   *sqlc.Queries
	Router    *chi.Mux
	Service   *service.Service
	Config    *auth.Config
	Auth      *auth.Service
	Hooks     *hook.Hooks
	Scheduler *schedule.Scheduler
}

func New(ctx context.Context, filename string) (*App, func(), error) {
	queries, cleanup, err := database.DB(ctx, filepath.Join(filename, "data.db"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	mailer := mail.New(&mail.Config{
		SMTPServer:   "localhost",
		SMTPUser:     "",
		SMTPPassword: "",
	})

	config := &auth.Config{
		URL:    "http://localhost:8080",
		Email:  "info@cataly-soar.com",
		Domain: "localhost",

		UserCreateConfig: &auth.UserCreateConfig{
			Active: false,
		},
		OIDC: &oidc.Config{
			OIDCAuth:          false,
			OIDCIssuer:        "",
			ClientID:          "",
			ClientSecret:      "",
			RedirectURL:       "",
			AuthURL:           "",
			OIDCClaimUsername: "",
			OIDCClaimEmail:    "",
			OIDCClaimName:     "",
		},
	}

	authService, err := auth.New(ctx, queries, mailer, config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create auth service: %w", err)
	}

	scheduler, err := schedule.New(ctx, config, authService, queries)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	hooks := hook.NewHooks()

	return &App{
		Hooks:     hooks,
		Queries:   queries,
		Router:    chi.NewRouter(),
		Service:   service.New(queries, hooks, scheduler),
		Config:    config,
		Auth:      authService,
		Scheduler: scheduler,
	}, cleanup, nil
}
