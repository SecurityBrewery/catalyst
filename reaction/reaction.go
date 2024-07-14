package reaction

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"

	"github.com/SecurityBrewery/catalyst/migrations"
)

type WebhookTrigger struct {
	Token string `json:"token"`
	Path  string `json:"path"`
}

func findReaction(app core.App, action string) (*models.Record, *WebhookTrigger, bool, error) {
	exp := dbx.HashExp{"trigger": "webhook"}

	records, err := app.Dao().FindRecordsByExpr(migrations.ReactionCollectionName, exp)
	if err != nil {
		return nil, nil, false, err
	}

	if len(records) == 0 {
		return nil, nil, false, nil
	}

	for _, record := range records {
		triggerdata := record.GetString("triggerdata")

		var trigger WebhookTrigger
		if err := json.Unmarshal([]byte(triggerdata), &trigger); err != nil {
			return nil, nil, false, err
		}

		if trigger.Path == action {
			return record, &trigger, true, nil
		}
	}

	return nil, nil, false, nil
}

type PythonReaction struct {
	Bootstrap string `json:"bootstrap"`
	Script    string `json:"script"`
}

type WebhookReaction struct {
	Headers map[string]string `json:"headers"`
	URL     string            `json:"url"`
}

func runReaction(ctx context.Context, record *models.Record, payload string) ([]byte, error) {
	switch record.GetString("reaction") {
	case "python":
		name := record.GetString("name")
		reactionData := record.GetString("reactiondata")

		var reaction PythonReaction
		if err := json.Unmarshal([]byte(reactionData), &reaction); err != nil {
			return nil, err
		}

		return runPythonReaction(ctx, name, reaction.Bootstrap, reaction.Script, payload)
	case "webhook":
		reactionData := record.GetString("reactiondata")

		var reaction WebhookReaction
		if err := json.Unmarshal([]byte(reactionData), &reaction); err != nil {
			return nil, err
		}

		return runWebhookReaction(ctx, reaction.URL, reaction.Headers, payload)
	default:
		return nil, errors.New("unsupported action type")
	}
}
