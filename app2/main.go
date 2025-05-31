package app2

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/martian/v3/cors"

	"github.com/SecurityBrewery/catalyst/app2/auth"
	"github.com/SecurityBrewery/catalyst/app2/database"
	"github.com/SecurityBrewery/catalyst/app2/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app2/openapi"
)

func App(ctx context.Context, filename string, _ bool) (*App2, error) {
	queries, _, err := database.DB(filepath.Join(filename, "data.db"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	authService, err := auth.New(ctx, queries, &auth.Config{
		Domain:       "localhost",
		CookieSecure: false,
		PasswordAuth: true,
		BearerAuth:   true,
		OIDCAuth:     false,
		OIDCIssuer:   "",
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "",
		AuthURL:      "",
		UserCreateConfig: &auth.UserCreateConfig{
			Active:            false,
			OIDCClaimUsername: "",
			OIDCClaimEmail:    "",
			OIDCClaimName:     "",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create auth service: %w", err)
	}

	return &App2{
		Queries: queries,
		Router:  chi.NewRouter(),
		Service: NewService(queries),
		Auth:    authService,
	}, nil
}

type App2 struct {
	Queries *sqlc.Queries
	Router  *chi.Mux
	Service *Service
	Auth    *auth.Service
}

func (a *App2) SetupRoutes() error {
	a.Router.Use(func(next http.Handler) http.Handler {
		return http.Handler(cors.NewHandler(next))
	})
	a.Router.Use(demomode(a.Queries))
	a.Router.Use(a.Auth.SessionManager.LoadAndSave)
	a.Router.Use(middleware.Timeout(time.Second * 60))
	a.Router.Use(middleware.RequestID)
	a.Router.Use(middleware.RealIP)
	a.Router.Use(middleware.Logger)
	a.Router.Use(middleware.Recoverer)

	apiHandler := openapi.Handler(openapi.NewStrictHandler(a.Service, nil))
	a.Router.With(a.Auth.Middleware).Mount("/api", http.StripPrefix("/api", apiHandler))
	a.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui/", http.StatusFound)
	})
	a.Router.Get("/ui/*", staticFiles)
	a.Router.Mount("/auth", a.Auth.Server())
	a.Router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if _, err := a.Queries.ListFeatures(r.Context(), sqlc.ListFeaturesParams{Offset: 0, Limit: 100}); err != nil {
			slog.ErrorContext(r.Context(), "Failed to get flags", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
	a.Router.Get("/config", func(w http.ResponseWriter, r *http.Request) {
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

		b, _ := json.Marshal(map[string]any{
			"flags": flags,
		})

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(b)
	})

	return nil
}

func (a *App2) Start(ctx context.Context) error {
	err := a.SetupRoutes()
	if err != nil {
		return err
	}

	return http.ListenAndServe(":8090", a.Router)
}

func staticFiles(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse("http://localhost:3000/")

	r.Host = r.URL.Host

	httputil.NewSingleHostReverseProxy(u).ServeHTTP(w, r)
}
