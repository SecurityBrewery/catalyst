package app

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/go-chi/chi/v5"

	"github.com/SecurityBrewery/catalyst/app/auth"
	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/hook"
	"github.com/SecurityBrewery/catalyst/app/service"
	"github.com/SecurityBrewery/catalyst/reaction/schedule"
)

type App struct {
	Queries   *sqlc.Queries
	Router    *chi.Mux
	Service   *service.Service
	Auth      *auth.Service
	Hooks     *hook.Hooks
	Scheduler *schedule.Scheduler
}

func New(ctx context.Context, filename string) (*App, error) {
	queries, _, err := database.DB(filepath.Join(filename, "data.db"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	authService, err := auth.New(ctx, queries, &auth.Config{
		Domain:       "localhost",
		CookieSecure: false,
		PasswordAuth: true,
		BearerAuth:   true,
		OIDCAuth:     false,
		OIDCIssuer:   "",
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "",
		AuthURL:      "",
		UserCreateConfig: &auth.UserCreateConfig{
			Active:            false,
			OIDCClaimUsername: "",
			OIDCClaimEmail:    "",
			OIDCClaimName:     "",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create auth service: %w", err)
	}

	scheduler, err := schedule.New(ctx, queries)
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	hooks := hook.NewHooks()

	return &App{
		Hooks:     hooks,
		Queries:   queries,
		Router:    chi.NewRouter(),
		Service:   service.New(queries, hooks, scheduler),
		Auth:      authService,
		Scheduler: scheduler,
	}, nil
}
