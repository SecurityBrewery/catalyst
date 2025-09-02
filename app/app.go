package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/hook"
	"github.com/SecurityBrewery/catalyst/app/mail"
	"github.com/SecurityBrewery/catalyst/app/migration"
	"github.com/SecurityBrewery/catalyst/app/reaction"
	"github.com/SecurityBrewery/catalyst/app/reaction/schedule"
	"github.com/SecurityBrewery/catalyst/app/router"
	"github.com/SecurityBrewery/catalyst/app/service"
	"github.com/SecurityBrewery/catalyst/app/upload"
	"github.com/SecurityBrewery/catalyst/app/webhook"
)

type App struct {
	Queries *sqlc.Queries
	Hooks   *hook.Hooks
	router  http.Handler
}

func New(ctx context.Context, dir string) (*App, func(), error) {
	uploader, err := upload.New(dir)
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

	scheduler, err := schedule.New(ctx, queries)
	if err != nil {
		return nil, cleanup, fmt.Errorf("failed to create scheduler: %w", err)
	}

	hooks := hook.NewHooks()

	service := service.New(queries, hooks, uploader, scheduler)

	router, err := router.New(service, queries, uploader, mailer)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create router: %w", err)
	}

	if err := reaction.BindHooks(hooks, router, queries, false); err != nil {
		return nil, nil, err
	}

	webhook.BindHooks(hooks, queries)

	app := &App{
		Queries: queries,
		Hooks:   hooks,
		router:  router,
	}

	return app, cleanup, nil
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}
