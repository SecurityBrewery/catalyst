package webhook

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"

	"github.com/SecurityBrewery/catalyst/migrations"
	"github.com/SecurityBrewery/catalyst/reaction/action"
	"github.com/SecurityBrewery/catalyst/reaction/action/webhook"
)

type Webhook struct {
	Token string `json:"token"`
	Path  string `json:"path"`
}

const prefix = "/reaction/"

func BindHooks(app core.App) {
	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		e.Router.Any(prefix, handle(e.App))

		return e.Next()
	})
}

func handle(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		record, payload, apiErr := parseRequest(app, e.Request)
		if apiErr != nil {
			return apiErr
		}

		output, err := action.Run(e.Request.Context(), app, record.GetString("action"), record.GetString("actiondata"), string(payload))
		if err != nil {
			return apis.NewApiError(http.StatusInternalServerError, err.Error(), nil)
		}

		return writeOutput(e, output)
	}
}

func parseRequest(dao core.App, r *http.Request) (*core.Record, []byte, *router.ApiError) {
	if !strings.HasPrefix(r.URL.Path, prefix) {
		return nil, nil, apis.NewApiError(http.StatusNotFound, "wrong prefix", nil)
	}

	reactionName := strings.TrimPrefix(r.URL.Path, prefix)

	record, trigger, found, err := findByWebhookTrigger(dao, reactionName)
	if err != nil {
		return nil, nil, apis.NewNotFoundError(err.Error(), nil)
	}

	if !found {
		return nil, nil, apis.NewNotFoundError("reaction not found", nil)
	}

	if trigger.Token != "" {
		auth := r.Header.Get("Authorization")

		if !strings.HasPrefix(auth, "Bearer ") {
			return nil, nil, apis.NewUnauthorizedError("missing token", nil)
		}

		if trigger.Token != strings.TrimPrefix(auth, "Bearer ") {
			return nil, nil, apis.NewUnauthorizedError("invalid token", nil)
		}
	}

	body, isBase64Encoded := webhook.EncodeBody(r.Body)

	payload, err := json.Marshal(&Request{
		Method:          r.Method,
		Path:            r.URL.EscapedPath(),
		Headers:         r.Header,
		Query:           r.URL.Query(),
		Body:            body,
		IsBase64Encoded: isBase64Encoded,
	})
	if err != nil {
		return nil, nil, apis.NewApiError(http.StatusInternalServerError, err.Error(), nil)
	}

	return record, payload, nil
}

func findByWebhookTrigger(app core.App, path string) (*core.Record, *Webhook, bool, error) {
	records, err := app.FindAllRecords(migrations.ReactionCollectionName, dbx.HashExp{"trigger": "webhook"})
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

func writeOutput(e *core.RequestEvent, output []byte) error {
	var catalystResponse webhook.Response
	if err := json.Unmarshal(output, &catalystResponse); err == nil && catalystResponse.StatusCode != 0 {
		for key, values := range catalystResponse.Headers {
			for _, value := range values {
				e.Response.Header().Add(key, value)
			}
		}

		if catalystResponse.IsBase64Encoded {
			output, err = base64.StdEncoding.DecodeString(catalystResponse.Body)
			if err != nil {
				return fmt.Errorf("error decoding base64 body: %w", err)
			}
		} else {
			output = []byte(catalystResponse.Body)
		}
	}

	if IsJSON(output) {
		return e.JSON(http.StatusOK, json.RawMessage(output))
	}

	return e.String(http.StatusOK, string(output))
}
