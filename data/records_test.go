package data_test

import (
	"testing"

	"github.com/SecurityBrewery/catalyst/data"
	catalystTesting "github.com/SecurityBrewery/catalyst/testing"
)

func Test_records(t *testing.T) {
	t.Parallel()

	app, _ := catalystTesting.App(t)

	data.GenerateFake(t, app.Queries, 2, 2)
}

func TestGenerate(t *testing.T) {
	t.Parallel()

	app, _ := catalystTesting.App(t)

	data.GenerateFake(t, app.Queries, 0, 0)
}
