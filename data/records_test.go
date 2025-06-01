package data_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/data"
	catalystTesting "github.com/SecurityBrewery/catalyst/testing"
)

func Test_records(t *testing.T) {
	t.Parallel()

	app, _ := catalystTesting.App(t)

	err := data.GenerateFake(t.Context(), app.Queries, 2, 2)
	require.NoError(t, err, "failed to generate fake data")
}

func TestGenerate(t *testing.T) {
	t.Parallel()

	app, _ := catalystTesting.App(t)

	err := data.GenerateFake(t.Context(), app.Queries, 0, 0)
	require.NoError(t, err, "failed to generate fake data")
}
