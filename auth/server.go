package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/SecurityBrewery/catalyst/database"
)

func Server(config *Config, catalystDatabase *database.Database, jar *Jar) *chi.Mux {
	server := chi.NewRouter()

	server.Get("/config", hasOIDC(config))

	if config.OIDCAuthEnable {
		server.Get("/callback", callback(config, jar))
		server.Get("/oidclogin", redirectToOIDCLogin(config, jar))
	}
	if config.SimpleAuthEnable {
		server.Post("/login", login(catalystDatabase, jar))
	}
	server.Post("/logout", logout())

	return server
}

func hasOIDC(config *Config) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		b, err := json.Marshal(map[string]any{
			"simple": config.SimpleAuthEnable,
			"oidc":   config.OIDCAuthEnable,
		})
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)

			return
		}

		_, _ = writer.Write(b)
	}
}
