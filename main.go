package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/urfave/cli/v3"

	"github.com/SecurityBrewery/catalyst/app"
	"github.com/SecurityBrewery/catalyst/app/auth"
	"github.com/SecurityBrewery/catalyst/app/data"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/mail"
	"github.com/SecurityBrewery/catalyst/app/reaction"
	"github.com/SecurityBrewery/catalyst/app/webhook"
	"github.com/SecurityBrewery/catalyst/upgradetest"
)

func main() {
	cmd := &cli.Command{
		Name:  "catalyst",
		Usage: "Catalyst CLI",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "app-url"},
			&cli.StringFlag{Name: "email"},
			&cli.StringSliceFlag{Name: "flags"},
			&cli.StringFlag{Name: "smtp-server"},
			&cli.StringFlag{Name: "smtp-user"},
			&cli.StringFlag{Name: "smtp-password"},
		},
		Commands: []*cli.Command{
			{
				Name:   "serve",
				Usage:  "Start the Catalyst server",
				Action: serve,
			},
			{
				Name: "fake-data",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "users", Usage: "Number of fake users to generate", Value: 10},
					&cli.IntFlag{Name: "tickets", Usage: "Number of fake tickets to generate", Value: 100},
				},
				Action: fakeData,
			},
			{
				Name:  "default-data",
				Usage: "Generate default data for Catalyst",
				Action: func(ctx context.Context, command *cli.Command) error {
					catalyst, cleanup, err := setup(ctx, command)
					if err != nil {
						return fmt.Errorf("failed to setup catalyst: %w", err)
					}

					defer cleanup()

					if err := upgradetest.GenerateUpgradeTestData(ctx, catalyst.Queries); err != nil {
						return fmt.Errorf("failed to generate default data: %w", err)
					}

					slog.InfoContext(ctx, "default data generated successfully")

					return nil
				},
			},
			{
				Name: "admin",
				Commands: []*cli.Command{
					{Name: "create", Action: adminCreate},
					{Name: "delete", Action: adminDelete},
					{Name: "set-password", Action: adminSetPassword},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func setup(ctx context.Context, command *cli.Command) (*app.App, func(), error) {
	catalyst, cleanup, err := app.New(ctx, "./catalyst_data", &app.Config{
		Auth: &auth.Config{
			AppSecret: "", // TODO: set a secure secret
			URL:       command.String("app-url"),
			Email:     command.String("email"),
		},
		Mail: &mail.Config{
			SMTPServer:   command.String("smtp-server"),
			SMTPUser:     command.String("smtp-user"),
			SMTPPassword: command.String("smtp-password"),
		},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize catalyst: %w", err)
	}

	if err := setFlags(ctx, command.StringSlice("flags"), catalyst); err != nil {
		return nil, nil, fmt.Errorf("failed to set flags: %w", err)
	}

	return catalyst, cleanup, nil
}

func serve(ctx context.Context, command *cli.Command) error {
	catalyst, cleanup, err := setup(ctx, command)
	if err != nil {
		return fmt.Errorf("failed to setup catalyst: %w", err)
	}

	if err := catalyst.SetupRoutes(); err != nil {
		return err
	}

	defer cleanup()

	if err := reaction.BindHooks(catalyst, false); err != nil {
		return err
	}

	webhook.BindHooks(catalyst)

	server := &http.Server{
		Addr:        ":8090",
		Handler:     catalyst.Router,
		ReadTimeout: 10 * time.Minute,
	}

	return server.ListenAndServe()
}

func fakeData(ctx context.Context, command *cli.Command) error {
	catalyst, cleanup, err := setup(ctx, command)
	if err != nil {
		return fmt.Errorf("failed to setup catalyst: %w", err)
	}

	if err := catalyst.SetupRoutes(); err != nil {
		return err
	}

	defer cleanup()

	if err := data.GenerateDemoData(ctx, catalyst.Queries, command.Int("users"), command.Int("tickets")); err != nil {
		return err
	}

	slog.InfoContext(ctx, "fake data generated", "users", command.Int("users"), "tickets", command.Int("tickets"))

	return nil
}

func setFlags(ctx context.Context, newFlags []string, catalyst *app.App) error {
	features, err := catalyst.Queries.ListFeatures(ctx, sqlc.ListFeaturesParams{})
	if err != nil {
		return err
	}

	var existingFlags []string

	for _, feature := range features {
		if !slices.Contains(newFlags, feature.Name) {
			if err := catalyst.Queries.DeleteFeature(ctx, feature.Name); err != nil {
				return err
			}

			slog.InfoContext(ctx, "deleted feature", "name", feature.Name)

			continue
		}

		existingFlags = append(existingFlags, feature.Name)
	}

	for _, flag := range newFlags {
		if slices.Contains(existingFlags, flag) {
			continue
		}

		if _, err := catalyst.Queries.CreateFeature(ctx, flag); err != nil {
			return err
		}
	}

	return nil
}
