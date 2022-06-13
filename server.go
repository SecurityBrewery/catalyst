package catalyst

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/SecurityBrewery/catalyst/auth"
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
	Auth            *auth.Config
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
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	if config.Auth.OIDCEnable {
		if err := config.Auth.Load(ctx); err != nil {
			return nil, err
		}
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
	secureJar := auth.NewJar(config.Secret)

	middlewares := []func(next http.Handler) http.Handler{auth.Authenticate(catalystDatabase, config.Auth, secureJar), auth.AuthorizeBlockedUser()}

	// create server
	apiServer := api.NewServer(catalystService, auth.AuthorizeRole, middlewares...)
	fileReadWrite := auth.AuthorizeRole([]string{role.FileReadWrite.String()})
	tudHandler := tusdUpload(catalystDatabase, bus, catalystStorage.S3(), config.ExternalAddress)
	apiServer.With(fileReadWrite).Head("/files/{ticketID}/tusd/{id}", tudHandler)
	apiServer.With(fileReadWrite).Patch("/files/{ticketID}/tusd/{id}", tudHandler)
	apiServer.With(fileReadWrite).Post("/files/{ticketID}/tusd", tudHandler)
	apiServer.With(fileReadWrite).Post("/files/{ticketID}/upload", upload(catalystDatabase, catalystStorage.S3(), catalystStorage.Uploader()))
	apiServer.With(fileReadWrite).Get("/files/{ticketID}/download/{key}", download(catalystStorage.Downloader()))

	apiServer.With(auth.AuthorizeRole([]string{role.BackupRead.String()})).Get("/backup/create", backupHandler(catalystStorage, dbConfig))
	apiServer.With(auth.AuthorizeRole([]string{role.BackupRestore.String()})).Post("/backup/restore", restoreHandler(catalystStorage, catalystDatabase, dbConfig))

	server := chi.NewRouter()
	server.Use(middleware.RequestID, middleware.RealIP, middleware.Logger, middleware.Recoverer)
	server.Mount("/api", apiServer)
	server.With(middlewares...).Handle("/wss", handleWebSocket(bus))
	server.Mount("/auth", auth.Server(config.Auth, catalystDatabase, secureJar))

	server.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui/", http.StatusFound)
	})

	return server, nil
}
