package main

import (
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

	middlewares := []func(next http.Handler) http.Handler{
		catalyst.Authenticate(theCatalyst.DB, config.Auth),
		catalyst.AuthorizeBlockedUser(),
	}
	fsys, _ := fs.Sub(ui.UI, "dist")
	theCatalyst.Server.With(middlewares...).Get("/ui/*", api.Static(fsys))

	if err := http.ListenAndServe(":8000", theCatalyst.Server); err != nil {
		log.Fatal(err)
	}
}
