package data_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/data"
	catalystTesting "github.com/SecurityBrewery/catalyst/testing"
)

func TestDefaultData(t *testing.T) {
	t.Parallel()

	app, _ := catalystTesting.App(t)

	require.NoError(t, data.GenerateDefaultData(t.Context(), app.Queries))

	data.ValidateDefaultData(t, app)
}
