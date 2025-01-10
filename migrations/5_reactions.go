package migrations

import (
	"fmt"

	"github.com/pocketbase/pocketbase/core"
)

const ReactionCollectionName = "reactions"

func reactionsUp(app core.App) error {
	triggers := []string{"webhook", "hook"}
	reactions := []string{"python", "webhook"}

	reactionCollection := core.NewBaseCollection(ReactionCollectionName)
	reactionCollection.Fields.Add(&core.TextField{Name: "name", Required: true})
	reactionCollection.Fields.Add(&core.SelectField{Name: "trigger", Required: true, Values: triggers})
	reactionCollection.Fields.Add(&core.JSONField{Name: "triggerdata", Required: true, MaxSize: 50_000})
	reactionCollection.Fields.Add(&core.SelectField{Name: "action", Required: true, Values: reactions})
	reactionCollection.Fields.Add(&core.JSONField{Name: "actiondata", Required: true, MaxSize: 50_000})

	return app.Save(internalCollection(reactionCollection))
}

func reactionsDown(app core.App) error {
	id, err := app.FindCollectionByNameOrId(ReactionCollectionName)
	if err != nil {
		return fmt.Errorf("failed to find collection %s: %w", ReactionCollectionName, err)
	}

	return app.Delete(id)
}
