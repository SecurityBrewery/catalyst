package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/urfave/cli/v3"

	"github.com/SecurityBrewery/catalyst/app"
	"github.com/SecurityBrewery/catalyst/app/data"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/settings"
)

func main() {
	cmd := &cli.Command{
		Name:  "catalyst",
		Usage: "Catalyst CLI",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "app-url"},
			&cli.StringSliceFlag{Name: "flags"},
		},
		Commands: []*cli.Command{
			{
				Name: "migrate",
				Action: func(ctx context.Context, command *cli.Command) error {
					_, cleanup, err := setup(ctx, command)
					if err != nil {
						return fmt.Errorf("failed to setup catalyst: %w", err)
					}

					defer cleanup()

					return nil
				},
			},
			{
				Name:  "serve",
				Usage: "Start the Catalyst server",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "http", Usage: "HTTP listen address", Value: ":8090"},
				},
				Action: serve,
			},
			{
				Name:    "demo-data",
				Aliases: []string{"fake-data"},
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "users", Usage: "Number of fake users to generate", Value: 10},
					&cli.IntFlag{Name: "tickets", Usage: "Number of fake tickets to generate", Value: 100},
				},
				Action: fakeData,
			},
			{
				Name:    "upgrade-test-data",
				Aliases: []string{"default-data"},
				Usage:   "Generate default data for Catalyst",
				Action:  defaultData,
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

	ctx := context.Background()
	if err := cmd.Run(ctx, os.Args); err != nil {
		slog.ErrorContext(ctx, "Error running catalyst", "error", err)

		os.Exit(1)
	}
}

func setup(ctx context.Context, command *cli.Command) (*app.App, func(), error) {
	catalyst, cleanup, err := app.New(ctx, "./catalyst_data")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize catalyst: %w", err)
	}

	if appURL := command.String("app-url"); appURL != "" {
		_, err := settings.Update(ctx, catalyst.Queries, func(settings *settings.Settings) {
			settings.Meta.AppURL = appURL
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to update app URL: %w", err)
		}
	}

	if err := setFlags(ctx, command.StringSlice("flags"), catalyst.Queries); err != nil {
		return nil, nil, fmt.Errorf("failed to set flags: %w", err)
	}

	return catalyst, cleanup, nil
}

func serve(ctx context.Context, command *cli.Command) error {
	catalyst, cleanup, err := setup(ctx, command)
	if err != nil {
		return fmt.Errorf("failed to setup catalyst: %w", err)
	}

	defer cleanup()

	server := &http.Server{
		Addr:        ":8090",
		Handler:     catalyst,
		ReadTimeout: 10 * time.Minute,
	}

	slog.InfoContext(ctx, "Starting Catalyst server", "address", server.Addr)

	return server.ListenAndServe()
}

func fakeData(ctx context.Context, command *cli.Command) error {
	catalyst, cleanup, err := setup(ctx, command)
	if err != nil {
		return fmt.Errorf("failed to setup catalyst: %w", err)
	}

	defer cleanup()

	if err := data.GenerateDemoData(ctx, catalyst.Queries, command.Int("users"), command.Int("tickets")); err != nil {
		return err
	}

	slog.InfoContext(ctx, "Fake data generated", "users", command.Int("users"), "tickets", command.Int("tickets"))

	return nil
}

func defaultData(ctx context.Context, command *cli.Command) error {
	catalyst, cleanup, err := setup(ctx, command)
	if err != nil {
		return fmt.Errorf("failed to setup catalyst: %w", err)
	}

	defer cleanup()

	if err := data.GenerateUpgradeTestData(ctx, catalyst.Queries); err != nil {
		return fmt.Errorf("failed to generate default data: %w", err)
	}

	slog.InfoContext(ctx, "default data generated successfully")

	return nil
}

func setFlags(ctx context.Context, newFlags []string, queries *sqlc.Queries) error {
	features, err := queries.ListFeatures(ctx, sqlc.ListFeaturesParams{})
	if err != nil {
		return err
	}

	var existingFlags []string

	for _, feature := range features {
		if !slices.Contains(newFlags, feature.Key) {
			if err := queries.DeleteFeature(ctx, feature.Key); err != nil {
				return err
			}

			slog.InfoContext(ctx, "deleted feature", "name", feature.Key)

			continue
		}

		existingFlags = append(existingFlags, feature.Key)
	}

	for _, flag := range newFlags {
		if slices.Contains(existingFlags, flag) {
			continue
		}

		if _, err := queries.CreateFeature(ctx, flag); err != nil {
			return err
		}
	}

	return nil
}
