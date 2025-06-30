package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"

	"github.com/SecurityBrewery/catalyst/app/auth/usercontext"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/openapi"
)

const bearerPrefix = "Bearer "

func Middleware(queries *sqlc.Queries) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/config" {
				next.ServeHTTP(w, r)

				return
			}

			authorizationHeader := r.Header.Get("Authorization")
			bearerToken := strings.TrimPrefix(authorizationHeader, bearerPrefix)

			user, claims, err := verifyAccessToken(r.Context(), bearerToken, queries)
			if err != nil {
				slog.ErrorContext(r.Context(), "invalid bearer token", "error", err)

				unauthorizedJSON(w, "invalid bearer token")

				return
			}

			scopes, err := scopes(claims)
			if err != nil {
				slog.ErrorContext(r.Context(), "failed to get scopes from token", "error", err)

				unauthorizedJSON(w, "failed to get scopes")

				return
			}

			// Set the user in the context
			r = usercontext.UserRequest(r, user)
			r = usercontext.PermissionRequest(r, scopes)

			next.ServeHTTP(w, r)
		})
	}
}

func ValidateFileScopes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requiredScopes := []string{"file:read"}
		if slices.Contains([]string{http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete}, r.Method) {
			requiredScopes = []string{"file:write"}
		}

		if err := validateScopes(r.Context(), requiredScopes); err != nil {
			slog.ErrorContext(r.Context(), "failed to validate scopes", "error", err)
			unauthorizedJSON(w, "missing required scopes")

			return
		}

		next.ServeHTTP(w, r)
	})
}

func ValidateScopesStrict(next strictnethttp.StrictHTTPHandlerFunc, _ string) strictnethttp.StrictHTTPHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (response any, err error) {
		requiredScopes, err := requiredScopes(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "failed to get required scopes", "error", err)
			unauthorizedJSON(w, "failed to get required scopes")

			return nil, fmt.Errorf("failed to get required scopes: %w", err)
		}

		if err := validateScopes(ctx, requiredScopes); err != nil {
			slog.ErrorContext(ctx, "failed to validate scopes", "error", err)
			unauthorizedJSON(w, "missing required scopes")

			return nil, fmt.Errorf("missing required scopes: %w", err)
		}

		return next(ctx, w, r, request)
	}
}

func LogError(next strictnethttp.StrictHTTPHandlerFunc, _ string) strictnethttp.StrictHTTPHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (response any, err error) {
		re, err := next(ctx, w, r, request)
		if err != nil {
			if err.Error() == "context canceled" {
				// This is a common error when the request is canceled, e.g., by the client.
				// We can ignore this error as it does not indicate a problem with the handler.
				return re, nil
			}

			slog.ErrorContext(ctx, "handler error", "error", err, "method", r.Method, "path", r.URL.Path)
		}

		return re, err
	}
}

func validateScopes(ctx context.Context, requiredScopes []string) error {
	if len(requiredScopes) > 0 {
		permissions, ok := usercontext.PermissionFromContext(ctx)
		if !ok {
			return errors.New("missing permissions")
		}

		if !hasScope(permissions, requiredScopes) {
			return fmt.Errorf("missing required scopes: %v", requiredScopes)
		}
	}

	return nil
}

func requiredScopes(ctx context.Context) ([]string, error) {
	requiredScopesValue := ctx.Value(openapi.OAuth2Scopes)
	if requiredScopesValue == nil {
		return nil, nil
	}

	requiredScopes, ok := requiredScopesValue.([]string)
	if !ok {
		return nil, fmt.Errorf("invalid required scopes type: %T", requiredScopesValue)
	}

	return requiredScopes, nil
}

func hasScope(scopes []string, requiredScopes []string) bool {
	if slices.Contains(scopes, "admin") {
		// If the user has admin scope, they can access everything
		return true
	}

	for _, s := range requiredScopes {
		if !slices.Contains(scopes, s) {
			return false
		}
	}

	return true
}
