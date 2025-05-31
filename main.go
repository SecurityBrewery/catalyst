package main

import (
	"log"

	"golang.org/x/net/context"

	"github.com/SecurityBrewery/catalyst/app2"
	"github.com/SecurityBrewery/catalyst/app2/fakedata"
)

func main() {
	catalyst, err := app2.App("./catalyst_data", false)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	if err := fakedata.GenerateDefaultData(ctx, catalyst.Queries); err != nil {
		log.Fatal(err)
	}

	if err := catalyst.Start(ctx); err != nil {
		log.Fatal(err)
	}
}
