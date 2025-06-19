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

type Config struct {
	Auth *auth.Config `json:"auth" yaml:"auth"`
	Mail *mail.Config `json:"mail" yaml:"mail"`
}

func New(ctx context.Context, filename string, config *Config) (*App, func(), error) {
	queries, cleanup, err := database.DB(ctx, filepath.Join(filename, "data.db"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	mailer := mail.New(config.Mail)

	authService := auth.New(queries, mailer, config.Auth)

	scheduler, err := schedule.New(ctx, config.Auth, authService, queries)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	hooks := hook.NewHooks()

	return &App{
		Hooks:     hooks,
		Queries:   queries,
		Router:    chi.NewRouter(),
		Service:   service.New(queries, hooks, scheduler),
		Config:    config.Auth,
		Auth:      authService,
		Scheduler: scheduler,
	}, cleanup, nil
}
