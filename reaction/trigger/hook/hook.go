package hook

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"go.uber.org/multierr"

	"github.com/SecurityBrewery/catalyst/migrations"
	"github.com/SecurityBrewery/catalyst/reaction/action"
	"github.com/SecurityBrewery/catalyst/webhook"
)

type Hook struct {
	Collections []string `json:"collections"`
	Events      []string `json:"events"`
}

func BindHooks(app core.App, test bool) {
	app.OnRecordCreateRequest().BindFunc(func(e *core.RecordRequestEvent) error {
		return hook(e, app, "create", e.Collection.Name, e.Record, test)
	})
	app.OnRecordUpdateRequest().BindFunc(func(e *core.RecordRequestEvent) error {
		return hook(e, app, "update", e.Collection.Name, e.Record, test)
	})
	app.OnRecordDeleteRequest().BindFunc(func(e *core.RecordRequestEvent) error {
		return hook(e, app, "delete", e.Collection.Name, e.Record, test)
	})
}

func hook(req *core.RecordRequestEvent, app core.App, event, collection string, record *core.Record, test bool) error {
	if err := req.Next(); err != nil {
		return err
	}

	if !test {
		go mustRunHook(app, collection, event, record, req.Auth)
	} else {
		mustRunHook(app, collection, event, record, req.Auth)
	}

	return nil
}

func mustRunHook(app core.App, collection, event string, record, auth *core.Record) {
	ctx := context.Background()

	if err := runHook(ctx, app, collection, event, record, auth); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("failed to run hook reaction: %v", err))
	}
}

func runHook(ctx context.Context, app core.App, collection, event string, record, auth *core.Record) error {
	payload, err := json.Marshal(&webhook.Payload{
		Action:     event,
		Collection: collection,
		Record:     record,
		Auth:       auth,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	hooks, err := findByHookTrigger(app, collection, event)
	if err != nil {
		return fmt.Errorf("failed to find hook by trigger: %w", err)
	}

	if len(hooks) == 0 {
		return nil
	}

	var errs error

	for _, hook := range hooks {
		_, err = action.Run(ctx, app, hook.GetString("action"), hook.GetString("actiondata"), string(payload))
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to run hook reaction: %w", err))
		}
	}

	return errs
}

func findByHookTrigger(app core.App, collection, event string) ([]*core.Record, error) {
	records, err := app.FindAllRecords(migrations.ReactionCollectionName, dbx.HashExp{"trigger": "hook"})
	if err != nil {
		return nil, fmt.Errorf("failed to find hook reaction: %w", err)
	}

	if len(records) == 0 {
		return nil, nil
	}

	var matchedRecords []*core.Record

	for _, record := range records {
		var hook Hook
		if err := json.Unmarshal([]byte(record.GetString("triggerdata")), &hook); err != nil {
			return nil, err
		}

		if slices.Contains(hook.Collections, collection) && slices.Contains(hook.Events, event) {
			matchedRecords = append(matchedRecords, record)
		}
	}

	return matchedRecords, nil
}
