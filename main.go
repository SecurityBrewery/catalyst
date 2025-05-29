package main

import (
	"log"

	"github.com/SecurityBrewery/catalyst/app2"
)

func main() {
	catalyst, err := app2.App("./catalyst_data", false)
	if err != nil {
		log.Fatal(err)
	}

	if err := catalyst.Start(); err != nil {
		log.Fatal(err)
	}
}
