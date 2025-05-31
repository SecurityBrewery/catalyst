package app2

import (
	"context"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/SecurityBrewery/catalyst/app2/database/sqlc"
)

func demomode(queries *sqlc.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if criticalPath(r) && criticalMethod(r) && demoMode(r.Context(), queries) {
				http.Error(w, "Cannot modify reactions or files in demo mode", http.StatusForbidden)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func criticalPath(r *http.Request) bool {
	// Define critical paths that should not be accessed in demo mode
	criticalPaths := []string{"/api/reactions", "/api/files"}

	for _, path := range criticalPaths {
		if strings.Contains(r.URL.Path, path) {
			return true
		}
	}

	return false
}

func criticalMethod(r *http.Request) bool {
	return !slices.Contains([]string{http.MethodHead, http.MethodGet}, r.Method)
}

func demoMode(ctx context.Context, queries *sqlc.Queries) bool {
	features, err := queries.ListFeatures(ctx, sqlc.ListFeaturesParams{Offset: 0, Limit: 100})
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get flags", "error", err)

		return false
	}

	for _, feature := range features {
		if feature.Name == "demo" {
			return true
		}
	}

	return false
}
