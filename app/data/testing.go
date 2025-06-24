package data

import (
	"testing"

	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

func NewTestDB(t *testing.T, dir string) *sqlc.Queries {
	t.Helper()

	queries := database.TestDB(t, dir)

	DefaultTestData(t, dir, queries)

	return queries
}
