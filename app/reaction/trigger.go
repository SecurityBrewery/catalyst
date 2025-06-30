package reaction

import (
	"github.com/go-chi/chi/v5"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/hook"
	reactionHook "github.com/SecurityBrewery/catalyst/app/reaction/trigger/hook"
	"github.com/SecurityBrewery/catalyst/app/reaction/trigger/webhook"
)

func BindHooks(hooks *hook.Hooks, router chi.Router, queries *sqlc.Queries, test bool) error {
	reactionHook.BindHooks(hooks, queries, test)
	webhook.BindHooks(router, queries)

	return nil
}
