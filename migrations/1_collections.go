package migrations

import (
	"fmt"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	CommentCollectionName  = "comments"
	FeatureCollectionName  = "features"
	LinkCollectionName     = "links"
	TaskCollectionName     = "tasks"
	TicketCollectionName   = "tickets"
	TimelineCollectionName = "timeline"
	TypeCollectionName     = "types"
	WebhookCollectionName  = "webhooks"
	fileCollectionName     = "files"

	UserCollectionID = "_pb_users_auth_"
)

func collectionsUp(app core.App) error {
	typeCollection := core.NewBaseCollection(TypeCollectionName)
	typeCollection.Fields.Add(&core.TextField{Name: "singular", Required: true})
	typeCollection.Fields.Add(&core.TextField{Name: "plural", Required: true})
	typeCollection.Fields.Add(&core.TextField{Name: "icon", Required: true})
	typeCollection.Fields.Add(&core.JSONField{Name: "schema", Required: true, MaxSize: 50_000})

	ticketCollection := core.NewBaseCollection(TicketCollectionName)
	ticketCollection.Fields.Add(&core.TextField{Name: "name", Required: true})
	ticketCollection.Fields.Add(&core.RelationField{Name: "type", CollectionId: typeCollection.Id, Required: true})
	ticketCollection.Fields.Add(&core.TextField{Name: "description"})
	ticketCollection.Fields.Add(&core.BoolField{Name: "open"})
	ticketCollection.Fields.Add(&core.TextField{Name: "resolution"})
	ticketCollection.Fields.Add(&core.JSONField{Name: "schema", MaxSize: 50_000})
	ticketCollection.Fields.Add(&core.JSONField{Name: "state", MaxSize: 50_000})
	ticketCollection.Fields.Add(&core.RelationField{Name: "owner", CollectionId: UserCollectionID})

	taskCollection := core.NewBaseCollection(TaskCollectionName)
	taskCollection.Fields.Add(&core.RelationField{Name: "ticket", CollectionId: ticketCollection.Id, Required: true})
	taskCollection.Fields.Add(&core.TextField{Name: "name", Required: true})
	taskCollection.Fields.Add(&core.BoolField{Name: "open"})
	taskCollection.Fields.Add(&core.RelationField{Name: "owner", CollectionId: UserCollectionID})

	commentCollection := core.NewBaseCollection(CommentCollectionName)
	commentCollection.Fields.Add(&core.RelationField{Name: "ticket", CollectionId: ticketCollection.Id, Required: true})
	commentCollection.Fields.Add(&core.RelationField{Name: "author", CollectionId: UserCollectionID})
	commentCollection.Fields.Add(&core.TextField{Name: "message", Required: true})

	timelineCollection := core.NewBaseCollection(TimelineCollectionName)
	timelineCollection.Fields.Add(&core.RelationField{Name: "ticket", CollectionId: ticketCollection.Id, Required: true})
	timelineCollection.Fields.Add(&core.DateField{Name: "time", Required: true})
	timelineCollection.Fields.Add(&core.TextField{Name: "message", Required: true})

	linkCollection := core.NewBaseCollection(LinkCollectionName)
	linkCollection.Fields.Add(&core.RelationField{Name: "ticket", CollectionId: ticketCollection.Id, Required: true})
	linkCollection.Fields.Add(&core.TextField{Name: "name", Required: true})
	linkCollection.Fields.Add(&core.URLField{Name: "url", Required: true})

	fileCollection := core.NewBaseCollection(fileCollectionName)
	fileCollection.Fields.Add(&core.RelationField{Name: "ticket", CollectionId: ticketCollection.Id, Required: true})
	fileCollection.Fields.Add(&core.TextField{Name: "name", Required: true})
	fileCollection.Fields.Add(&core.NumberField{Name: "size", Required: true})
	fileCollection.Fields.Add(&core.FileField{Name: "blob", Required: true, MaxSize: 1024 * 1024 * 100})

	featureCollection := core.NewBaseCollection(FeatureCollectionName)
	featureCollection.Fields.Add(&core.TextField{Name: "name", Required: true})
	featureCollection.Fields.Add(&core.AutodateField{Name: "created", OnCreate: true})
	featureCollection.Fields.Add(&core.AutodateField{Name: "updated", OnUpdate: true})
	featureCollection.ListRule = types.Pointer("@request.auth.id != ''")
	featureCollection.ViewRule = types.Pointer("@request.auth.id != ''")
	featureCollection.Indexes = types.JSONArray[string]{
		fmt.Sprintf("CREATE UNIQUE INDEX `unique_name` ON `%s` (`name`)", FeatureCollectionName),
	}

	collections := []*core.Collection{
		internalCollection(typeCollection),
		internalCollection(ticketCollection),
		internalCollection(taskCollection),
		internalCollection(commentCollection),
		internalCollection(timelineCollection),
		internalCollection(linkCollection),
		internalCollection(fileCollection),
		featureCollection,
	}

	for _, c := range collections {
		if err := app.Save(c); err != nil {
			return err
		}
	}

	return nil
}

func internalCollection(c *core.Collection) *core.Collection {
	c.Fields.Add(&core.AutodateField{Name: "created", OnCreate: true})
	c.Fields.Add(&core.AutodateField{Name: "updated", OnUpdate: true})

	c.ListRule = types.Pointer("@request.auth.id != ''")
	c.ViewRule = types.Pointer("@request.auth.id != ''")
	c.CreateRule = types.Pointer("@request.auth.id != ''")
	c.UpdateRule = types.Pointer("@request.auth.id != ''")
	c.DeleteRule = types.Pointer("@request.auth.id != ''")

	return c
}

func collectionsDown(app core.App) error {
	collections := []string{
		fileCollectionName,
		LinkCollectionName,
		TaskCollectionName,
		CommentCollectionName,
		TimelineCollectionName,
		FeatureCollectionName,
		TicketCollectionName,
		TypeCollectionName,
	}

	for _, name := range collections {
		id, err := app.FindCollectionByNameOrId(name)
		if err != nil {
			return err
		}

		if err := app.Delete(id); err != nil {
			return err
		}
	}

	return nil
}
