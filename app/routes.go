package app

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/martian/v3/cors"

	"github.com/SecurityBrewery/catalyst/app/auth"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/openapi"
	"github.com/SecurityBrewery/catalyst/ui"
)

func (a *App) SetupRoutes() error {
	// middleware for the router
	a.Router.Use(func(next http.Handler) http.Handler {
		return http.Handler(cors.NewHandler(next))
	})
	a.Router.Use(demoMode(a.Queries))
	a.Router.Use(middleware.RequestID)
	a.Router.Use(middleware.RealIP)
	a.Router.Use(middleware.Logger)
	a.Router.Use(middleware.Timeout(time.Second * 60))
	a.Router.Use(middleware.Recoverer)

	// base routes
	a.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui/", http.StatusFound)
	})
	a.Router.Get("/ui/*", a.staticFiles)
	a.Router.Get("/health", a.healthHandler)
	a.Router.Get("/config", a.configHandler)

	// auth routes
	a.Router.Mount("/auth", a.Auth.Server())

	// API routes
	apiHandler := openapi.Handler(openapi.NewStrictHandler(a.Service, []openapi.StrictMiddlewareFunc{auth.ValidateScopes}))
	a.Router.With(a.Auth.Middleware).Mount("/api", http.StripPrefix("/api", apiHandler))

	return nil
}

func (a *App) staticFiles(w http.ResponseWriter, r *http.Request) {
	if devServer := os.Getenv("UI_DEVSERVER"); devServer != "" {
		u, _ := url.Parse(devServer)

		r.Host = r.URL.Host

		httputil.NewSingleHostReverseProxy(u).ServeHTTP(w, r)

		return
	}

	fsys := http.FS(ui.UI())

	path := strings.TrimPrefix(r.URL.Path, "/ui")
	if path == "" || path == "/" {
		path = "/index.html"
	}

	if _, err := fsys.Open(strings.TrimPrefix(path, "/")); err != nil {
		path = "/index.html"
	}

	r.URL.Path = path
	http.FileServer(fsys).ServeHTTP(w, r)
}

func (a *App) healthHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := a.Queries.ListFeatures(r.Context(), sqlc.ListFeaturesParams{Offset: 0, Limit: 100}); err != nil {
		slog.ErrorContext(r.Context(), "Failed to get flags", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte("OK"))
}

func (a *App) configHandler(w http.ResponseWriter, r *http.Request) {
	features, err := a.Queries.ListFeatures(r.Context(), sqlc.ListFeaturesParams{Offset: 0, Limit: 100})
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to get flags", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	flags := make([]string, 0, len(features))
	for _, feature := range features {
		flags = append(flags, feature.Name)
	}

	b, _ := json.Marshal(map[string]any{ //nolint:errchkjson
		"flags": flags,
	})

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(b)
}
