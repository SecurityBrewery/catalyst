package app_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/app"
	"github.com/SecurityBrewery/catalyst/migrations"
	catalystTesting "github.com/SecurityBrewery/catalyst/testing"
)

func Test_MigrateDBsDown(t *testing.T) {
	t.Parallel()

	catalystApp, _, cleanup := catalystTesting.App(t)
	defer cleanup()

	_, err := catalystApp.FindCollectionByNameOrId(migrations.ReactionCollectionName)
	require.NoError(t, err)

	require.NoError(t, app.MigrateDBsDown(catalystApp))
}
