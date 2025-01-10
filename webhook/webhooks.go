package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/migrations"
)

const webhooksCollectionName = "webhooks"

type Webhook struct {
	ID          string `db:"id"          json:"id"`
	Name        string `db:"name"        json:"name"`
	Collection  string `db:"collection"  json:"collection"`
	Destination string `db:"destination" json:"destination"`
}

func BindHooks(app core.App) {
	migrations.Register(func(app core.App) error {
		webhooksCollection := core.NewBaseCollection(webhooksCollectionName)
		webhooksCollection.System = true
		webhooksCollection.Fields.Add(&core.TextField{Name: "name", Required: true})
		webhooksCollection.Fields.Add(&core.TextField{Name: "collection", Required: true})
		webhooksCollection.Fields.Add(&core.URLField{Name: "destination", Required: true})

		return app.Save(webhooksCollection)
	}, nil, "1690000000_webhooks.go")

	app.OnRecordCreateRequest().BindFunc(func(e *core.RecordRequestEvent) error {
		if err := e.Next(); err != nil {
			return err
		}

		return event(app, "create", e.Collection.Name, e.Record, e)
	})
	app.OnRecordUpdateRequest().BindFunc(func(e *core.RecordRequestEvent) error {
		if err := e.Next(); err != nil {
			return err
		}

		return event(app, "update", e.Collection.Name, e.Record, e)
	})
	app.OnRecordDeleteRequest().BindFunc(func(e *core.RecordRequestEvent) error {
		if err := e.Next(); err != nil {
			return err
		}

		return event(app, "delete", e.Collection.Name, e.Record, e)
	})
}

type Payload struct {
	Action     string       `json:"action"`
	Collection string       `json:"collection"`
	Record     *core.Record `json:"record"`
	Auth       *core.Record `json:"auth,omitempty"`
}

func event(app core.App, event, collection string, record *core.Record, req *core.RecordRequestEvent) error {
	var webhooks []Webhook
	if err := app.DB().
		Select().
		From(webhooksCollectionName).
		Where(dbx.HashExp{"collection": collection}).
		All(&webhooks); err != nil {
		return err
	}

	if len(webhooks) == 0 {
		return nil
	}

	payload, err := json.Marshal(&Payload{
		Action:     event,
		Collection: collection,
		Record:     record,
		Auth:       req.Auth,
	})
	if err != nil {
		return err
	}

	for _, webhook := range webhooks {
		if err := sendWebhook(req.Request.Context(), webhook, payload); err != nil {
			app.Logger().Error("failed to send webhook", "action", event, "name", webhook.Name, "collection", webhook.Collection, "destination", webhook.Destination, "error", err.Error())
		} else {
			app.Logger().Info("webhook sent", "action", event, "name", webhook.Name, "collection", webhook.Collection, "destination", webhook.Destination)
		}
	}

	return nil
}

func sendWebhook(ctx context.Context, webhook Webhook, payload []byte) error {
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
