package upgradetest

import (
	"testing"

	"github.com/stretchr/testify/require"

	catalystTesting "github.com/SecurityBrewery/catalyst/app/testing"
)

func TestUpgradeTestData(t *testing.T) {
	t.Parallel()

	app, cleanup, _ := catalystTesting.App(t)

	t.Cleanup(cleanup)

	require.NoError(t, GenerateUpgradeTestData(t.Context(), app.Queries))

	validateUpgradeTestData(t, app)
}
