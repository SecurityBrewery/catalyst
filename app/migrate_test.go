package app_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/app"
	"github.com/SecurityBrewery/catalyst/migrations"
	catalystTesting "github.com/SecurityBrewery/catalyst/testing"
)

func Test_MigrateDBsDown(t *testing.T) {
	catalystApp, cleanup := catalystTesting.App(t)
	defer cleanup()

	_, err := catalystApp.Dao().FindCollectionByNameOrId(migrations.ReactionCollectionName)
	require.NoError(t, err)

	require.NoError(t, app.MigrateDBsDown(catalystApp))
}
