package app

import (
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"

	"github.com/SecurityBrewery/catalyst/migrations"
	"github.com/SecurityBrewery/catalyst/reaction"
	"github.com/SecurityBrewery/catalyst/webhook"
)

func init() {
	migrations.Register()
}

func App(dir string) *pocketbase.PocketBase {
	app := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDev:     dev(),
		DefaultDataDir: dir,
	})

	BindHooks(app)

	// Register additional commands
	app.RootCmd.AddCommand(bootstrapCmd(app))
	app.RootCmd.AddCommand(fakeDataCmd(app))
	app.RootCmd.AddCommand(setFeatureFlagsCmd(app))

	return app
}

func BindHooks(app core.App) {
	webhook.BindHooks(app)
	reaction.BindHooks(app)

	app.OnBeforeServe().Add(addRoutes())
}

func dev() bool {
	return strings.HasPrefix(os.Args[0], os.TempDir())
}
