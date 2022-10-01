package catalyst

import (
	"context"
	"net/http"
	"time"

	maut "github.com/cugu/maut/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/SecurityBrewery/catalyst/bus"
	"github.com/SecurityBrewery/catalyst/busservice"
	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/api"
	"github.com/SecurityBrewery/catalyst/generated/model"
	"github.com/SecurityBrewery/catalyst/hooks"
	"github.com/SecurityBrewery/catalyst/index"
	"github.com/SecurityBrewery/catalyst/role"
	"github.com/SecurityBrewery/catalyst/service"
	"github.com/SecurityBrewery/catalyst/storage"
)

type Config struct {
	IndexPath string
	DB        *database.Config
	Storage   *storage.Config

	Secret          []byte
	Auth            *maut.Config
	ExternalAddress string
	InternalAddress string
	InitialAPIKey   string
	Network         string
	Port            int
}

type Server struct {
	Bus     *bus.Bus
	DB      *database.Database
	Index   *index.Index
	Storage *storage.Storage
	Server  chi.Router
}

func New(hooks *hooks.Hooks, config *Config) (*Server, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Minute*10)
	defer cancel()

	/*
		if config.Auth.OIDCAuthEnable {
			if err := config.Auth.Load(ctx); err != nil {
				return nil, err
			}
		}
	*/
	a, err := maut.NewAuthenticator(ctx, config.Auth, nil)
	if err != nil {
		return nil, err
	}

	catalystStorage, err := storage.New(config.Storage)
	if err != nil {
		return nil, err
	}

	catalystIndex, err := index.New(config.IndexPath)
	if err != nil {
		return nil, err
	}

	catalystBus := bus.New()

	catalystDatabase, err := database.New(ctx, catalystIndex, catalystBus, hooks, config.DB)
	if err != nil {
		return nil, err
	}

	busservice.New(config.InternalAddress+"/api", config.InitialAPIKey, config.Network, catalystBus, catalystDatabase)

	catalystService, err := service.New(catalystBus, catalystDatabase, catalystStorage, GetVersion())
	if err != nil {
		return nil, err
	}

	if config.InitialAPIKey != "" {
		_ = catalystDatabase.UserDelete(ctx, "setup")

		ctx = busdb.UserContext(ctx, &model.UserResponse{
			ID:      "setup",
			Roles:   role.Strings(role.Explode(role.Admin)),
			Apikey:  false,
			Blocked: false,
		})
		_, err = catalystDatabase.UserCreateSetupAPIKey(ctx, config.InitialAPIKey)
		if err != nil {
			return nil, err
		}
	}

	apiServer, err := setupAPI(a, catalystService, catalystStorage, catalystDatabase, config.DB, catalystBus, config)
	if err != nil {
		return nil, err
	}

	return &Server{
		Bus:     catalystBus,
		DB:      catalystDatabase,
		Index:   catalystIndex,
		Storage: catalystStorage,
		Server:  apiServer,
	}, nil
}

func setupAPI(a *maut.Authenticator, catalystService *service.Service, catalystStorage *storage.Storage, catalystDatabase *database.Database, dbConfig *database.Config, bus *bus.Bus, config *Config) (chi.Router, error) {
	middlewares := []func(next http.Handler) http.Handler{
		a.Authenticate(),
		a.AuthorizeBlockedUser(),
	}

	// create server
	apiServer := api.NewServer(catalystService, autRole(a), middlewares...)
	fileReadWrite := a.AuthorizePermission(role.FileReadWrite.String())
	tudHandler := tusdUpload(catalystDatabase, bus, catalystStorage.S3(), config.ExternalAddress)
	apiServer.With(fileReadWrite).Head("/files/{ticketID}/tusd/{id}", tudHandler)
	apiServer.With(fileReadWrite).Patch("/files/{ticketID}/tusd/{id}", tudHandler)
	apiServer.With(fileReadWrite).Post("/files/{ticketID}/tusd", tudHandler)
	apiServer.With(fileReadWrite).Post("/files/{ticketID}/upload", upload(catalystDatabase, catalystStorage.S3(), catalystStorage.Uploader()))
	apiServer.With(fileReadWrite).Get("/files/{ticketID}/download/{key}", download(catalystStorage.Downloader()))

	apiServer.With(a.AuthorizePermission(role.BackupRead.String())).Get("/backup/create", backupHandler(catalystStorage, dbConfig))
	apiServer.With(a.AuthorizePermission(role.BackupRestore.String())).Post("/backup/restore", restoreHandler(catalystStorage, catalystDatabase, dbConfig))

	server := chi.NewRouter()
	server.Use(middleware.RequestID, middleware.RealIP, middleware.Logger, middleware.Recoverer)
	server.Mount("/api", apiServer)
	server.With(middlewares...).Handle("/wss", handleWebSocket(bus))
	server.Mount("/auth", a.Server())

	server.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui/", http.StatusFound)
	})

	return server, nil
}

func autRole(a *maut.Authenticator) func([]string) func(http.Handler) http.Handler {
	return func(permissions []string) func(http.Handler) http.Handler {
		return a.AuthorizePermission(permissions...)
	}
}
