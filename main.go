package main

import (
	"log"

	"golang.org/x/net/context"

	"github.com/SecurityBrewery/catalyst/app2"
	"github.com/SecurityBrewery/catalyst/app2/fakedata"
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

	if err := catalyst.Start(ctx); err != nil {
		log.Fatal(err)
	}
}
