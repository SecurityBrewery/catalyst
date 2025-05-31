package hook

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"

	"go.uber.org/multierr"

	"github.com/SecurityBrewery/catalyst/app2"
	"github.com/SecurityBrewery/catalyst/app2/database/sqlc"
	"github.com/SecurityBrewery/catalyst/reaction/action"
	"github.com/SecurityBrewery/catalyst/webhook"
)

type Hook struct {
	Collections []string `json:"collections"`
	Events      []string `json:"events"`
}

func BindHooks(app *app2.App2, test bool) {
	/* app.OnRecordAfterCreateRequest().Add(func(e *core.RecordCreateEvent) error { // TODO
		return hook(e.HttpContext, app, "create", e.Collection.Name, e.Record, test)
	}) */
	/* app.OnRecordAfterUpdateRequest().Add(func(e *core.RecordUpdateEvent) error { // TODO
		return hook(e.HttpContext, app, "update", e.Collection.Name, e.Record, test)
	}) */
	/* app.OnRecordAfterDeleteRequest().Add(func(e *core.RecordDeleteEvent) error { // TODO
		return hook(e.HttpContext, app, "delete", e.Collection.Name, e.Record, test)
	}) */
}

func hook(ctx context.Context, app *app2.App2, event, collection string, record any, test bool) error {
	// auth, _ := ctx.Get(apis.ContextAuthRecordKey).(*models.Record) // TODO
	// admin, _ := ctx.Get(apis.ContextAdminKey).(*models.Admin)
	auth, admin := &sqlc.User{}, &sqlc.User{}

	if !test {
		go mustRunHook(app, collection, event, record, auth, admin)
	} else {
		mustRunHook(app, collection, event, record, auth, admin)
	}

	return nil
}

func mustRunHook(app *app2.App2, collection, event string, record any, auth *sqlc.User, admin *sqlc.User) {
	ctx := context.Background()

	if err := runHook(ctx, app, collection, event, record, auth, admin); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("failed to run hook reaction: %v", err))
	}
}

func runHook(ctx context.Context, app *app2.App2, collection, event string, record any, auth *sqlc.User, admin *sqlc.User) error {
	payload, err := json.Marshal(&webhook.Payload{
		Action:     event,
		Collection: collection,
		Record:     record,
		Auth:       auth,
		Admin:      admin,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	hooks, err := findByHookTrigger(ctx, app.Queries, collection, event)
	if err != nil {
		return fmt.Errorf("failed to find hook by trigger: %w", err)
	}

	if len(hooks) == 0 {
		return nil
	}

	var errs error

	for _, hook := range hooks {
		_, err = action.Run(ctx, app, hook.Action, hook.Actiondata, string(payload))
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to run hook reaction: %w", err))
		}
	}

	return errs
}

func findByHookTrigger(ctx context.Context, queries *sqlc.Queries, collection, event string) ([]*sqlc.ListReactionsRow, error) {
	reactions, err := queries.ListReactions(ctx, sqlc.ListReactionsParams{
		Offset: 0,
		Limit:  100,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find hook reaction: %w", err)
	}

	if len(reactions) == 0 {
		return nil, nil
	}

	var matchedRecords []*sqlc.ListReactionsRow

	for _, reaction := range reactions {
		if reaction.Trigger != "hook" {
			continue
		}

		var hook Hook
		if err := json.Unmarshal([]byte(reaction.Triggerdata), &hook); err != nil {
			return nil, err
		}

		if slices.Contains(hook.Collections, collection) && slices.Contains(hook.Events, event) {
			matchedRecords = append(matchedRecords, &reaction)
		}
	}

	return matchedRecords, nil
}
