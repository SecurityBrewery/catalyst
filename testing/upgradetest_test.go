package testing

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/app/data"
)

func TestUpgradeTestData(t *testing.T) {
	t.Parallel()

	app, cleanup, _ := App(t)

	t.Cleanup(cleanup)

	require.NoError(t, data.GenerateUpgradeTestData(t.Context(), app.Queries))

	ValidateUpgradeTestData(t, app.Queries)
}
