package app

import (
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"

	"github.com/SecurityBrewery/catalyst/migrations"
	"github.com/SecurityBrewery/catalyst/reaction"
	"github.com/SecurityBrewery/catalyst/webhook"
)

func init() { //nolint:gochecknoinits
	migrations.Register()
}

func App(dir string, test bool) (*pocketbase.PocketBase, error) {
	app := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDev:     false, // dev(),
		DefaultDataDir: dir,
	})

	webhook.BindHooks(app)
	reaction.BindHooks(app, test)

	app.OnBeforeServe().Add(addRoutes())

	// Register additional commands
	app.RootCmd.AddCommand(fakeDataCmd(app))
	app.RootCmd.AddCommand(setFeatureFlagsCmd(app))

	if err := app.Bootstrap(); err != nil {
		return nil, err
	}

	if err := MigrateDBs(app); err != nil {
		return nil, err
	}

	return app, nil
}

func dev() bool {
	return strings.HasPrefix(os.Args[0], os.TempDir())
}
