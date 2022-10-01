package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/cugu/maut/api"
)

func (a *Authenticator) bearerAuth(authHeader string, iss string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.HasPrefix(authHeader, "Bearer ") {
				api.JSONErrorStatus(w, http.StatusUnauthorized, errors.New("no bearer token"))

				return
			}

			username, claims, apiError := a.verifyClaims(r, authHeader[7:])
			if apiError != nil {
				api.JSONErrorStatus(w, apiError.Status, apiError.Internal)

				return
			}

			if claims["iss"] != iss {
				api.JSONErrorStatus(w, http.StatusUnauthorized, fmt.Errorf("wrong issuer, expected %s, got %s", iss, claims["iss"]))

				return
			}

			r, err := a.setContextClaims(r.Context(), r, username)
			if err != nil {
				api.JSONErrorStatus(w, http.StatusInternalServerError, fmt.Errorf("could not load user: %w", err))

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (a *Authenticator) keyAuth(keyHeader string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key, err := a.resolver.UserAPIKeyByHash(r.Context(), keyHeader)
			if err != nil {
				api.JSONErrorStatus(w, http.StatusUnauthorized, fmt.Errorf("could not verify private token: %w", err))

				return
			}

			r = a.setUserContext(r, key)

			next.ServeHTTP(w, r)
		})
	}
}

func (a *Authenticator) sessionAuth() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, noCookie := a.jar.userSession(r)

			if noCookie {
				a.redirectToLogin()(w, r)

				return
			}

			r, err := a.setContextClaims(r.Context(), r, userID)
			if err != nil {
				api.JSONErrorStatus(w, http.StatusInternalServerError, fmt.Errorf("could not load user: %w", err))

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (a *Authenticator) redirectToLogin() func(http.ResponseWriter, *http.Request) {
	if a.config.SimpleAuthEnable {
		return func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}

	if a.config.OIDCAuthEnable {
		return a.redirectToOIDCLogin()
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		api.JSONErrorStatus(writer, http.StatusForbidden, errors.New("unauthenticated"))
	}
}
