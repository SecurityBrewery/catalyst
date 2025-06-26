package hook

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"

	"go.uber.org/multierr"

	"github.com/SecurityBrewery/catalyst/app"
	"github.com/SecurityBrewery/catalyst/app/auth/usercontext"
	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/permission"
	"github.com/SecurityBrewery/catalyst/app/reaction/action"
	"github.com/SecurityBrewery/catalyst/app/webhook"
)

type Hook struct {
	Collections []string `json:"collections"`
	Events      []string `json:"events"`
}

func BindHooks(app *app.App, test bool) {
	app.Hooks.OnRecordAfterCreateRequest.Subscribe(func(ctx context.Context, table string, record any) {
		hook(ctx, app, permission.CreateAction, table, record, test)
	})
	app.Hooks.OnRecordAfterUpdateRequest.Subscribe(func(ctx context.Context, table string, record any) {
		hook(ctx, app, permission.UpdateAction, table, record, test)
	})
	app.Hooks.OnRecordAfterDeleteRequest.Subscribe(func(ctx context.Context, table string, record any) {
		hook(ctx, app, permission.DeleteAction, table, record, test)
	})
}

func hook(ctx context.Context, app *app.App, event, collection string, record any, test bool) {
	user, ok := usercontext.UserFromContext(ctx)
	if !ok {
		slog.ErrorContext(ctx, "failed to get user from session")

		return
	}

	if !test {
		go mustRunHook(context.Background(), app, collection, event, record, user) //nolint:contextcheck
	} else {
		mustRunHook(ctx, app, collection, event, record, user)
	}
}

func mustRunHook(ctx context.Context, app *app.App, collection, event string, record any, auth *sqlc.User) {
	if err := runHook(ctx, app, collection, event, record, auth); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("failed to run hook reaction: %v", err))
	}
}

func runHook(ctx context.Context, app *app.App, collection, event string, record any, auth *sqlc.User) error {
	payload, err := json.Marshal(&webhook.Payload{
		Action:     event,
		Collection: collection,
		Record:     record,
		Auth:       auth,
		Admin:      nil,
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

	settings, err := database.LoadSettings(ctx, app.Queries)
	if err != nil {
		return fmt.Errorf("failed to load settings: %w", err)
	}

	for _, hook := range hooks {
		_, err = action.Run(ctx, settings.Meta.AppURL, app.Auth, app.Queries, hook.Action, hook.Actiondata, payload)
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
		if err := json.Unmarshal(reaction.Triggerdata, &hook); err != nil {
			return nil, err
		}

		if slices.Contains(hook.Collections, collection) && slices.Contains(hook.Events, event) {
			matchedRecords = append(matchedRecords, &reaction)
		}
	}

	return matchedRecords, nil
}
