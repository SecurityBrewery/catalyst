package database_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

func TestDBForeignKeyConstraints(t *testing.T) {
	t.Parallel()

	queries := database.TestDB(t, t.TempDir())

	assert.Error(t, queries.AssignGroupToUser(t.Context(), sqlc.AssignGroupToUserParams{
		UserID:  "does_not_exist",
		GroupID: "also_missing",
	}))
}
