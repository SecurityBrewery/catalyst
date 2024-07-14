package reaction

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

const prefix = "/reaction/"

func Handle(app core.App) func(c echo.Context) error {
	return func(c echo.Context) error {
		r := c.Request()
		w := c.Response()

		if !strings.HasPrefix(r.URL.Path, prefix) {
			return apis.NewApiError(http.StatusNotFound, "wrong prefix", nil)
		}

		reactionName := strings.TrimPrefix(r.URL.Path, prefix)

		reaction, trigger, found, err := findReaction(app, reactionName)
		if err != nil {
			return apis.NewNotFoundError(err.Error(), nil)
		}

		if !found {
			return apis.NewNotFoundError("reaction not found", nil)
		}

		if trigger.Token != "" && bearerToken(r) != trigger.Token {
			return apis.NewUnauthorizedError("invalid token", nil)
		}

		payload, err := requestToPayload(r)
		if err != nil {
			return apis.NewApiError(http.StatusInternalServerError, err.Error(), nil)
		}

		output, err := runReaction(r.Context(), reaction, payload)
		if err != nil {
			return apis.NewApiError(http.StatusInternalServerError, err.Error(), nil)
		}

		outputToResponse(app.Logger(), w, output)

		return nil
	}
}

func bearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")

	if !strings.HasPrefix(auth, "Bearer ") {
		return ""
	}

	return strings.TrimPrefix(auth, "Bearer ")
}
