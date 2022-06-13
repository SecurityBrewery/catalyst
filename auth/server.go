package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/SecurityBrewery/catalyst/database"
)

func Server(config *Config, catalystDatabase *database.Database, jar *Jar) *chi.Mux {
	server := chi.NewRouter()

	server.Get("/hasoidc", hasOIDC(config.OIDCEnable))

	if config.OIDCEnable {
		server.Get("/callback", callback(config, jar))
	}
	if config.SimpleAuthEnable {
		server.Post("/login", login(catalystDatabase, jar))
	}
	if config.SimpleAuthEnable && config.OIDCEnable {
		server.Get("/oidclogin", redirectToOIDCLogin(config, jar))
	}
	server.Post("/logout", logout())

	return server
}

func hasOIDC(oidcEnable bool) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		b, err := json.Marshal(map[string]any{
			"hasoidc": oidcEnable,
		})
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)

			return
		}

		_, _ = writer.Write(b)
	}
}
