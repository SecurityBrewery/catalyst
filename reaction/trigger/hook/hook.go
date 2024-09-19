package hook

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/multierr"
	"log/slog"
	"slices"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"

	"github.com/SecurityBrewery/catalyst/migrations"
	"github.com/SecurityBrewery/catalyst/reaction/action"
	"github.com/SecurityBrewery/catalyst/webhook"
)

type Hook struct {
	Collections []string `json:"collections"`
	Events      []string `json:"events"`
}

func BindHooks(pb *pocketbase.PocketBase, test bool) {
	pb.App.OnRecordAfterCreateRequest().Add(func(e *core.RecordCreateEvent) error {
		return hook(e.HttpContext, pb.App, "create", e.Collection.Name, e.Record, test)
	})
	pb.App.OnRecordAfterUpdateRequest().Add(func(e *core.RecordUpdateEvent) error {
		return hook(e.HttpContext, pb.App, "update", e.Collection.Name, e.Record, test)
	})
	pb.App.OnRecordAfterDeleteRequest().Add(func(e *core.RecordDeleteEvent) error {
		return hook(e.HttpContext, pb.App, "delete", e.Collection.Name, e.Record, test)
	})
}

func hook(ctx echo.Context, app core.App, event, collection string, record *models.Record, test bool) error {
	auth, _ := ctx.Get(apis.ContextAuthRecordKey).(*models.Record)
	admin, _ := ctx.Get(apis.ContextAdminKey).(*models.Admin)

	if !test {
		go mustRunHook(app, collection, event, record, auth, admin)
	} else {
		mustRunHook(app, collection, event, record, auth, admin)
	}

	return nil
}

func mustRunHook(app core.App, collection, event string, record, auth *models.Record, admin *models.Admin) {
	ctx := context.Background()

	if err := runHook(ctx, app, collection, event, record, auth, admin); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("failed to run hook reaction: %v", err))
	}
}

func runHook(ctx context.Context, app core.App, collection, event string, record, auth *models.Record, admin *models.Admin) error {
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

	hooks, err := findByHookTrigger(app.Dao(), collection, event)
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

func findByHookTrigger(dao *daos.Dao, collection, event string) ([]*models.Record, error) {
	records, err := dao.FindRecordsByExpr(migrations.ReactionCollectionName, dbx.HashExp{"trigger": "hook"})
	if err != nil {
		return nil, fmt.Errorf("failed to find hook reaction: %w", err)
	}

	if len(records) == 0 {
		return nil, nil
	}

	var matchedRecords []*models.Record

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
