package main

import (
	"log"

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

	if err = theCatalyst.Server.RunWithSigHandler(); err != nil {
		log.Fatal(err)
	}
}
