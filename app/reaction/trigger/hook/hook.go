package hook

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"slices"

	"github.com/SecurityBrewery/catalyst/app/auth/usercontext"
	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/hook"
	"github.com/SecurityBrewery/catalyst/app/reaction/action"
	"github.com/SecurityBrewery/catalyst/app/webhook"
)

type Hook struct {
	Collections []string `json:"collections"`
	Events      []string `json:"events"`
}

func BindHooks(hooks *hook.Hooks, queries *sqlc.Queries, test bool) {
	hooks.OnRecordAfterCreateRequest.Subscribe(func(ctx context.Context, table string, record any) {
		bindHook(ctx, queries, database.CreateAction, table, record, test)
	})
	hooks.OnRecordAfterUpdateRequest.Subscribe(func(ctx context.Context, table string, record any) {
		bindHook(ctx, queries, database.UpdateAction, table, record, test)
	})
	hooks.OnRecordAfterDeleteRequest.Subscribe(func(ctx context.Context, table string, record any) {
		bindHook(ctx, queries, database.DeleteAction, table, record, test)
	})
}

func bindHook(ctx context.Context, queries *sqlc.Queries, event, collection string, record any, test bool) {
	user, ok := usercontext.UserFromContext(ctx)
	if !ok {
		slog.ErrorContext(ctx, "failed to get user from session")

		return
	}

	if !test {
		go mustRunHook(context.Background(), queries, collection, event, record, user) //nolint:contextcheck
	} else {
		mustRunHook(ctx, queries, collection, event, record, user)
	}
}

func mustRunHook(ctx context.Context, queries *sqlc.Queries, collection, event string, record any, auth *sqlc.User) {
	if err := runHook(ctx, queries, collection, event, record, auth); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("failed to run hook reaction: %v", err))
	}
}

func runHook(ctx context.Context, queries *sqlc.Queries, collection, event string, record any, auth *sqlc.User) error {
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

	hooks, err := findByHookTrigger(ctx, queries, collection, event)
	if err != nil {
		return fmt.Errorf("failed to find hook by trigger: %w", err)
	}

	if len(hooks) == 0 {
		return nil
	}

	settings, err := database.LoadSettings(ctx, queries)
	if err != nil {
		return fmt.Errorf("failed to load settings: %w", err)
	}

	var errs []error

	for _, hook := range hooks {
		_, err = action.Run(ctx, settings.Meta.AppURL, queries, hook.Action, hook.Actiondata, payload)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to run hook reaction: %w", err))
		}
	}

	return errors.Join(errs...)
}

func findByHookTrigger(ctx context.Context, queries *sqlc.Queries, collection, event string) ([]*sqlc.ListReactionsByTriggerRow, error) {
	reactions, err := database.PaginateItems(ctx, func(ctx context.Context, offset, limit int64) ([]sqlc.ListReactionsByTriggerRow, error) {
		return queries.ListReactionsByTrigger(ctx, sqlc.ListReactionsByTriggerParams{Trigger: "hook", Limit: limit, Offset: offset})
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find hook reaction: %w", err)
	}

	if len(reactions) == 0 {
		return nil, nil
	}

	var matchedRecords []*sqlc.ListReactionsByTriggerRow

	for _, reaction := range reactions {
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
