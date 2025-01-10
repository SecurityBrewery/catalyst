package reaction

import (
	"github.com/pocketbase/pocketbase/core"

	"github.com/SecurityBrewery/catalyst/reaction/schedule"
	"github.com/SecurityBrewery/catalyst/reaction/trigger/hook"
	"github.com/SecurityBrewery/catalyst/reaction/trigger/webhook"
)

func BindHooks(app core.App, test bool) {
	schedule.Start(app)
	hook.BindHooks(app, test)
	webhook.BindHooks(app)
}
