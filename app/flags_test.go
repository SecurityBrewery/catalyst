package app_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/app"
	catalystTesting "github.com/SecurityBrewery/catalyst/testing"
)

func Test_flags(t *testing.T) {
	t.Parallel()

	catalystApp, _, cleanup := catalystTesting.App(t)
	defer cleanup()

	got, err := app.Flags(catalystApp)
	require.NoError(t, err)

	want := []string{}
	assert.ElementsMatch(t, want, got)
}

func Test_setFlags(t *testing.T) {
	t.Parallel()

	catalystApp, _, cleanup := catalystTesting.App(t)
	defer cleanup()

	require.NoError(t, app.SetFlags(catalystApp, []string{"test"}))

	got, err := app.Flags(catalystApp)
	require.NoError(t, err)

	assert.ElementsMatch(t, []string{"test"}, got)

	require.NoError(t, app.SetFlags(catalystApp, []string{"test2"}))

	got, err = app.Flags(catalystApp)
	require.NoError(t, err)

	assert.ElementsMatch(t, []string{"test2"}, got)
}
