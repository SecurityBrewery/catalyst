package reaction

import (
	"github.com/SecurityBrewery/catalyst/app"
	"github.com/SecurityBrewery/catalyst/reaction/trigger/hook"
	"github.com/SecurityBrewery/catalyst/reaction/trigger/webhook"
)

func BindHooks(app *app.App, test bool) error {
	hook.BindHooks(app, test)
	webhook.BindHooks(app)

	return nil
}
