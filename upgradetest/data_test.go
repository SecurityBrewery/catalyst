package upgradetest

import (
	"testing"

	"github.com/stretchr/testify/require"

	catalystTesting "github.com/SecurityBrewery/catalyst/testing"
)

func TestUpgradeTestData(t *testing.T) {
	t.Parallel()

	app, _ := catalystTesting.App(t)

	require.NoError(t, generateUpgradeTestData(t.Context(), app.Queries))

	validateUpgradeTestData(t, app)
}
