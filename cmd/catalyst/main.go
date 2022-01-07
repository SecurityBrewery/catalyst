package main

import (
	"log"
	"net/http"

	"github.com/SecurityBrewery/catalyst"
	"github.com/SecurityBrewery/catalyst/cmd"
	"github.com/SecurityBrewery/catalyst/hooks"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config, err := cmd.ParseCatalystConfig()
	if err != nil {
		log.Fatal(err)
	}

	theCatalyst, err := catalyst.New(&hooks.Hooks{}, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":8000", theCatalyst.Server); err != nil {
		log.Fatal(err)
	}
}
