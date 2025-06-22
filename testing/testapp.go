package testing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/app"
	"github.com/SecurityBrewery/catalyst/app/counter"
	"github.com/SecurityBrewery/catalyst/app/data"
	"github.com/SecurityBrewery/catalyst/app/reaction"
)

func App(t *testing.T) (*app.App, func(), *counter.Counter) {
	t.Helper()

	baseApp, cleanup, err := app.New(t.Context(), t.TempDir())
	require.NoError(t, err)

	err = reaction.BindHooks(baseApp, true)
	require.NoError(t, err)

	data.DefaultTestData(t, baseApp.Queries)

	return baseApp, cleanup, countEvents(baseApp)
}

func countEvents(app *app.App) *counter.Counter {
	c := counter.NewCounter()

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

func count(c *counter.Counter, name string) func(ctx context.Context, table string, record any) {
	return func(_ context.Context, _ string, _ any) {
		c.Increment(name)
	}
}
