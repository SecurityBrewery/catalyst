package hook

import (
	"encoding/json"
	"slices"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
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
	// Bind hooks
	app.OnRecordAfterCreateRequest().Add(func(e *core.RecordCreateEvent) error {
		return hook(app, "create", e.Collection.Name, e.Record, e.HttpContext)
	})
	app.OnRecordAfterUpdateRequest().Add(func(e *core.RecordUpdateEvent) error {
		return hook(app, "update", e.Collection.Name, e.Record, e.HttpContext)
	})
	app.OnRecordAfterDeleteRequest().Add(func(e *core.RecordDeleteEvent) error {
		return hook(app, "delete", e.Collection.Name, e.Record, e.HttpContext)
	})
}

func hook(app core.App, event, collection string, record *models.Record, ctx echo.Context) error {
	auth, _ := ctx.Get(apis.ContextAuthRecordKey).(*models.Record)
	admin, _ := ctx.Get(apis.ContextAdminKey).(*models.Admin)

	hook, found, err := findByHookTrigger(app, collection, event)
	if err != nil {
		app.Logger().Error("failed to find hook reaction", "error", err.Error())

		return nil
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
		app.Logger().Error("failed to marshal payload", "error", err.Error())

		return nil
	}

	_, err = action.Run(ctx.Request().Context(), hook.GetString("action"), hook.GetString("actiondata"), string(payload))
	if err != nil {
		app.Logger().Error("failed to run hook reaction", "error", err.Error())

		return nil
	}

	return nil
}

func findByHookTrigger(app core.App, collection, event string) (*models.Record, bool, error) {
	records, err := app.Dao().FindRecordsByExpr(migrations.ReactionCollectionName, dbx.HashExp{"trigger": "hook"})
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
