package app2

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/martian/v3/cors"

	"github.com/SecurityBrewery/catalyst/app2/database"
	"github.com/SecurityBrewery/catalyst/app2/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app2/openapi"
)

func App(filename string, _ bool) (*App2, error) {
	queries, _, err := database.DB(filepath.Join(filename, "data.db"))
	if err != nil {
		return nil, err
	}

	return &App2{
		Queries: queries,
	}, nil
}

type App2 struct {
	Queries *sqlc.Queries
}

func (a *App2) Start() error {
	service := &Service{
		Queries: a.Queries,
	}

	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.Handler(cors.NewHandler(next))
	})
	r.Use(middleware.Timeout(time.Second * 60))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	apiHandler := openapi.Handler(openapi.NewStrictHandler(service, nil))
	r.Mount("/api", http.StripPrefix("/api", apiHandler))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui/", http.StatusFound)
	})
	r.Get("/ui/*", staticFiles)

	return http.ListenAndServe(":8090", r)
}

func staticFiles(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "staticFiles", "path", r.URL.Path)

	u, _ := url.Parse("http://localhost:3000/")

	r.Host = r.URL.Host

	httputil.NewSingleHostReverseProxy(u).ServeHTTP(w, r)
}
