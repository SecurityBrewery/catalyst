package migrations

import (
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models/schema"
)

func reactionsUpdateUp(db dbx.Builder) error {
	dao := daos.New(db)

	triggers := []string{"webhook", "hook", "schedule"}

	col, err := dao.FindCollectionByNameOrId(ReactionCollectionName)
	if err != nil {
		return fmt.Errorf("failed to find collection %s: %w", ReactionCollectionName, err)
	}

	field := col.Schema.GetFieldByName("trigger")

	field.Options = &schema.SelectOptions{MaxSelect: 1, Values: triggers}

	col.Schema.AddField(field)

	return dao.SaveCollection(col)
}
