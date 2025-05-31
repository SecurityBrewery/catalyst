package action

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/reaction/action/python"
	"github.com/SecurityBrewery/catalyst/reaction/action/webhook"
)

func Run(ctx context.Context, queries *sqlc.Queries, actionName, actionData, payload string) ([]byte, error) {
	action, err := decode(actionName, actionData)
	if err != nil {
		return nil, err
	}

	if a, ok := action.(authenticatedAction); ok {
		token, err := systemToken(ctx, queries)
		if err != nil {
			return nil, fmt.Errorf("failed to get system token: %w", err)
		}

		a.SetEnv([]string{
			"CATALYST_APP_URL=", // + app.Settings().Meta.AppUrl, // TODO
			"CATALYST_TOKEN=" + token,
		})
	}

	return action.Run(ctx, payload)
}

type action interface {
	Run(ctx context.Context, payload string) ([]byte, error)
}

type authenticatedAction interface {
	SetEnv(env []string)
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

func systemToken(ctx context.Context, queries *sqlc.Queries) (string, error) {
	_, err := queries.UserByUserName(ctx, "system")
	if err != nil {
		return "", fmt.Errorf("failed to find system auth record: %w", err)
	}

	/* TODO
	return security.NewJWT(
		jwt.MapClaims{
			"id":           authRecord.ID,
			"type":         tokens.TypeAuthRecord,
			"collectionId": "users",
		},
		authRecord.Tokenkey, // +app.Settings().RecordAuthToken.Secret,
		int64(time.Second*60),
	)
	*/
	return "", nil
}
