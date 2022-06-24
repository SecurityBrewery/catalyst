package cmd

import (
	"fmt"
	"github.com/SecurityBrewery/catalyst/key"

	"github.com/alecthomas/kong"
	kongyaml "github.com/alecthomas/kong-yaml"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/exp/slices"
	"golang.org/x/oauth2"

	"github.com/SecurityBrewery/catalyst"
	"github.com/SecurityBrewery/catalyst/auth"
	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/role"
	"github.com/SecurityBrewery/catalyst/storage"
)

type cli struct {
	Secret           string `env:"SECRET"              required:""                    help:"A random secret value (can be created with 'openssl rand -hex 32')"`
	ExternalAddress  string `env:"EXTERNAL_ADDRESS"    required:""                    help:"The external address of the Catalyst server"`
	CatalystAddress  string `env:"CATALYST_ADDRESS"    default:"http://catalyst:8000" help:"The internal address of the Catalyst server"`
	Network          string `env:"CATALYST_NETWORK"    default:"catalyst"             help:"The network of the Catalyst server"`
	Port             int    `env:"PORT"                default:"8000"                 help:"The port of the Catalyst server"`
	InitialAPIKey    string `env:"INITIAL_API_KEY"     required:""                    help:"Setup an initial API key "`
	SimpleAuthEnable bool   `env:"SIMPLE_AUTH_ENABLE"  default:"true"                 help:"Enable simple authentication"`
	APIKeyAuthEnable bool   `env:"API_KEY_AUTH_ENABLE" default:"true"                 help:"Enable API key authentication"`
	OIDCAuthEnable   bool   `env:"OIDC_AUTH_ENABLE"    default:"true"                 help:"Enable OIDC authentication"`
	IndexPath        string `env:"INDEX_PATH"          default:"index.bleve"          help:"Path for the bleve index"`
	ArangoDBHost     string `env:"ARANGO_DB_HOST"      default:"http://arangodb:8529" help:"The host of the ArangoDB server"`
	ArangoDBUser     string `env:"ARANGO_DB_USER"      default:"root"                 help:"The user of the ArangoDB server"`
	ArangoDBPassword string `env:"ARANGO_DB_PASSWORD"  required:""                    help:"The password of the ArangoDB server"`
	S3Host           string `env:"S3_HOST"             default:"http://minio:9000"    help:"The host of the S3 server"     name:"s3-host"`
	S3User           string `env:"S3_USER"             default:"minio"                help:"The user of the S3 server"     name:"s3-user"`
	S3Password       string `env:"S3_PASSWORD"         required:""                    help:"The password of the S3 server" name:"s3-password"`

	OIDCIssuer        string   `env:"OIDC_ISSUER"         default:""                   help:"The url of the OIDC provider"`
	OIDCClientID      string   `env:"OIDC_CLIENT_ID"      default:"catalyst"           help:"The client ID for OIDC"`
	OIDCClientSecret  string   `env:"OIDC_CLIENT_SECRET"  default:""                   help:"The client secret for OIDC"`
	OIDCScopes        []string `env:"OIDC_SCOPES"                                      help:"Additional scopes to request, ['oidc', 'profile', 'email'] are always added." placeholder:"customscopes"`
	OIDCClaimUsername string   `env:"OIDC_CLAIM_USERNAME" default:"preferred_username" help:"Username field in the OIDC claim"`
	OIDCClaimEmail    string   `env:"OIDC_CLAIM_EMAIL"    default:"email"              help:"Email field in the OIDC claim"`
	OIDCClaimName     string   `env:"OIDC_CLAIM_NAME"     default:"name"               help:"Name field in the OIDC claim"`
	AuthBlockNew      bool     `env:"AUTH_BLOCK_NEW"      default:"true"               help:"Block newly created users"`
	AuthDefaultRoles  []string `env:"AUTH_DEFAULT_ROLES"                               help:"Default roles for new users"`
	AuthAdminUsers    []string `env:"AUTH_ADMIN_USERS"                                 help:"Usernames to grant admin rights"`
}

func (c *cli) Validate() error {
	switch {
	case c.OIDCAuthEnable && c.OIDCIssuer == "":
		return fmt.Errorf("missing flags: --oidc-issuer=STRING (is required when oidc-auth-enable is true)")
	case c.OIDCAuthEnable && c.OIDCClientSecret == "":
		return fmt.Errorf("missing flags: --oidc-client-secret=STRING (is required when oidc-auth-enable is true)")
	}

	return nil
}

func ParseCatalystConfig() (*catalyst.Config, error) {
	var cli cli

	kong.Parse(
		&cli,
		kong.Configuration(kong.JSON, "/etc/catalyst.json", ".catalyst.json"),
		kong.Configuration(kongyaml.Loader, "/etc/catalyst.yaml", ".catalyst.yaml"),
	)

	return MapConfig(&cli)
}

func MapConfig(cli *cli) (*catalyst.Config, error) {
	if cli.InitialAPIKey == "" {
		cli.InitialAPIKey = key.Generate()
	}

	config := &catalyst.Config{
		IndexPath: cli.IndexPath,
		Network:   cli.Network,
		DB: &database.Config{
			Host:     cli.ArangoDBHost,
			User:     cli.ArangoDBUser,
			Password: cli.ArangoDBPassword,
		},
		Storage: &storage.Config{
			Host:     cli.S3Host,
			User:     cli.S3User,
			Password: cli.S3Password,
		},
		Secret:          []byte(cli.Secret),
		ExternalAddress: cli.ExternalAddress,
		InternalAddress: cli.CatalystAddress,
		Port:            cli.Port,
		Auth: &auth.Config{
			SimpleAuthEnable: cli.SimpleAuthEnable,
			APIKeyAuthEnable: cli.APIKeyAuthEnable,
			OIDCAuthEnable:   cli.OIDCAuthEnable,
		},
		InitialAPIKey: cli.InitialAPIKey,
	}

	if cli.OIDCAuthEnable {
		roles := role.Explode(role.Analyst)
		roles = append(roles, role.Explodes(cli.AuthDefaultRoles)...)
		roles = role.Explodes(role.Strings(roles))

		scopes := slices.Compact(append([]string{oidc.ScopeOpenID, "profile", "email"}, cli.OIDCScopes...))

		config.Auth.OIDC = &auth.OIDCConfig{
			Issuer:           cli.OIDCIssuer,
			AuthBlockNew:     cli.AuthBlockNew,
			AuthDefaultRoles: roles,
			AuthAdminUsers:   cli.AuthAdminUsers,
			ClaimUsername:    cli.OIDCClaimUsername,
			ClaimEmail:       cli.OIDCClaimEmail,
			ClaimName:        cli.OIDCClaimName,
			OAuth2: &oauth2.Config{
				ClientID:     cli.OIDCClientID,
				ClientSecret: cli.OIDCClientSecret,
				RedirectURL:  cli.ExternalAddress + "/auth/callback",
				Scopes:       scopes,
			},
		}
	}

	return config, nil
}
