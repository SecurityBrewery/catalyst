package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/SecurityBrewery/catalyst/app"
	"github.com/SecurityBrewery/catalyst/data"
	"github.com/SecurityBrewery/catalyst/reaction"
	"github.com/SecurityBrewery/catalyst/webhook"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	catalyst, cleanup, err := app.New(ctx, "./catalyst_data")
	if err != nil {
		return err
	}

	defer cleanup()

	if _, err := catalyst.Queries.CreateFeature(ctx, "dev"); err != nil {
		return err
	}

	testUser, err := catalyst.Queries.UserByUserName(ctx, "u_test")
	if err != nil {
		return err
	}

	if err := catalyst.Queries.DeleteUser(ctx, testUser.ID); err != nil {
		return err
	}

	slog.InfoContext(ctx, "deleting test user", "id", testUser.ID)

	if _, err := catalyst.Queries.GetUser(ctx, testUser.ID); err == nil {
		return errors.New("test user still exists after deletion")
	}

	if err := data.GenerateDemoData(ctx, catalyst.Queries, 10, 100); err != nil {
		return err
	}

	err = catalyst.SetupRoutes()
	if err != nil {
		return err
	}

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
