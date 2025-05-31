package reaction

import (
	"github.com/SecurityBrewery/catalyst/app2"
	"github.com/SecurityBrewery/catalyst/reaction/schedule"
	"github.com/SecurityBrewery/catalyst/reaction/trigger/hook"
	"github.com/SecurityBrewery/catalyst/reaction/trigger/webhook"
)

func BindHooks(app *app2.App2, test bool) {
	schedule.Start(app)
	hook.BindHooks(app, test)
	webhook.BindHooks(app)
}
