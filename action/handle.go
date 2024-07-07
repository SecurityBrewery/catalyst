package action

import (
	"github.com/pocketbase/pocketbase/core"
	"net/http"
	"strings"
)

const prefix = "/action/"

func Handle(app core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, prefix) {
			errResponse(app.Logger(), w, http.StatusNotFound, "wrong prefix")

			return
		}

		actionName := strings.TrimPrefix(r.URL.Path, prefix)

		action, found, err := findAction(app, actionName)
		if err != nil {
			errResponse(app.Logger(), w, http.StatusInternalServerError, err.Error())

			return
		}

		if !found {
			errResponse(app.Logger(), w, http.StatusNotFound, "action not found")

			return
		}

		token := action.GetString("token")

		if token != "" && bearerToken(r) != token {
			errResponse(app.Logger(), w, http.StatusUnauthorized, "invalid token")

			return
		}

		payload, err := requestToPayload(r)
		if err != nil {
			errResponse(app.Logger(), w, http.StatusInternalServerError, err.Error())

			return
		}

		output, err := runAction(
			action.GetString("type"),
			action.GetString("name"),
			action.GetString("bootstrap"),
			action.GetString("script"),
			payload,
		)
		if err != nil {
			errResponse(app.Logger(), w, http.StatusInternalServerError, err.Error())

			return
		}

		outputToResponse(app.Logger(), w, output)
	}
}

func bearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")

	if !strings.HasPrefix(auth, "Bearer ") {
		return ""
	}

	return strings.TrimPrefix(auth, "Bearer ")
}
