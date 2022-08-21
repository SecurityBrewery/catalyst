package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/SecurityBrewery/catalyst"
	"github.com/SecurityBrewery/catalyst/cmd"
	"github.com/SecurityBrewery/catalyst/generated/api"
	"github.com/SecurityBrewery/catalyst/hooks"
	"github.com/SecurityBrewery/catalyst/ui"
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

	fsys, _ := fs.Sub(ui.UI, "dist")
	staticHandlerFunc := http.HandlerFunc(api.VueStatic(fsys))
	theCatalyst.Server.Get("/ui/*", http.StripPrefix("/ui", staticHandlerFunc).ServeHTTP)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), theCatalyst.Server); err != nil {
		log.Fatal(err)
	}
}
