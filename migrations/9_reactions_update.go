package migrations

import (
	"fmt"

	"github.com/pocketbase/pocketbase/core"
)

func reactionsUpdateUp(app core.App) error {
	triggers := []string{"webhook", "hook", "schedule"}

	col, err := app.FindCollectionByNameOrId(ReactionCollectionName)
	if err != nil {
		return fmt.Errorf("failed to find collection %s: %w", ReactionCollectionName, err)
	}

	field := col.Fields.GetByName("trigger")

	selectField, ok := field.(*core.SelectField)
	if !ok {
		return fmt.Errorf("field %s is not a select field", field.GetName())
	}

	selectField.Values = triggers

	col.Fields.Add(field)

	return app.Save(col)
}
