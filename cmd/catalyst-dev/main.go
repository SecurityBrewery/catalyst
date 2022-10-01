package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/arangodb/go-driver"
	maut "github.com/cugu/maut/auth"

	"github.com/SecurityBrewery/catalyst"
	"github.com/SecurityBrewery/catalyst/cmd"
	"github.com/SecurityBrewery/catalyst/generated/api"
	"github.com/SecurityBrewery/catalyst/generated/model"
	"github.com/SecurityBrewery/catalyst/generated/pointer"
	"github.com/SecurityBrewery/catalyst/hooks"
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

	demoUser := &maut.User{ID: "demo", Roles: []string{maut.AdminRole}}
	ctx := maut.UserContext(context.Background(), demoUser, nil)
	if err := test.SetupTestData(ctx, theCatalyst.DB); err != nil {
		log.Fatal(err)
	}

	_, _ = theCatalyst.DB.UserCreate(context.Background(), &model.UserForm{ID: "eve", Roles: []string{"admin"}, Password: pointer.String("eve")})
	_ = theCatalyst.DB.UserDataCreate(context.Background(), "eve", &model.UserData{
		Name:  pointer.String("Eve"),
		Email: pointer.String("eve@example.com"),
		Image: &avatarEve,
	})
	_, _ = theCatalyst.DB.UserCreate(context.Background(), &model.UserForm{ID: "kevin", Roles: []string{"admin"}, Password: pointer.String("kevin")})
	_ = theCatalyst.DB.UserDataCreate(context.Background(), "kevin", &model.UserData{
		Name:  pointer.String("Kevin"),
		Email: pointer.String("kevin@example.com"),
		Image: &avatarKevin,
	})

	_, _ = theCatalyst.DB.UserCreate(context.Background(), &model.UserForm{ID: "tom", Roles: []string{"admin"}, Password: pointer.String("tom")})
	_ = theCatalyst.DB.UserDataCreate(context.Background(), "tom", &model.UserData{
		Name:  pointer.String("tom"),
		Email: pointer.String("tom@example.com"),
		Image: &avatarKevin,
	})

	// proxy static requests
	theCatalyst.Server.Get("/ui/*", func(writer http.ResponseWriter, request *http.Request) {
		log.Println("proxy request", request.URL.Path)

		api.Proxy("http://localhost:8080/")(writer, request)
	})

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", config.Port),
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           theCatalyst.Server,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
