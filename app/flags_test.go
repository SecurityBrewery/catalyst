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

	// stage 1
	require.NoError(t, app.SetFlags(catalystApp, []string{"test"}))

	got, err := app.Flags(catalystApp)
	require.NoError(t, err)

	assert.ElementsMatch(t, []string{"test"}, got)

	// stage 2
	require.NoError(t, app.SetFlags(catalystApp, []string{"test2"}))

	got, err = app.Flags(catalystApp)
	require.NoError(t, err)

	assert.ElementsMatch(t, []string{"test2"}, got)

	// stage 3
	require.NoError(t, app.SetFlags(catalystApp, []string{"test", "test2"}))

	got, err = app.Flags(catalystApp)
	require.NoError(t, err)

	assert.ElementsMatch(t, []string{"test", "test2"}, got)
}
