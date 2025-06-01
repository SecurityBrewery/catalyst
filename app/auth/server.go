package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (s *Service) Server() http.Handler {
	router := chi.NewRouter()

	router.Get("/user", s.handleUser)
	router.Get("/config", func(writer http.ResponseWriter, _ *http.Request) {
		b, _ := json.Marshal(map[string]any{
			"oidc": s.config.OIDCAuth,
		})

		_, _ = writer.Write(b)
	})

	router.Post("/local/login", s.handleLogin)
	router.Post("/local/reset-password-request", s.handlePasswordResetRequest)
	router.Get("/local/reset-password", s.handlePasswordReset)
	router.Post("/local/reset-password", s.handlePasswordResetPost)

	if s.config.OIDCAuth {
		router.Get("/oidc/login", s.oidcLogin)
		router.Get("/oidc/callback", s.oidcCallback)
	}

	return router
}

func (s *Service) handleUser(w http.ResponseWriter, r *http.Request) {
	authorizationHeader := r.Header.Get("Authorization")
	bearerToken := strings.TrimPrefix(authorizationHeader, bearerPrefix)

	user, err := s.verifyAuthToken(r.Context(), bearerToken)
	if err != nil {
		_, _ = w.Write([]byte("null"))

		return
	}

	b, err := json.Marshal(user)
	if err != nil {
		scimError(w, http.StatusInternalServerError, err.Error())

		return
	}

	r.Header.Set("Content-Type", "application/json")

	_, _ = w.Write(b)
}
