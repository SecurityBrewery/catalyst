package fakedata_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/app2/fakedata"
	catalystTesting "github.com/SecurityBrewery/catalyst/testing"
)

func Test_records(t *testing.T) {
	t.Parallel()

	app, _ := catalystTesting.App(t)

	err := fakedata.Generate(t.Context(), app.Queries, 2, 2)
	require.NoError(t, err)
}

func TestGenerate(t *testing.T) {
	t.Parallel()

	app, _ := catalystTesting.App(t)

	err := fakedata.Generate(t.Context(), app.Queries, 0, 0)
	require.NoError(t, err)
}
