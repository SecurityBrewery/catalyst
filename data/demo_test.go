package data_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/data"
	catalystTesting "github.com/SecurityBrewery/catalyst/testing"
)

func TestGenerate(t *testing.T) {
	t.Parallel()

	app, _ := catalystTesting.App(t)

	_ = app.Queries.DeleteUser(t.Context(), "u_admin")
	_ = app.Queries.DeleteUser(t.Context(), "u_bob_analyst")
	_ = app.Queries.DeleteRole(t.Context(), "r_admin")
	_ = app.Queries.DeleteRole(t.Context(), "r_analyst")

	err := data.GenerateDemoData(t.Context(), app.Queries, 0, 0)
	require.NoError(t, err, "failed to generate fake data")
}
