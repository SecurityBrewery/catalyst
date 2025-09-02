package router

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/martian/v3/cors"

	"github.com/SecurityBrewery/catalyst/app/auth"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/mail"
	"github.com/SecurityBrewery/catalyst/app/service"
	"github.com/SecurityBrewery/catalyst/app/upload"
)

func New(service *service.Service, queries *sqlc.Queries, uploader *upload.Uploader, mailer *mail.Mailer) (*chi.Mux, error) {
	r := chi.NewRouter()

	// middleware for the router
	r.Use(func(next http.Handler) http.Handler {
		return http.Handler(cors.NewHandler(next))
	})
	r.Use(demoMode(queries))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(time.Second * 60))
	r.Use(middleware.Recoverer)

	// base routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui/", http.StatusFound)
	})
	r.Get("/ui/*", staticFiles)
	r.Get("/health", healthHandler(queries))

	// auth routes
	r.Mount("/auth", auth.Server(queries, mailer))

	// API routes
	r.With(auth.Middleware(queries)).Mount("/api", http.StripPrefix("/api", service))

	uploadHandler, err := tusRoutes(queries, uploader)
	if err != nil {
		return nil, err
	}

	r.Mount("/files", http.StripPrefix("/files", uploadHandler))

	return r, nil
}

func healthHandler(queries *sqlc.Queries) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := queries.ListFeatures(r.Context(), sqlc.ListFeaturesParams{Offset: 0, Limit: 100}); err != nil {
			slog.ErrorContext(r.Context(), "Failed to get flags", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte("OK"))
	}
}
