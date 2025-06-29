package app

import (
	"context"
	"fmt"

	"github.com/go-chi/chi/v5"

	"github.com/SecurityBrewery/catalyst/app/auth"
	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/hook"
	"github.com/SecurityBrewery/catalyst/app/mail"
	"github.com/SecurityBrewery/catalyst/app/migration"
	"github.com/SecurityBrewery/catalyst/app/reaction/schedule"
	"github.com/SecurityBrewery/catalyst/app/service"
	"github.com/SecurityBrewery/catalyst/app/upload/uploader"
)

type App struct {
	Queries   *sqlc.Queries
	Router    *chi.Mux
	Service   *service.Service
	Auth      *auth.Service
	Hooks     *hook.Hooks
	Scheduler *schedule.Scheduler
	Uploader  *uploader.Uploader
}

func New(ctx context.Context, dir string) (*App, func(), error) {
	uploader, err := uploader.New(dir)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create uploader: %w", err)
	}

	queries, cleanup, err := database.DB(ctx, dir)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := migration.Apply(ctx, queries, dir, uploader); err != nil {
		return nil, nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	mailer := mail.New(queries)

	authService := auth.New(queries, mailer)

	scheduler, err := schedule.New(ctx, authService, queries)
	if err != nil {
		return nil, cleanup, fmt.Errorf("failed to create scheduler: %w", err)
	}

	hooks := hook.NewHooks()

	app := &App{
		Hooks:     hooks,
		Queries:   queries,
		Router:    chi.NewRouter(),
		Service:   service.New(queries, hooks, uploader, scheduler),
		Auth:      authService,
		Uploader:  uploader,
		Scheduler: scheduler,
	}

	if err := app.setupRoutes(); err != nil {
		return nil, cleanup, fmt.Errorf("failed to setup routes: %w", err)
	}

	return app, cleanup, nil
}
