package main

import (
	"log"

	"golang.org/x/net/context"

	"github.com/SecurityBrewery/catalyst/app2"
)

func main() {
	catalyst, err := app2.App("./catalyst_data", false)
	if err != nil {
		log.Fatal(err)
	}

	if err := catalyst.Start(context.Background()); err != nil {
		log.Fatal(err)
	}
}
