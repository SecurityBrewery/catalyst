package app

import (
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"

	"github.com/SecurityBrewery/catalyst/migrations"
)

func App(dir string) *pocketbase.PocketBase {
	migrations.Register()

	app := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDev:     dev(),
		DefaultDataDir: dir,
	})

	attachWebhooks(app)

	app.OnBeforeServe().Add(addRoutes())

	// Register additional commands
	app.RootCmd.AddCommand(bootstrapCmd(app))
	app.RootCmd.AddCommand(fakeDataCmd(app))
	app.RootCmd.AddCommand(setFeatureFlagsCmd(app))

	return app
}

func dev() bool {
	return strings.HasPrefix(os.Args[0], os.TempDir())
}
