package service

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/SecurityBrewery/catalyst/app/auth"
	"github.com/SecurityBrewery/catalyst/app/openapi"
)

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	middlewareFuncs := []openapi.StrictMiddlewareFunc{auth.ValidateScopesStrict, auth.LogError}
	apiHandler := openapi.Handler(openapi.NewStrictHandlerWithOptions(s, middlewareFuncs, openapi.StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  jsonError,
		ResponseErrorHandlerFunc: jsonError,
	}))

	apiHandler.ServeHTTP(w, r)
}

func jsonError(w http.ResponseWriter, r *http.Request, err error) {
	b, err := json.Marshal(openapi.Error{
		Status:  http.StatusInternalServerError,
		Error:   "An internal error occurred",
		Message: err.Error(),
	})
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to marshal error response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write(b)
}
