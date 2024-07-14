package fakedata_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/fakedata"
	catalystTesting "github.com/SecurityBrewery/catalyst/testing"
)

func Test_records(t *testing.T) {
	app, cleanup := catalystTesting.App(t)
	defer cleanup()

	got, err := fakedata.Records(app, 2, 2)
	require.NoError(t, err)

	assert.Greater(t, len(got), 2)
}

func TestGenerate(t *testing.T) {
	app, cleanup := catalystTesting.App(t)
	defer cleanup()

	err := fakedata.Generate(app, 0, 0)
	require.NoError(t, err)
}
