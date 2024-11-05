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

	var appURL string
	app.RootCmd.PersistentFlags().StringVar(&appURL, "app-url", "", "the app's URL")

	var flags []string
	app.RootCmd.PersistentFlags().StringSliceVar(&flags, "flags", nil, "feature flags")

	_ = app.RootCmd.ParseFlags(os.Args[1:])

	app.RootCmd.AddCommand(fakeDataCmd(app))

	webhook.BindHooks(app)
	reaction.BindHooks(app, test)

	app.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
		if err := MigrateDBs(e.App); err != nil {
			return err
		}

		if err := SetFlags(e.App, flags); err != nil {
			return err
		}

		if HasFlag(e.App, "demo") {
			bindDemoHooks(e.App)
		}

		if appURL != "" {
			s := e.App.Settings()
			s.Meta.AppUrl = appURL
			if err := e.App.Dao().SaveSettings(s); err != nil {
				return err
			}
		}

		return e.App.RefreshSettings()
	})

	app.OnBeforeServe().Add(addRoutes())

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
