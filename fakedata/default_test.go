package fakedata_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/fakedata"
	catalystTesting "github.com/SecurityBrewery/catalyst/testing"
)

func TestDefaultData(t *testing.T) {
	t.Parallel()

	app, _ := catalystTesting.App(t)

	require.NoError(t, fakedata.GenerateDefaultData(t.Context(), app.Queries))

	require.NoError(t, fakedata.ValidateDefaultData(t.Context(), t, app))
}
