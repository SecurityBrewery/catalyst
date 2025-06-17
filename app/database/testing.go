package database

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

func NewTestDB(t *testing.T) *sqlc.Queries {
	t.Helper()

	tmpDBPath := filepath.Join(t.TempDir(), "data.db")
	queries, cleanup, err := DB(t.Context(), tmpDBPath)
	require.NoError(t, err, "failed to create test database")

	t.Cleanup(cleanup)

	DefaultTestData(t, queries)

	return queries
}
