package action

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SecurityBrewery/catalyst/reaction/action/python"
	"github.com/SecurityBrewery/catalyst/reaction/action/webhook"
)

func Run(ctx context.Context, actionName, actionData, payload string) ([]byte, error) {
	action, err := decode(actionName, actionData)
	if err != nil {
		return nil, err
	}

	return action.Run(ctx, payload)
}

type action interface {
	Run(ctx context.Context, payload string) ([]byte, error)
}

func decode(actionName, actionData string) (action, error) {
	switch actionName {
	case "python":
		var reaction python.Python
		if err := json.Unmarshal([]byte(actionData), &reaction); err != nil {
			return nil, err
		}

		return &reaction, nil
	case "webhook":
		var reaction webhook.Webhook
		if err := json.Unmarshal([]byte(actionData), &reaction); err != nil {
			return nil, err
		}

		return &reaction, nil
	default:
		return nil, fmt.Errorf("action %q not found", actionName)
	}
}
