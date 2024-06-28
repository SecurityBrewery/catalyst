package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

const webhooksCollection = "webhooks"

type Webhook struct {
	ID          string `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Collection  string `db:"collection" json:"collection"`
	Destination string `db:"destination" json:"destination"`
}

func attachWebhooks(app *pocketbase.PocketBase) {
	migrations.Register(func(db dbx.Builder) error {
		return daos.New(db).SaveCollection(&models.Collection{
			Name:   webhooksCollection,
			Type:   models.CollectionTypeBase,
			System: true,
			Schema: schema.NewSchema(
				&schema.SchemaField{
					Name:     "name",
					Type:     schema.FieldTypeText,
					Required: true,
				},
				&schema.SchemaField{
					Name:     "collection",
					Type:     schema.FieldTypeText,
					Required: true,
				},
				&schema.SchemaField{
					Name:     "destination",
					Type:     schema.FieldTypeText,
					Required: true,
				},
			),
		})
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		id, err := dao.FindCollectionByNameOrId(webhooksCollection)
		if err != nil {
			return err
		}

		return dao.DeleteCollection(id)
	}, "1690000000_webhooks.go")

	app.OnRecordAfterCreateRequest().Add(func(e *core.RecordCreateEvent) error {
		return event(app, "create", e.Collection.Name, e.Record, e.HttpContext)
	})
	app.OnRecordAfterUpdateRequest().Add(func(e *core.RecordUpdateEvent) error {
		return event(app, "update", e.Collection.Name, e.Record, e.HttpContext)
	})
	app.OnRecordAfterDeleteRequest().Add(func(e *core.RecordDeleteEvent) error {
		return event(app, "delete", e.Collection.Name, e.Record, e.HttpContext)
	})
}

type Payload struct {
	Action     string         `json:"action"`
	Collection string         `json:"collection"`
	Record     *models.Record `json:"record"`
	Auth       *models.Record `json:"auth,omitempty"`
	Admin      *models.Admin  `json:"admin,omitempty"`
}

func event(app *pocketbase.PocketBase, action, collection string, record *models.Record, ctx echo.Context) error {
	auth, _ := ctx.Get(apis.ContextAuthRecordKey).(*models.Record)
	admin, _ := ctx.Get(apis.ContextAdminKey).(*models.Admin)

	var webhooks []Webhook
	if err := app.Dao().DB().
		Select().
		From(webhooksCollection).
		Where(dbx.HashExp{"collection": collection}).
		All(&webhooks); err != nil {
		return err
	}

	if len(webhooks) == 0 {
		return nil
	}

	b, err := json.Marshal(&Payload{
		Action:     action,
		Collection: collection,
		Record:     record,
		Auth:       auth,
		Admin:      admin,
	})
	if err != nil {
		return err
	}

	for _, webhook := range webhooks {
		resp, err := http.Post(webhook.Destination, "application/json", bytes.NewReader(b))
		if err != nil {
			app.Logger().Error("failed to send webhook", "error", err.Error())

			continue
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			app.Logger().Info("webhook sent", "action", action, "name", webhook.Name, "collection", webhook.Collection, "destination", webhook.Destination)
		} else {
			b, _ := io.ReadAll(resp.Body)

			app.Logger().Error("failed to send webhook", "action", action, "name", webhook.Name, "collection", webhook.Collection, "destination", webhook.Destination, "status", resp.Status, "details", string(b))
		}
	}

	return nil
}
