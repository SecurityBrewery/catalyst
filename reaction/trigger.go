package reaction

import (
	"github.com/pocketbase/pocketbase/core"

	"github.com/SecurityBrewery/catalyst/reaction/trigger/hook"
	"github.com/SecurityBrewery/catalyst/reaction/trigger/webhook"
)

func BindHooks(app core.App) {
	hook.BindHooks(app)
	webhook.BindHooks(app)
}
