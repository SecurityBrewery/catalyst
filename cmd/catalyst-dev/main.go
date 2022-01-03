package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/SecurityBrewery/catalyst"
	"github.com/SecurityBrewery/catalyst/cmd"
	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/model"
	"github.com/SecurityBrewery/catalyst/hooks"
	"github.com/SecurityBrewery/catalyst/role"
	"github.com/SecurityBrewery/catalyst/test"
	"github.com/arangodb/go-driver"
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
	theCatalyst.Server.NotFound(proxy)
	// theCatalyst.Server.NoRoute(
	// 	sessions.Sessions(catalyst.SessionName, cookie.NewStore(config.Secret)),
	// 	catalyst.Authenticate(theCatalyst.DB, config.Auth),
	// 	catalyst.AuthorizeBlockedUser,
	// 	proxy,
	// )

	if err := http.ListenAndServe(":8000", theCatalyst.Server); err != nil {
		log.Fatal(err)
	}
}

func proxy(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse("http://localhost:8080")
	proxy := httputil.NewSingleHostReverseProxy(u)

	r.Host = r.URL.Host

	proxy.ServeHTTP(w, r)
}
