package catalyst

import (
	"context"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/SecurityBrewery/catalyst/automation"
	"github.com/SecurityBrewery/catalyst/bus"
	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/models"
	"github.com/SecurityBrewery/catalyst/generated/restapi"
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

	UISettings      *models.Settings
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
	Server  *restapi.Server
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

	err = automation.New(config.Bus.APIUrl, config.InitialAPIKey, catalystBus, catalystDatabase)
	if err != nil {
		return nil, err
	}

	catalystService, err := service.New(catalystBus, catalystDatabase, catalystStorage, config.UISettings)
	if err != nil {
		return nil, err
	}

	if config.InitialAPIKey != "" {
		_ = catalystDatabase.UserDelete(ctx, "setup")

		ctx = busdb.UserContext(ctx, &models.UserResponse{
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

func setupAPI(catalystService *service.Service, catalystStorage *storage.Storage, catalystDatabase *database.Database, dbConfig *database.Config, bus *bus.Bus, config *Config) (*restapi.Server, error) {
	// session
	store := cookie.NewStore(config.Secret)
	setSession := sessions.Sessions(SessionName, store)

	authenticate := Authenticate(catalystDatabase, config.Auth)

	// create server
	apiServer := restapi.New(catalystService, &restapi.Config{Address: "0.0.0.0:8000", InsecureHTTP: true})
	apiServer.UseRawPath = true

	apiServer.ApiGroup.Use(setSession, authenticate, AuthorizeBlockedUser)
	apiServer.RoleAuth = AuthorizeRole

	apiServer.ConfigureRoutes()
	apiServer.ApiGroup.HEAD("/files/:ticketID/upload/:id", AuthorizeRole([]role.Role{role.FileReadWrite}), upload(catalystStorage.S3(), config.ExternalAddress))
	apiServer.ApiGroup.PATCH("/files/:ticketID/upload/:id", AuthorizeRole([]role.Role{role.FileReadWrite}), upload(catalystStorage.S3(), config.ExternalAddress))
	apiServer.ApiGroup.POST("/files/:ticketID/upload", AuthorizeRole([]role.Role{role.FileReadWrite}), upload(catalystStorage.S3(), config.ExternalAddress))
	apiServer.ApiGroup.GET("/files/:ticketID/download/:key", AuthorizeRole([]role.Role{role.FileReadWrite}), download(catalystStorage.Downloader()))

	apiServer.ApiGroup.GET("/backup/create", AuthorizeRole([]role.Role{role.BackupRead}), BackupHandler(catalystStorage, dbConfig))
	apiServer.ApiGroup.POST("/backup/restore", AuthorizeRole([]role.Role{role.BackupRestore}), RestoreHandler(catalystStorage, catalystDatabase, dbConfig))

	apiServer.GET("/callback", setSession, callback(config.Auth))
	apiServer.Any("/wss", setSession, authenticate, AuthorizeBlockedUser, handleWebSocket(bus))
	apiServer.NoRoute(setSession, authenticate, AuthorizeBlockedUser, static)

	apiServer.Use(cors.Default())
	return apiServer, nil
}
