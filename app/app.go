package app

import (
	"context"
	"fmt"
	"github.com/SecurityBrewery/catalyst/analysis"
	"net/url"
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

	ctx := context.Background()
	engine := analysis.NewEngine(ctx, app.Dao())

	app.OnBeforeServe().Add(addRoutes(engine))

	app.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
		if HasFlag(e.App, "demo") {
			bindDemoHooks(e.App)
		}

		return nil
	})

	app.OnRecordBeforeCreateRequest("resources").Add(func(e *core.RecordCreateEvent) error {
		e

		if err := engine.SetDao(e.App.Dao()); err != nil {
			return err
		}

		enrichments, err := engine.Enrich(
			e.HttpContext.Request().Context(),
			e.Record.GetString("value"),
			1,
		)
		if err != nil || len(enrichments) == 0 {
			urlx := ""
			resourceType := "unknown"
			icon := "Box"
			if _, err := url.Parse(e.Record.GetString("value")); err == nil {
				urlx = e.Record.GetString("value")
				resourceType = "url"
				icon = "Link"
			}

			e.Record.Set("service", "catalyst")
			e.Record.Set("type", resourceType)
			e.Record.Set("resource", e.Record.GetString("value"))
			e.Record.Set("name", e.Record.GetString("value"))
			e.Record.Set("icon", icon)
			e.Record.Set("description", "")
			e.Record.Set("url", urlx)
			e.Record.Set("attributes", "")

			return nil
		}

		enrichment := enrichments[0]

		e.Record.Set("service", enrichment.ServiceID)
		e.Record.Set("type", enrichment.Resource.Type)
		e.Record.Set("resource", enrichment.Resource.ID)
		e.Record.Set("name", enrichment.Resource.Name)
		e.Record.Set("icon", enrichment.Resource.Icon)
		e.Record.Set("description", enrichment.Resource.Description)
		e.Record.Set("url", enrichment.Resource.URL)
		e.Record.Set("attributes", enrichment.Resource.Attributes)

		return nil
	})

	// Register additional commands
	app.RootCmd.AddCommand(fakeDataCmd(app))
	app.RootCmd.AddCommand(setFeatureFlagsCmd(app))
	app.RootCmd.AddCommand(setAppURL(app))

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
