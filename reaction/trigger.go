package reaction

import (
	"github.com/pocketbase/pocketbase"

	"github.com/SecurityBrewery/catalyst/reaction/schedule"
	"github.com/SecurityBrewery/catalyst/reaction/trigger/hook"
	"github.com/SecurityBrewery/catalyst/reaction/trigger/webhook"
)

func BindHooks(pb *pocketbase.PocketBase, test bool) {
	schedule.Start(pb)
	hook.BindHooks(pb, test)
	webhook.BindHooks(pb)
}
