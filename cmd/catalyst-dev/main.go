package main

import (
	"context"
	"log"
	"net/http/httputil"
	"net/url"

	"github.com/arangodb/go-driver"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/SecurityBrewery/catalyst"
	"github.com/SecurityBrewery/catalyst/cmd"
	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/models"
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

	demoUser := &models.UserResponse{ID: "demo", Roles: []string{role.Admin}}
	ctx := busdb.UserContext(context.Background(), demoUser)
	if err := test.SetupTestData(ctx, theCatalyst.DB); err != nil {
		log.Fatal(err)
	}

	// proxy static requests
	theCatalyst.Server.NoRoute(
		sessions.Sessions(catalyst.SessionName, cookie.NewStore(config.Secret)),
		catalyst.Authenticate(theCatalyst.DB, config.Auth),
		catalyst.AuthorizeBlockedUser,
		proxy,
	)

	if err = theCatalyst.Server.RunWithSigHandler(); err != nil {
		log.Fatal(err)
	}
}

func proxy(ctx *gin.Context) {
	u, _ := url.Parse("http://localhost:8080")
	proxy := httputil.NewSingleHostReverseProxy(u)

	ctx.Request.Host = ctx.Request.URL.Host

	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}
