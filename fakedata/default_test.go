package fakedata_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/fakedata"
	catalystTesting "github.com/SecurityBrewery/catalyst/testing"
)

func TestDefaultData(t *testing.T) {
	t.Parallel()

	app, _, cleanup := catalystTesting.App(t)
	defer cleanup()

	require.NoError(t, fakedata.GenerateDefaultData(app))

	require.NoError(t, fakedata.ValidateDefaultData(app))
}
