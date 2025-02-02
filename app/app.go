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

func init() { //nolint:gochecknoinits
	migrations.Register()
}

func App(dir string, test bool) (*pocketbase.PocketBase, error) {
	app := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDev:     test || dev(),
		DefaultDataDir: dir,
	})

	var appURL string

	app.RootCmd.PersistentFlags().StringVar(&appURL, "app-url", "", "the app's URL")

	var flags []string

	app.RootCmd.PersistentFlags().StringSliceVar(&flags, "flags", nil, "feature flags")

	_ = app.RootCmd.ParseFlags(os.Args[1:])

	app.RootCmd.AddCommand(fakeDataCmd(app))
	app.RootCmd.AddCommand(defaultDataCmd(app))

	webhook.BindHooks(app)
	reaction.BindHooks(app, test)

	app.OnBootstrap().BindFunc(func(e *core.BootstrapEvent) error {
		if err := e.Next(); err != nil {
			return err
		}

		return MigrateDBs(e.App)
	})

	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		if err := setupServer(appURL, flags)(e); err != nil {
			return err
		}

		return e.Next()
	})

	return app, nil
}

func dev() bool {
	return strings.HasPrefix(os.Args[0], os.TempDir())
}
