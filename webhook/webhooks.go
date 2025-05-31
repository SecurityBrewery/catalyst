package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/SecurityBrewery/catalyst/app2"
	"github.com/SecurityBrewery/catalyst/app2/database/sqlc"
)

const webhooksCollection = "webhooks"

type Webhook struct {
	ID          string `db:"id"          json:"id"`
	Name        string `db:"name"        json:"name"`
	Collection  string `db:"collection"  json:"collection"`
	Destination string `db:"destination" json:"destination"`
}

func BindHooks(app app2.App2) {
	/*app.OnRecordAfterCreateRequest().Add(func(e *core.RecordCreateEvent) error {
		return event(app, "create", e.Collection.Name, e.Record, e.HttpContext)
	})
	app.OnRecordAfterUpdateRequest().Add(func(e *core.RecordUpdateEvent) error {
		return event(app, "update", e.Collection.Name, e.Record, e.HttpContext)
	})
	app.OnRecordAfterDeleteRequest().Add(func(e *core.RecordDeleteEvent) error {
		return event(app, "delete", e.Collection.Name, e.Record, e.HttpContext)
	})*/
}

type Payload struct {
	Action     string     `json:"action"`
	Collection string     `json:"collection"`
	Record     any        `json:"record"`
	Auth       *sqlc.User `json:"auth,omitempty"`
	Admin      *sqlc.User `json:"admin,omitempty"`
}

func event(app app2.App2, event, collection string, record any, ctx context.Context) error {
	// auth, _ := ctx.Get(apis.ContextAuthRecordKey).(*models.Record) // TODO
	// admin, _ := ctx.Get(apis.ContextAdminKey).(*models.Admin)     // TODO
	auth, admin := &sqlc.User{}, &sqlc.User{}

	webhooks, err := app.Queries.ListWebhooks(ctx, sqlc.ListWebhooksParams{
		Offset: 0,
		Limit:  100,
	})
	if err != nil {
		return err
	}

	if len(webhooks) == 0 {
		return nil
	}

	payload, err := json.Marshal(&Payload{
		Action:     event,
		Collection: collection,
		Record:     record,
		Auth:       auth,
		Admin:      admin,
	})
	if err != nil {
		return err
	}

	for _, webhook := range webhooks {
		if err := sendWebhook(ctx, webhook, payload); err != nil {
			slog.ErrorContext(ctx, "failed to send webhook", "action", event, "name", webhook.Name, "collection", webhook.Collection, "destination", webhook.Destination, "error", err.Error())
		} else {
			slog.InfoContext(ctx, "webhook sent", "action", event, "name", webhook.Name, "collection", webhook.Collection, "destination", webhook.Destination)
		}
	}

	return nil
}

func sendWebhook(ctx context.Context, webhook sqlc.ListWebhooksRow, payload []byte) error {
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, webhook.Destination, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)

		return fmt.Errorf("failed to send webhook: %s", string(b))
	}

	return nil
}
