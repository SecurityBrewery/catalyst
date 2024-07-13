package reaction

import (
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"

	"github.com/SecurityBrewery/catalyst/migrations"
)

func findReaction(app core.App, action string) (*models.Record, bool, error) {
	records, err := app.Dao().FindRecordsByExpr(migrations.ReactionCollectionName, dbx.HashExp{"name": action})
	if err != nil {
		return nil, false, err
	}

	if len(records) == 0 {
		return nil, false, nil
	}

	return records[0], true, nil
}

func runReaction(actionType, name, bootstrap, script, payload string) ([]byte, error) {
	switch actionType {
	case "python":
		return runPythonReaction(name, bootstrap, script, payload)
	default:
		return nil, errors.New("unsupported action type")
	}
}
