package hook

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
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

func BindHooks(app core.App) {
	app.OnRecordAfterCreateRequest().Add(func(e *core.RecordCreateEvent) error {
		if err := hook(app.Dao(), "create", e.Collection.Name, e.Record, e.HttpContext); err != nil {
			app.Logger().Error("failed to find hook reaction", "error", err.Error())
		}

		return nil
	})
	app.OnRecordAfterUpdateRequest().Add(func(e *core.RecordUpdateEvent) error {
		if err := hook(app.Dao(), "update", e.Collection.Name, e.Record, e.HttpContext); err != nil {
			app.Logger().Error("failed to find hook reaction", "error", err.Error())
		}

		return nil
	})
	app.OnRecordAfterDeleteRequest().Add(func(e *core.RecordDeleteEvent) error {
		if err := hook(app.Dao(), "delete", e.Collection.Name, e.Record, e.HttpContext); err != nil {
			app.Logger().Error("failed to find hook reaction", "error", err.Error())
		}

		return nil
	})
}

func hook(dao *daos.Dao, event, collection string, record *models.Record, ctx echo.Context) error {
	auth, _ := ctx.Get(apis.ContextAuthRecordKey).(*models.Record)
	admin, _ := ctx.Get(apis.ContextAdminKey).(*models.Admin)

	hook, found, err := findByHookTrigger(dao, collection, event)
	if err != nil {
		return fmt.Errorf("failed to find hook reaction: %w", err)
	}

	if !found {
		return nil
	}

	payload, err := json.Marshal(&webhook.Payload{
		Action:     event,
		Collection: collection,
		Record:     record,
		Auth:       auth,
		Admin:      admin,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	_, err = action.Run(ctx.Request().Context(), hook.GetString("action"), hook.GetString("actiondata"), string(payload))
	if err != nil {
		return fmt.Errorf("failed to run hook reaction: %w", err)
	}

	return nil
}

func findByHookTrigger(dao *daos.Dao, collection, event string) (*models.Record, bool, error) {
	records, err := dao.FindRecordsByExpr(migrations.ReactionCollectionName, dbx.HashExp{"trigger": "hook"})
	if err != nil {
		return nil, false, err
	}

	if len(records) == 0 {
		return nil, false, nil
	}

	for _, record := range records {
		var hook Hook
		if err := json.Unmarshal([]byte(record.GetString("triggerdata")), &hook); err != nil {
			return nil, false, err
		}

		if slices.Contains(hook.Collections, collection) && slices.Contains(hook.Events, event) {
			return record, true, nil
		}
	}

	return nil, false, nil
}
