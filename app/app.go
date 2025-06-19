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

func New(ctx context.Context, filename string) (*App, func(), error) {
	queries, cleanup, err := database.DB(ctx, filepath.Join(filename, "data.db"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	settings, err := database.LoadSettings(ctx, queries)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get settings: %w", err)
	}

	mailer := mail.New(&mail.Config{
		SMTPServer:   settings.SMTP.Host,
		SMTPUser:     settings.SMTP.Username,
		SMTPPassword: settings.SMTP.Password,
	})

	authConfig := &auth.Config{
		AppSecret: settings.RecordAuthToken.Secret, // TODO: support more secrets
		URL:       settings.Meta.AppURL,
		Email:     settings.Meta.SenderAddress,
	}

	authService := auth.New(queries, mailer, authConfig)

	scheduler, err := schedule.New(ctx, authConfig, authService, queries)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	hooks := hook.NewHooks()

	return &App{
		Hooks:     hooks,
		Queries:   queries,
		Router:    chi.NewRouter(),
		Service:   service.New(queries, hooks, scheduler),
		Config:    authConfig,
		Auth:      authService,
		Scheduler: scheduler,
	}, cleanup, nil
}
