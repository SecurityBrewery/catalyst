package example

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/oauth2"

	"github.com/cugu/maut/auth"
)

type OIDCProvider interface {
	AddClient(id, redirectURL string) string
	Close()
	URL() string
}

func newAuthTestServerPort(t *testing.T, config *auth.Config, oidcProvider OIDCProvider, host, port string) *httptest.Server {
	t.Helper()

	// init test server
	testServer := httptest.NewUnstartedServer(nil)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		t.Fatal(err)
	}
	testServer.Listener = listener

	redirectURL := fmt.Sprintf("http://%s:%s/auth/callback", host, port)
	secret := oidcProvider.AddClient(fmt.Sprintf("apitest-%s", port), redirectURL)

	config.OIDCIssuer = oidcProvider.URL()
	config.OAuth2.ClientID = fmt.Sprintf("apitest-%s", port)
	config.OAuth2.ClientSecret = secret
	config.OAuth2.RedirectURL = redirectURL

	testServer.Config.Handler = newAuthServer(t, config)

	return testServer
}

func newAuthServer(t *testing.T, config *auth.Config) *chi.Mux {
	t.Helper()

	ctx := context.Background()
	resolver := auth.NewMemoryResolver()
	err := resolver.RoleCreateIfNotExists(ctx, &auth.Role{
		Name:        "analyst",
		Permissions: []string{"automation:read"},
	})
	if err != nil {
		t.Error(err)
	}
	authenticator, err := auth.NewAuthenticator(ctx, config, resolver)
	if err != nil {
		t.Fatal(err)
	}
	server := chi.NewRouter()
	server.Use(middleware.RequestID)
	server.Use(middleware.RealIP)
	server.Use(middleware.Logger)
	server.Use(middleware.Recoverer)

	server.
		With(
			authenticator.Authenticate(),
			authenticator.AuthorizeBlockedUser(),
			authenticator.AuthorizePermission("automation:read"),
		).
		Get("/", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("success"))
		})
	server.Mount("/auth", authenticator.Server())

	return server
}

func TestOIDC(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	t.Parallel()
	setup := []struct {
		name      string
		oidcSetup func() OIDCProvider
		authSetup func(*testing.T, *auth.Config, OIDCProvider, string, string) *httptest.Server
		host      string
		portLow   int
		login     func(client *http.Client, port string, initialResp *http.Response) (*http.Response, error)
	}{
		{name: "oidc internal", oidcSetup: newInternalOIDCServer, authSetup: newAuthTestServerPort, host: "127.0.0.1", portLow: 90, login: internalLogin},
		{name: "keycloak", oidcSetup: newKeycloakServer, authSetup: newAuthTestServerPort, host: "127.0.0.1", portLow: 91, login: keycloakLogin},
		{name: "authelia", oidcSetup: newAutheliaServer, authSetup: newAuthTestServerPort, host: "localhost", portLow: 92, login: autheliaLogin},
	}
	tests := []struct {
		name     string
		portHigh int
		config   *auth.Config
		want     *HTTPResponse
	}{
		{
			name:     "success",
			portHigh: 90,
			config: &auth.Config{
				CookieSecret:   []byte("secret"),
				OIDCAuthEnable: true,
				OAuth2: &oauth2.Config{
					Scopes: []string{oidc.ScopeOpenID, "profile", "email"}, // TODO: add groups
				},
				UserCreateConfig: &auth.UserCreateConfig{
					OIDCClaimUsername: "preferred_username",
					OIDCClaimEmail:    "email",
					OIDCClaimName:     "name",
					OIDCClaimGroups:   "groups",
					AuthDefaultRoles:  []string{"analyst"},
				},
			},
			want: &HTTPResponse{
				StatusCode: http.StatusOK,
				Body:       "success",
			},
		},
		{
			name:     "block new",
			portHigh: 91,
			config: &auth.Config{
				CookieSecret:   []byte("secret"),
				OIDCAuthEnable: true,
				OAuth2: &oauth2.Config{
					Scopes: []string{oidc.ScopeOpenID, "profile", "email"}, // TODO: add groups
				},
				UserCreateConfig: &auth.UserCreateConfig{
					OIDCClaimUsername: "preferred_username",
					OIDCClaimEmail:    "email",
					OIDCClaimName:     "name",
					OIDCClaimGroups:   "groups",
					AuthDefaultRoles:  []string{"analyst"},
					AuthBlockNew:      true,
				},
			},
			want: &HTTPResponse{
				StatusCode: http.StatusForbidden,
				Body:       `{"error":"user is blocked"}`,
			},
		},
		{
			name:     "invalid default role",
			portHigh: 92,
			config: &auth.Config{
				CookieSecret:   []byte("secret"),
				OIDCAuthEnable: true,
				OAuth2: &oauth2.Config{
					Scopes: []string{oidc.ScopeOpenID, "profile", "email"}, // TODO: add groups
				},
				UserCreateConfig: &auth.UserCreateConfig{
					OIDCClaimUsername: "preferred_username",
					OIDCClaimEmail:    "email",
					OIDCClaimName:     "name",
					OIDCClaimGroups:   "groups",
					AuthDefaultRoles:  []string{"viewer"},
					AuthBlockNew:      true,
				},
			},
			want: &HTTPResponse{
				StatusCode: http.StatusForbidden,
				Body:       `{"error":"user is blocked"}`,
			},
		},
	}

	for _, su := range setup {
		su := su
		t.Run(su.name, func(t *testing.T) {
			t.Parallel()

			// create oidc test server
			oidcServer := su.oidcSetup()
			defer oidcServer.Close()

			for _, tt := range tests {
				tt := tt
				t.Run(tt.name, func(t *testing.T) {
					// create test server
					port := fmt.Sprint(tt.portHigh*100 + su.portLow)
					authServer := su.authSetup(t, tt.config.Clone(), oidcServer, su.host, port)
					authServer.Start()
					defer authServer.Close()

					// create cookie jar
					client, err := testClient(t)
					if err != nil {
						t.Error(err)
					}

					// perform initial request
					initialResp, err := client.Get(authServer.URL + "/")
					if err != nil {
						t.Fatal(err)
					}
					defer initialResp.Body.Close()

					// send password
					loginResp, err := su.login(client, port, initialResp)
					if err != nil {
						t.Fatal(err)
					}
					defer loginResp.Body.Close()

					assertResult(t, loginResp, tt.want)
				})
			}
		})
	}
}
