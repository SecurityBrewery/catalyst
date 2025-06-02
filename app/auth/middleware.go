package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"

	"github.com/SecurityBrewery/catalyst/app/auth/usercontext"
	"github.com/SecurityBrewery/catalyst/app/openapi"
)

const bearerPrefix = "Bearer "

func (s *Service) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		bearerToken := strings.TrimPrefix(authorizationHeader, bearerPrefix)

		user, claims, err := s.verifyAccessToken(r.Context(), bearerToken)
		if err != nil {
			slog.ErrorContext(r.Context(), "invalid bearer token", "error", err)

			scimUnauthorized(w, "invalid bearer token")

			return
		}

		scopes, err := scopes(claims)
		if err != nil {
			slog.ErrorContext(r.Context(), "failed to get scopes from token", "error", err)

			scimUnauthorized(w, "failed to get scopes")

			return
		}

		// Set the user in the context
		r = usercontext.UserRequest(r, user)
		r = usercontext.PermissionRequest(r, scopes)

		next.ServeHTTP(w, r)
	})
}

func (s *Service) ValidateScopes(next strictnethttp.StrictHTTPHandlerFunc, _ string) strictnethttp.StrictHTTPHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (response interface{}, err error) {
		requiredScopes, err := requiredScopes(r)
		if err != nil {
			slog.ErrorContext(ctx, "failed to get required scopes", "error", err)

			scimUnauthorized(w, "failed to get required scopes")

			return
		}

		if len(requiredScopes) > 0 {
			permissions, ok := usercontext.PermissionFromContext(ctx)
			if !ok {
				slog.ErrorContext(ctx, "missing permissions")
				scimUnauthorized(w, "missing permissions")

				return
			}

			if !hasScope(permissions, requiredScopes) {
				slog.ErrorContext(ctx, "missing required scopes", "required", requiredScopes, "permissions", permissions)

				scimUnauthorized(w, "missing required scopes")

				return
			}
		}

		return next(ctx, w, r, request)
	}
}

func requiredScopes(r *http.Request) ([]string, error) {
	requiredScopesValue := r.Context().Value(openapi.OAuth2Scopes)
	if requiredScopesValue == nil {
		slog.InfoContext(r.Context(), "no required scopes", "request", r.URL.Path)

		return nil, nil
	}

	requiredScopes, ok := requiredScopesValue.([]string)
	if !ok {
		return nil, fmt.Errorf("invalid required scopes type: %T", requiredScopesValue)
	}

	return requiredScopes, nil
}

func hasScope(scopes []string, requiredScopes []string) bool {
	for _, s := range requiredScopes {
		if !slices.Contains(scopes, s) {
			return false
		}
	}

	return true
}
