package auth

import (
	"encoding/base64"
	"errors"
	"log/slog"
	"net/http"
	"strings"
)

const bearerPrefix = "Bearer "

func (s *Service) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Authorization"), bearerPrefix) {
			// delegate to session auth middleware
			s.SessionAuth(next).ServeHTTP(w, r)
		} else {
			// delegate to bearer auth middleware
			s.BearerAuth(next).ServeHTTP(w, r)
		}
	})
}

func (s *Service) SessionAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// password auth is disabled
		if !s.config.PasswordAuth {
			scimUnauthorized(w, "Password auth is disabled")

			return
		}

		user, authErr, err := s.SessionManager.Get(r.Context())
		if err != nil {
			slog.ErrorContext(r.Context(), "SessionAuth", "error", err)

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)

			return
		}

		if authErr != nil {
			scimUnauthorized(w, authErr.Error())

			return
		}

		if !user.Verified {
			scimUnauthorized(w, "User is not verified")

			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Service) BearerAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// bearer auth is disabled
		if !s.config.BearerAuth {
			scimUnauthorized(w, "Bearer auth is disabled")

			return
		}

		authorizationHeader := r.Header.Get("Authorization")

		bearerToken := strings.TrimPrefix(authorizationHeader, bearerPrefix)

		token, err := base64.StdEncoding.DecodeString(bearerToken)
		if err != nil {
			scimUnauthorized(w, "Invalid token")

			return
		}

		username, password, ok := strings.Cut(string(token), ":")
		if !ok {
			scimUnauthorized(w, "Invalid token")

			return
		}

		if err := s.loginWithUsername(r.Context(), username, password); err != nil {
			if errors.Is(err, ErrUserInactive) {
				scimUnauthorized(w, "User is inactive")

				return
			}

			scimUnauthorized(w, "Login failed")

			return
		}

		next.ServeHTTP(w, r)
	})
}
