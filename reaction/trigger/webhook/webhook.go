package webhook

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/SecurityBrewery/catalyst/app2"
	"github.com/SecurityBrewery/catalyst/app2/database/sqlc"
	"github.com/SecurityBrewery/catalyst/reaction/action"
	"github.com/SecurityBrewery/catalyst/reaction/action/webhook"
)

type Webhook struct {
	Token string `json:"token"`
	Path  string `json:"path"`
}

const prefix = "/reaction/"

func BindHooks(app *app2.App2) {
	app.Router.HandleFunc(prefix+"*", handle(app))
}

func handle(app *app2.App2) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reaction, payload, status, err := parseRequest(app.Queries, r)
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}

		output, err := action.Run(r.Context(), app, reaction.Action, reaction.Actiondata, string(payload))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		writeOutput(w, output)
	}
}

func parseRequest(queries *sqlc.Queries, r *http.Request) (*sqlc.ListReactionsRow, []byte, int, error) {
	if !strings.HasPrefix(r.URL.Path, prefix) {
		return nil, nil, http.StatusNotFound, fmt.Errorf("wrong prefix")
	}

	reactionName := strings.TrimPrefix(r.URL.Path, prefix)

	reaction, trigger, found, err := findByWebhookTrigger(r.Context(), queries, reactionName)
	if err != nil {
		return nil, nil, http.StatusNotFound, err
	}

	if !found {
		return nil, nil, http.StatusNotFound, fmt.Errorf("reaction not found")
	}

	if trigger.Token != "" {
		auth := r.Header.Get("Authorization")

		if !strings.HasPrefix(auth, "Bearer ") {
			return nil, nil, http.StatusUnauthorized, fmt.Errorf("missing token")
		}

		if trigger.Token != strings.TrimPrefix(auth, "Bearer ") {
			return nil, nil, http.StatusUnauthorized, fmt.Errorf("invalid token")
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
		return nil, nil, http.StatusInternalServerError, err
	}

	return reaction, payload, http.StatusOK, nil
}

func findByWebhookTrigger(ctx context.Context, queries *sqlc.Queries, path string) (*sqlc.ListReactionsRow, *Webhook, bool, error) {
	reactions, err := queries.ListReactions(ctx, sqlc.ListReactionsParams{
		Offset: 0,
		Limit:  100,
	})
	if err != nil {
		return nil, nil, false, err
	}

	if len(reactions) == 0 {
		return nil, nil, false, nil
	}

	for _, reaction := range reactions {
		if reaction.Trigger != "webhook" {
			continue
		}

		var webhook Webhook
		if err := json.Unmarshal([]byte(reaction.Triggerdata), &webhook); err != nil {
			return nil, nil, false, err
		}

		if webhook.Path == path {
			return &reaction, &webhook, true, nil
		}
	}

	return nil, nil, false, nil
}

func writeOutput(w http.ResponseWriter, output []byte) error {
	var catalystResponse webhook.Response
	if err := json.Unmarshal(output, &catalystResponse); err == nil && catalystResponse.StatusCode != 0 {
		for key, values := range catalystResponse.Headers {
			for _, value := range values {
				w.Header().Add(key, value)
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

	if json.Valid(output) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(output)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(output)
	}

	return nil
}
