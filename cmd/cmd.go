package cmd

import (
	"github.com/alecthomas/kong"
	kongyaml "github.com/alecthomas/kong-yaml"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"github.com/SecurityBrewery/catalyst"
	"github.com/SecurityBrewery/catalyst/bus"
	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/role"
	"github.com/SecurityBrewery/catalyst/storage"
)

type CLI struct {
	Secret          string `env:"SECRET"           required:""  help:"A random secret value (can be created with 'openssl rand -hex 32')"`
	ExternalAddress string `env:"EXTERNAL_ADDRESS" required:""`
	CatalystAddress string `env:"CATALYST_ADDRESS" default:"http://catalyst:8000"`
	Network         string `env:"CATALYST_NETWORK" default:"catalyst"`

	OIDCIssuer        string   `env:"OIDC_ISSUER"         required:""`
	OIDCClientID      string   `env:"OIDC_CLIENT_ID"      default:"catalyst"`
	OIDCClientSecret  string   `env:"OIDC_CLIENT_SECRET"  required:""`
	OIDCScopes        []string `env:"OIDC_SCOPES"                                      help:"Additional scopes, ['oidc', 'profile', 'email'] are always added." placeholder:"customscopes"`
	OIDCClaimUsername string   `env:"OIDC_CLAIM_USERNAME" default:"preferred_username" help:"username field in the OIDC claim"`
	OIDCClaimEmail    string   `env:"OIDC_CLAIM_EMAIL"    default:"email"              help:"email field in the OIDC claim"`
	OIDCClaimName     string   `env:"OIDC_CLAIM_NAME"     default:"name"               help:"name field in the OIDC claim"`
	AuthBlockNew      bool     `env:"AUTH_BLOCK_NEW"      default:"true"               help:"Block newly created users"`
	AuthDefaultRoles  []string `env:"AUTH_DEFAULT_ROLES"                               help:"Default roles for new users"`
	AuthAdminUsers    []string `env:"AUTH_ADMIN_USERS"                                 help:"Username of admins"`

	IndexPath string `env:"INDEX_PATH" default:"index.bleve" help:"Path for the bleve index"`

	ArangoDBHost     string `env:"ARANGO_DB_HOST"     default:"http://arangodb:8529"`
	ArangoDBUser     string `env:"ARANGO_DB_USER"     default:"root"`
	ArangoDBPassword string `env:"ARANGO_DB_PASSWORD" required:""`

	S3Host     string `env:"S3_HOST"     default:"http://minio:9000" name:"s3-host"`
	S3User     string `env:"S3_USER"     default:"minio"             name:"s3-user"`
	S3Password string `env:"S3_PASSWORD" required:""                 name:"s3-password"`

	EmitterIOHost string `env:"EMITTER_IO_HOST" default:"tcp://emitter:8080"`
	EmitterIORKey string `env:"EMITTER_IO_KEY"  required:""`

	InitialAPIKey string `env:"INITIAL_API_KEY"`
}

func ParseCatalystConfig() (*catalyst.Config, error) {
	var cli CLI
	kong.Parse(
		&cli,
		kong.Configuration(kong.JSON, "/etc/catalyst.json", ".catalyst.json"),
		kong.Configuration(kongyaml.Loader, "/etc/catalyst.yaml", ".catalyst.yaml"),
	)

	return MapConfig(cli)
}

func MapConfig(cli CLI) (*catalyst.Config, error) {
	roles := role.Explode(role.Analyst)
	roles = append(roles, role.Explodes(cli.AuthDefaultRoles)...)
	roles = role.Explodes(role.Strings(roles))

	scopes := unique(append([]string{oidc.ScopeOpenID, "profile", "email"}, cli.OIDCScopes...))
	config := &catalyst.Config{
		IndexPath:       cli.IndexPath,
		Network:         cli.Network,
		DB:              &database.Config{Host: cli.ArangoDBHost, User: cli.ArangoDBUser, Password: cli.ArangoDBPassword},
		Storage:         &storage.Config{Host: cli.S3Host, User: cli.S3User, Password: cli.S3Password},
		Secret:          []byte(cli.Secret),
		ExternalAddress: cli.ExternalAddress,
		Auth: &catalyst.AuthConfig{
			OIDCIssuer:        cli.OIDCIssuer,
			OAuth2:            &oauth2.Config{ClientID: cli.OIDCClientID, ClientSecret: cli.OIDCClientSecret, RedirectURL: cli.ExternalAddress + "/callback", Scopes: scopes},
			OIDCClaimUsername: cli.OIDCClaimUsername,
			OIDCClaimEmail:    cli.OIDCClaimEmail,
			OIDCClaimName:     cli.OIDCClaimName,
			AuthBlockNew:      cli.AuthBlockNew,
			AuthDefaultRoles:  roles,
			AuthAdminUsers:    cli.AuthAdminUsers,
		},
		Bus:           &bus.Config{Host: cli.EmitterIOHost, Key: cli.EmitterIORKey, APIUrl: cli.CatalystAddress + "/api"},
		InitialAPIKey: cli.InitialAPIKey,
	}
	return config, nil
}

func unique(l []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range l {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
