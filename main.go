package main

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/SecurityBrewery/catalyst/app"
	"github.com/SecurityBrewery/catalyst/reaction"
	"github.com/SecurityBrewery/catalyst/webhook"
)

func main() {
	ctx := context.Background()

	catalyst, err := app.New(ctx, "./catalyst_data")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := catalyst.Queries.CreateFeature(ctx, "dev"); err != nil {
		log.Fatal(err)
	}

	/*
		testUser, err := catalyst.Queries.UserByUserName(ctx, "u_test")
		if err != nil {
			log.Fatal(err)
		}

		if err := catalyst.Queries.DeleteUser(ctx, testUser.ID); err != nil {
			return
		}


			if err := data.GenerateDefaultData(ctx, catalyst.Queries); err != nil {
			log.Fatal(err)
		}
	*/

	err = catalyst.SetupRoutes()
	if err != nil {
		log.Fatal(err)
	}

	if err := reaction.BindHooks(catalyst, false); err != nil {
		log.Fatal(err)
	}

	webhook.BindHooks(catalyst)

	server := &http.Server{
		Addr:        ":8090",
		Handler:     catalyst.Router,
		ReadTimeout: 10 * time.Minute,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
