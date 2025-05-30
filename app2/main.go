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
	"github.com/SecurityBrewery/catalyst/app2/fakedata"
	"github.com/SecurityBrewery/catalyst/app2/openapi"
)

func App(filename string, _ bool) (*App2, error) {
	queries, _, err := database.DB(filepath.Join(filename, "data.db"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := fakedata.Generate(queries, 1, 10); err != nil {
		return nil, fmt.Errorf("failed to generate fake data: %w", err)
	}

	if _, err := queries.CreateFeature(context.Background(), "dev"); err != nil {
		return nil, err
	}

	return &App2{
		Queries: queries,
	}, nil
}

type App2 struct {
	Queries *sqlc.Queries
}

func (a *App2) Start(ctx context.Context) error {
	service := &Service{
		Queries: a.Queries,
	}

	authService, err := auth.New(ctx, a.Queries, &auth.Config{
		Domain:       "localhost",
		CookieSecure: false,
		PasswordAuth: true,
		BearerAuth:   false,
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
		return fmt.Errorf("failed to create auth service: %w", err)
	}

	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.Handler(cors.NewHandler(next))
	})
	r.Use(authService.SessionManager.LoadAndSave)
	r.Use(middleware.Timeout(time.Second * 60))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	apiHandler := openapi.Handler(openapi.NewStrictHandler(service, nil))
	r.With(authService.Middleware).Mount("/api", http.StripPrefix("/api", apiHandler))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui/", http.StatusFound)
	})
	r.Get("/ui/*", staticFiles)
	r.Mount("/auth", authService.Server())
	r.Get("/config", func(w http.ResponseWriter, r *http.Request) {
		features, err := a.Queries.ListFeatures(r.Context())
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

	return http.ListenAndServe(":8090", r)
}

func staticFiles(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse("http://localhost:3000/")

	r.Host = r.URL.Host

	httputil.NewSingleHostReverseProxy(u).ServeHTTP(w, r)
}
