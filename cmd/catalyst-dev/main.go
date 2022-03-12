package main

import (
	"context"
	"log"
	"net/http"

	"github.com/arangodb/go-driver"

	"github.com/SecurityBrewery/catalyst"
	"github.com/SecurityBrewery/catalyst/cmd"
	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/api"
	"github.com/SecurityBrewery/catalyst/generated/model"
	"github.com/SecurityBrewery/catalyst/hooks"
	"github.com/SecurityBrewery/catalyst/role"
	"github.com/SecurityBrewery/catalyst/test"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config, err := cmd.ParseCatalystConfig()
	if err != nil {
		log.Fatal(err)
	}

	// create app and clear db after start
	theCatalyst, err := catalyst.New(&hooks.Hooks{
		DatabaseAfterConnectFuncs: []func(ctx context.Context, client driver.Client, name string){test.Clear},
	}, config)
	if err != nil {
		log.Fatal(err)
	}

	demoUser := &model.UserResponse{ID: "demo", Roles: []string{role.Admin}}
	ctx := busdb.UserContext(context.Background(), demoUser)
	if err := test.SetupTestData(ctx, theCatalyst.DB); err != nil {
		log.Fatal(err)
	}

	// proxy static requests
	middlewares := []func(next http.Handler) http.Handler{
		catalyst.Authenticate(theCatalyst.DB, config.Auth),
		catalyst.AuthorizeBlockedUser(),
	}
	theCatalyst.Server.With(middlewares...).NotFound(api.Proxy("http://localhost:8080"))

	if err := http.ListenAndServe(":8000", theCatalyst.Server); err != nil {
		log.Fatal(err)
	}
}
