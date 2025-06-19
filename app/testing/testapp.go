package testing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/app"
	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/reaction"
)

func App(t *testing.T) (*app.App, func(), *Counter) {
	t.Helper()

	baseApp, cleanup, err := app.New(t.Context(), t.TempDir())
	require.NoError(t, err)

	err = baseApp.SetupRoutes()
	require.NoError(t, err)

	err = reaction.BindHooks(baseApp, true)
	require.NoError(t, err)

	database.DefaultTestData(t, baseApp.Queries)

	counter := countEvents(baseApp)

	return baseApp, cleanup, counter
}

func countEvents(app *app.App) *Counter {
	c := NewCounter()

	app.Hooks.OnRecordsListRequest.Subscribe(count(c, "OnRecordsListRequest"))
	app.Hooks.OnRecordViewRequest.Subscribe(count(c, "OnRecordViewRequest"))
	app.Hooks.OnRecordBeforeCreateRequest.Subscribe(count(c, "OnRecordBeforeCreateRequest"))
	app.Hooks.OnRecordAfterCreateRequest.Subscribe(count(c, "OnRecordAfterCreateRequest"))
	app.Hooks.OnRecordBeforeUpdateRequest.Subscribe(count(c, "OnRecordBeforeUpdateRequest"))
	app.Hooks.OnRecordAfterUpdateRequest.Subscribe(count(c, "OnRecordAfterUpdateRequest"))
	app.Hooks.OnRecordBeforeDeleteRequest.Subscribe(count(c, "OnRecordBeforeDeleteRequest"))
	app.Hooks.OnRecordAfterDeleteRequest.Subscribe(count(c, "OnRecordAfterDeleteRequest"))

	return c
}

func count(c *Counter, name string) func(ctx context.Context, table string, record any) {
	return func(_ context.Context, _ string, _ any) {
		c.Increment(name)
	}
}
