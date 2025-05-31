package main

import (
	"log"
	"net/http"

	"golang.org/x/net/context"

	"github.com/SecurityBrewery/catalyst/app2"
	"github.com/SecurityBrewery/catalyst/app2/fakedata"
	"github.com/SecurityBrewery/catalyst/reaction"
)

func main() {
	ctx := context.Background()

	catalyst, err := app2.App(ctx, "./catalyst_data", false)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := catalyst.Queries.CreateFeature(ctx, "dev"); err != nil {
		log.Fatal(err)
	}

	if err := fakedata.GenerateDefaultData(ctx, catalyst.Queries); err != nil {
		log.Fatal(err)
	}

	err = catalyst.SetupRoutes()
	if err != nil {
		log.Fatal(err)
	}

	reaction.BindHooks(catalyst, false)

	server := &http.Server{
		Addr:        ":8090",
		Handler:     catalyst.Router,
		ReadTimeout: 10 * 60 * 1000, // 10 minutes
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
