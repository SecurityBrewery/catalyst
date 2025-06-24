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
	"github.com/SecurityBrewery/catalyst/app/upload"
)

type App struct {
	Queries   *sqlc.Queries
	Router    *chi.Mux
	Service   *service.Service
	Auth      *auth.Service
	Hooks     *hook.Hooks
	Scheduler *schedule.Scheduler
	Uploader  *upload.Uploader
}

func New(ctx context.Context, dir string) (*App, func(), error) {
	queries, cleanup, err := database.DB(ctx, filepath.Join(dir, "data.db"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	mailer := mail.New(queries)

	authService := auth.New(queries, mailer)

	scheduler, err := schedule.New(ctx, authService, queries)
	if err != nil {
		return nil, cleanup, fmt.Errorf("failed to create scheduler: %w", err)
	}

	hooks := hook.NewHooks()

	uploader := upload.NewUploader(dir, authService, queries)

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
