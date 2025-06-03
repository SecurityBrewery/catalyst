package database

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

func NewTestDB(t *testing.T) *sqlc.Queries {
	queries, _, err := DB(filepath.Join(t.TempDir(), "data.db"))
	require.NoError(t, err, "failed to create test database")

	DefaultTestData(t, queries)

	return queries
}
