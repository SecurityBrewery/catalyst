package data

import (
	"testing"

	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

func NewTestDB(t *testing.T) *sqlc.Queries {
	t.Helper()

	queries := database.TestDB(t)

	DefaultTestData(t, queries)

	return queries
}
