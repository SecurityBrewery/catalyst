package main

import (
	"log"

	"github.com/SecurityBrewery/catalyst/app"
)

func main() {
	if err := app.App("./catalyst_data").Start(); err != nil {
		log.Fatal(err)
	}
}
