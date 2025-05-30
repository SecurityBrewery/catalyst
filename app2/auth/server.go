package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Service) Server() http.Handler {
	router := chi.NewRouter()

	router.Get("/user", s.HandleUser)
	router.Get("/config", func(writer http.ResponseWriter, request *http.Request) {
		b, _ := json.Marshal(map[string]any{
			"password": s.config.PasswordAuth,
			"oidc":     s.config.OIDCAuth,
		})

		_, _ = writer.Write(b)
	})

	if s.config.PasswordAuth {
		router.Post("/local/login", s.handleLogin)
		router.Get("/local/logout", s.handleLogout)
	}

	if s.config.OIDCAuth {
		router.Get("/oidc/login", s.oidcLogin)
		router.Get("/oidc/callback", s.oidcCallback)
	}

	return router
}

func (s *Service) HandleUser(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	user, authErr, err := s.SessionManager.Get(r.Context())
	if err != nil {
		scimError(w, http.StatusInternalServerError, err.Error())

		return
	}

	if authErr != nil {
		_, _ = w.Write([]byte("null"))

		return
	}

	b, err := json.Marshal(user)
	if err != nil {
		scimError(w, http.StatusInternalServerError, err.Error())

		return
	}

	_, _ = w.Write(b)
}
