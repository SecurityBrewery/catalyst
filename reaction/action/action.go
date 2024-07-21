package action

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tokens"
	"github.com/pocketbase/pocketbase/tools/security"

	"github.com/SecurityBrewery/catalyst/migrations"
	"github.com/SecurityBrewery/catalyst/reaction/action/python"
	"github.com/SecurityBrewery/catalyst/reaction/action/webhook"
)

func Run(ctx context.Context, app core.App, actionName, actionData, payload string) ([]byte, error) {
	action, err := decode(actionName, actionData)
	if err != nil {
		return nil, err
	}

	if a, ok := action.(authenticatedAction); ok {
		token, err := systemToken(app)
		if err != nil {
			return nil, fmt.Errorf("failed to get system token: %w", err)
		}

		a.SetToken(token)
	}

	return action.Run(ctx, payload)
}

type action interface {
	Run(ctx context.Context, payload string) ([]byte, error)
}

type authenticatedAction interface {
	SetToken(token string)
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

func systemToken(app core.App) (string, error) {
	authRecord, err := app.Dao().FindAuthRecordByUsername(migrations.UserCollectionName, migrations.SystemUserID)
	if err != nil {
		return "", fmt.Errorf("failed to find system auth record: %w", err)
	}

	return security.NewJWT(
		jwt.MapClaims{
			"id":           authRecord.Id,
			"type":         tokens.TypeAuthRecord,
			"collectionId": authRecord.Collection().Id,
		},
		authRecord.TokenKey()+app.Settings().RecordAuthToken.Secret,
		int64(time.Second*60),
	)
}
