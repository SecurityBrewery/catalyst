package auth

import (
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/exp/slices"

	"github.com/cugu/maut/api"
)

func (a *Authenticator) Middleware(permissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return a.Authenticate()(
			a.AuthorizeBlockedUser()(
				a.AuthorizePermission(permissions...)(next),
			),
		)
	}
}

func (a *Authenticator) Authenticate() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			keyHeader := r.Header.Get("PRIVATE-TOKEN")
			authHeader := r.Header.Get("Authorization")

			switch {
			case keyHeader != "":
				if a.config.APIKeyAuthEnable {
					a.keyAuth(keyHeader)(next).ServeHTTP(w, r)
				} else {
					api.JSONErrorStatus(w, http.StatusUnauthorized, errors.New("API Key authentication not enabled"))
				}
			case authHeader != "":
				if a.config.OIDCAuthEnable {
					iss := a.config.OIDCIssuer
					a.bearerAuth(authHeader, iss)(next).ServeHTTP(w, r)
				} else {
					api.JSONErrorStatus(w, http.StatusUnauthorized, errors.New("OIDC authentication not enabled"))
				}
			default:
				a.sessionAuth()(next).ServeHTTP(w, r)
			}
		})
	}
}

func (a *Authenticator) AuthorizeBlockedUser() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, _, ok := UserFromContext(r.Context())
			if !ok {
				api.JSONErrorStatus(w, http.StatusInternalServerError, errors.New("no user in context"))

				return
			}

			if user.Blocked {
				api.JSONErrorStatus(w, http.StatusForbidden, errors.New("user is blocked"))

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (a *Authenticator) AuthorizePermission(permissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, userPermissions, ok := UserFromContext(r.Context())
			if !ok {
				api.JSONErrorStatus(w, http.StatusInternalServerError, errors.New("no user in context"))

				return
			}

			authorized := false
			if slices.Contains(user.Roles, "admin") {
				authorized = true
			} else {
				for _, permission := range permissions {
					if slices.Contains(userPermissions, permission) {
						authorized = true

						break
					}
				}
			}

			if !authorized {
				api.JSONErrorStatus(w, http.StatusForbidden, fmt.Errorf("missing permissions %s has %s", permissions, userPermissions))

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
