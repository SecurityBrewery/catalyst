package data

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/migration"
	"github.com/SecurityBrewery/catalyst/app/upload/uploader"
)

func NewTestDB(t *testing.T, dir string) *sqlc.Queries {
	t.Helper()

	queries := database.TestDB(t, dir)
	uploader, err := uploader.New(dir)
	require.NoError(t, err)

	err = migration.Apply(t.Context(), queries, dir, uploader)
	require.NoError(t, err)

	DefaultTestData(t, dir, queries)

	return queries
}
