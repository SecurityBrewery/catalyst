package app

import (
	"fmt"
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

	webhook.BindHooks(app)
	reaction.BindHooks(app, test)

	app.OnBeforeServe().Add(addRoutes())

	app.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
		if HasFlag(e.App, "demo") {
			bindDemoHooks(e.App)
		}

		return nil
	})

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

func bindDemoHooks(app core.App) {
	app.OnRecordBeforeCreateRequest("files", "reactions").Add(func(e *core.RecordCreateEvent) error {
		return fmt.Errorf("cannot create %s in demo mode", e.Record.Collection().Name)
	})
	app.OnRecordBeforeUpdateRequest("files", "reactions").Add(func(e *core.RecordUpdateEvent) error {
		return fmt.Errorf("cannot update %s in demo mode", e.Record.Collection().Name)
	})
	app.OnRecordBeforeDeleteRequest("files", "reactions").Add(func(e *core.RecordDeleteEvent) error {
		return fmt.Errorf("cannot delete %s in demo mode", e.Record.Collection().Name)
	})
}

func dev() bool {
	return strings.HasPrefix(os.Args[0], os.TempDir())
}
