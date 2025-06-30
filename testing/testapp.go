package testing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/app"
	"github.com/SecurityBrewery/catalyst/app/counter"
	"github.com/SecurityBrewery/catalyst/app/data"
	"github.com/SecurityBrewery/catalyst/app/hook"
)

func App(t *testing.T) (*app.App, func(), *counter.Counter) {
	t.Helper()

	dir := t.TempDir()

	catalyst, cleanup, err := app.New(t.Context(), dir)
	require.NoError(t, err)

	data.DefaultTestData(t, dir, catalyst.Queries)

	return catalyst, cleanup, countEvents(catalyst.Hooks)
}

func countEvents(hooks *hook.Hooks) *counter.Counter {
	c := counter.NewCounter()

	hooks.OnRecordsListRequest.Subscribe(count(c, "OnRecordsListRequest"))
	hooks.OnRecordViewRequest.Subscribe(count(c, "OnRecordViewRequest"))
	hooks.OnRecordBeforeCreateRequest.Subscribe(count(c, "OnRecordBeforeCreateRequest"))
	hooks.OnRecordAfterCreateRequest.Subscribe(count(c, "OnRecordAfterCreateRequest"))
	hooks.OnRecordBeforeUpdateRequest.Subscribe(count(c, "OnRecordBeforeUpdateRequest"))
	hooks.OnRecordAfterUpdateRequest.Subscribe(count(c, "OnRecordAfterUpdateRequest"))
	hooks.OnRecordBeforeDeleteRequest.Subscribe(count(c, "OnRecordBeforeDeleteRequest"))
	hooks.OnRecordAfterDeleteRequest.Subscribe(count(c, "OnRecordAfterDeleteRequest"))

	return c
}

func count(c *counter.Counter, name string) func(ctx context.Context, table string, record any) {
	return func(_ context.Context, _ string, _ any) {
		c.Increment(name)
	}
}
