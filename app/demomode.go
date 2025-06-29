package app

import (
	"context"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

func demoMode(queries *sqlc.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isCriticalPath(r) && isCriticalMethod(r) && isDemoMode(r.Context(), queries) {
				http.Error(w, "Cannot modify reactions or files in demo mode", http.StatusForbidden)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isCriticalPath(r *http.Request) bool {
	// Define critical paths that should not be accessed in demo mode
	criticalPaths := []string{
		"/api/files",
		"/api/groups",
		"/api/reactions",
		"/api/settings",
		"/api/users",
		"/api/webhooks",
	}

	for _, path := range criticalPaths {
		if strings.Contains(r.URL.Path, path) {
			return true
		}
	}

	return false
}

func isCriticalMethod(r *http.Request) bool {
	return !slices.Contains([]string{http.MethodHead, http.MethodGet}, r.Method)
}

func isDemoMode(ctx context.Context, queries *sqlc.Queries) bool {
	var demoMode bool

	if err := database.Paginate(ctx, func(ctx context.Context, offset, limit int64) (nextPage bool, err error) {
		features, err := queries.ListFeatures(ctx, sqlc.ListFeaturesParams{Offset: offset, Limit: limit})
		if err != nil {
			return false, err
		}

		for _, feature := range features {
			if feature.Key == "demo" {
				demoMode = true

				return false, nil // Stop pagination if demo mode is found
			}
		}

		return true, nil
	}); err != nil {
		slog.ErrorContext(ctx, "Failed to check demo mode", "error", err)

		return false
	}

	return demoMode
}
