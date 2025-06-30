package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/SecurityBrewery/catalyst/app/auth/usercontext"
	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/hook"
)

type Webhook struct {
	ID          string `db:"id"          json:"id"`
	Name        string `db:"name"        json:"name"`
	Collection  string `db:"collection"  json:"collection"`
	Destination string `db:"destination" json:"destination"`
}

func BindHooks(hooks *hook.Hooks, queries *sqlc.Queries) {
	hooks.OnRecordAfterCreateRequest.Subscribe(func(ctx context.Context, table string, record any) {
		event(ctx, queries, database.CreateAction, table, record)
	})
	hooks.OnRecordAfterUpdateRequest.Subscribe(func(ctx context.Context, table string, record any) {
		event(ctx, queries, database.UpdateAction, table, record)
	})
	hooks.OnRecordAfterDeleteRequest.Subscribe(func(ctx context.Context, table string, record any) {
		event(ctx, queries, database.DeleteAction, table, record)
	})
}

type Payload struct {
	Action     string     `json:"action"`
	Collection string     `json:"collection"`
	Record     any        `json:"record"`
	Auth       *sqlc.User `json:"auth,omitempty"`
	Admin      *sqlc.User `json:"admin,omitempty"`
}

func event(ctx context.Context, queries *sqlc.Queries, event, collection string, record any) {
	user, ok := usercontext.UserFromContext(ctx)
	if !ok {
		slog.ErrorContext(ctx, "failed to get auth session")

		return
	}

	webhooks, err := database.PaginateItems(ctx, func(ctx context.Context, offset, limit int64) ([]sqlc.ListWebhooksRow, error) {
		return queries.ListWebhooks(ctx, sqlc.ListWebhooksParams{Limit: limit, Offset: offset})
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to list webhooks", "error", err.Error())

		return
	}

	if len(webhooks) == 0 {
		return
	}

	payload, err := json.Marshal(&Payload{
		Action:     event,
		Collection: collection,
		Record:     record,
		Auth:       user,
		Admin:      nil,
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to marshal payload", "error", err.Error())

		return
	}

	for _, webhook := range webhooks {
		if err := sendWebhook(ctx, webhook, payload); err != nil {
			slog.ErrorContext(ctx, "failed to send webhook", "action", event, "name", webhook.Name, "collection", webhook.Collection, "destination", webhook.Destination, "error", err.Error())
		} else {
			slog.InfoContext(ctx, "webhook sent", "action", event, "name", webhook.Name, "collection", webhook.Collection, "destination", webhook.Destination)
		}
	}
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
