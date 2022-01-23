package catalyst

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

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
	Bus       *bus.Config

	UISettings      *model.Settings
	Secret          []byte
	Auth            *AuthConfig
	ExternalAddress string
	InitialAPIKey   string
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
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	err := config.Auth.Load(ctx)
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

	catalystBus, err := bus.New(config.Bus)
	if err != nil {
		return nil, err
	}

	catalystDatabase, err := database.New(ctx, catalystIndex, catalystBus, hooks, config.DB)
	if err != nil {
		return nil, err
	}

	err = busservice.New(config.Bus.APIUrl, config.InitialAPIKey, catalystBus, catalystDatabase)
	if err != nil {
		return nil, err
	}

	catalystService, err := service.New(catalystBus, catalystDatabase, catalystStorage, config.UISettings)
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

	apiServer, err := setupAPI(catalystService, catalystStorage, catalystDatabase, config.DB, catalystBus, config)
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

func setupAPI(catalystService *service.Service, catalystStorage *storage.Storage, catalystDatabase *database.Database, dbConfig *database.Config, bus *bus.Bus, config *Config) (chi.Router, error) {
	// create server
	allowAll := cors.AllowAll().Handler
	apiServer := api.NewServer(
		catalystService,
		AuthorizeRole,
		allowAll, Authenticate(catalystDatabase, config.Auth), AuthorizeBlockedUser(),
	)

	apiServer.With(AuthorizeRole([]string{role.FileReadWrite.String()})).Head("/files/{ticketID}/tusd/{id}", tusdUpload(catalystDatabase, bus, catalystStorage.S3(), config.ExternalAddress))
	apiServer.With(AuthorizeRole([]string{role.FileReadWrite.String()})).Patch("/files/{ticketID}/tusd/{id}", tusdUpload(catalystDatabase, bus, catalystStorage.S3(), config.ExternalAddress))
	apiServer.With(AuthorizeRole([]string{role.FileReadWrite.String()})).Post("/files/{ticketID}/tusd", tusdUpload(catalystDatabase, bus, catalystStorage.S3(), config.ExternalAddress))
	apiServer.With(AuthorizeRole([]string{role.FileReadWrite.String()})).Post("/files/{ticketID}/upload", upload(catalystDatabase, catalystStorage.S3(), catalystStorage.Uploader()))
	apiServer.With(AuthorizeRole([]string{role.FileReadWrite.String()})).Get("/files/{ticketID}/download/{key}", download(catalystStorage.Downloader()))

	apiServer.With(AuthorizeRole([]string{role.BackupRead.String()})).Get("/backup/create", backupHandler(catalystStorage, dbConfig))
	apiServer.With(AuthorizeRole([]string{role.BackupRestore.String()})).Post("/backup/restore", restoreHandler(catalystStorage, catalystDatabase, dbConfig))

	server := chi.NewRouter()
	server.Use(middleware.RequestID, middleware.RealIP, middleware.Logger, middleware.Recoverer, allowAll)
	server.Mount("/api", apiServer)

	server.Get("/callback", callback(config.Auth))
	server.With(Authenticate(catalystDatabase, config.Auth), AuthorizeBlockedUser()).Handle("/wss", handleWebSocket(bus))
	// server.With(Authenticate(catalystDatabase, config.Auth), AuthorizeBlockedUser()).NotFound(static)
	server.NotFound(func(w http.ResponseWriter, r *http.Request) {
		log.Println("not found", r.URL.RawPath)
		Authenticate(catalystDatabase, config.Auth)(AuthorizeBlockedUser()(http.HandlerFunc(static))).ServeHTTP(w, r)
	})

	return server, nil
}
