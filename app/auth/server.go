package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/mail"
)

func Server(queries *sqlc.Queries, mailer *mail.Mailer) http.Handler {
	router := chi.NewRouter()

	router.Get("/user", handleUser(queries))
	router.Post("/local/login", handleLogin(queries))
	router.Post("/local/reset-password-mail", handleResetPasswordMail(queries, mailer))
	router.Post("/local/reset-password", handlePassword(queries))

	return router
}

func handleUser(queries *sqlc.Queries) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		bearerToken := strings.TrimPrefix(authorizationHeader, bearerPrefix)

		user, _, err := verifyAccessToken(r.Context(), bearerToken, queries)
		if err != nil {
			_, _ = w.Write([]byte("null"))

			return
		}

		permissions, err := queries.ListUserPermissions(r.Context(), user.ID)
		if err != nil {
			errorJSON(w, http.StatusInternalServerError, err.Error())

			return
		}

		b, err := json.Marshal(map[string]any{
			"user":        user,
			"permissions": permissions,
		})
		if err != nil {
			errorJSON(w, http.StatusInternalServerError, err.Error())

			return
		}

		r.Header.Set("Content-Type", "application/json")

		_, _ = w.Write(b)
	}
}
