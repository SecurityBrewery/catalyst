package webhook

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"

	"github.com/SecurityBrewery/catalyst/migrations"
	"github.com/SecurityBrewery/catalyst/reaction/action"
	"github.com/SecurityBrewery/catalyst/reaction/action/webhook"
)

type Webhook struct {
	Token string `json:"token"`
	Path  string `json:"path"`
}

func BindHooks(app core.App) {
	// Bind webhooks
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.Any("/reaction/*", handle(e.App))

		return nil
	})
}

const prefix = "/reaction/"

func handle(app core.App) func(c echo.Context) error {
	return func(c echo.Context) error {
		r := c.Request()
		w := c.Response()

		if !strings.HasPrefix(r.URL.Path, prefix) {
			return apis.NewApiError(http.StatusNotFound, "wrong prefix", nil)
		}

		reactionName := strings.TrimPrefix(r.URL.Path, prefix)

		record, trigger, found, err := findByWebhookTrigger(app, reactionName)
		if err != nil {
			return apis.NewNotFoundError(err.Error(), nil)
		}

		if !found {
			return apis.NewNotFoundError("reaction not found", nil)
		}

		if trigger.Token != "" && bearerToken(r) != trigger.Token {
			return apis.NewUnauthorizedError("invalid token", nil)
		}

		payload, err := webhook.RequestFromHTTPRequest(r)
		if err != nil {
			return apis.NewApiError(http.StatusInternalServerError, err.Error(), nil)
		}

		output, err := action.Run(r.Context(), record.GetString("action"), record.GetString("actiondata"), payload)
		if err != nil {
			return apis.NewApiError(http.StatusInternalServerError, err.Error(), nil)
		}

		webhook.OutputToResponse(w, output)

		return nil
	}
}

func findByWebhookTrigger(app core.App, path string) (*models.Record, *Webhook, bool, error) {
	records, err := app.Dao().FindRecordsByExpr(migrations.ReactionCollectionName, dbx.HashExp{"trigger": "webhook"})
	if err != nil {
		return nil, nil, false, err
	}

	if len(records) == 0 {
		return nil, nil, false, nil
	}

	for _, record := range records {
		var webhook Webhook
		if err := json.Unmarshal([]byte(record.GetString("triggerdata")), &webhook); err != nil {
			return nil, nil, false, err
		}

		if webhook.Path == path {
			return record, &webhook, true, nil
		}
	}

	return nil, nil, false, nil
}

func bearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")

	if !strings.HasPrefix(auth, "Bearer ") {
		return ""
	}

	return strings.TrimPrefix(auth, "Bearer ")
}
