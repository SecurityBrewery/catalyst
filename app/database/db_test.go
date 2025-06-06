package database_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

func TestDBInitialization(t *testing.T) {
	t.Parallel()

	queries, cleanup, err := database.DB(filepath.Join(t.TempDir(), "data.db"))
	require.NoError(t, err)
	t.Cleanup(cleanup)

	user, err := queries.SystemUser(t.Context())
	require.NoError(t, err)
	assert.Equal(t, "system", user.ID)

	types, err := queries.ListTypes(t.Context(), sqlc.ListTypesParams{Offset: 0, Limit: 10})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(types), 1)
}

func TestDBForeignKeyConstraints(t *testing.T) {
	t.Parallel()

	queries, cleanup, err := database.DB(filepath.Join(t.TempDir(), "data.db"))
	require.NoError(t, err)
	t.Cleanup(cleanup)

	err = queries.AssignRoleToUser(t.Context(), sqlc.AssignRoleToUserParams{
		UserID: "does_not_exist",
		RoleID: "also_missing",
	})
	assert.Error(t, err)
}

func TestNewTestDBDefaultData(t *testing.T) {
	t.Parallel()

	queries := database.NewTestDB(t)

	user, err := queries.UserByEmail(t.Context(), database.AdminEmail)
	require.NoError(t, err)
	assert.Equal(t, database.AdminEmail, user.Email)

	ticket, err := queries.Ticket(t.Context(), "test-ticket")
	require.NoError(t, err)
	assert.Equal(t, "test-ticket", ticket.ID)

	comment, err := queries.GetComment(t.Context(), "c_test_comment")
	require.NoError(t, err)
	assert.Equal(t, "c_test_comment", comment.ID)

	timeline, err := queries.GetTimeline(t.Context(), "h_test_timeline")
	require.NoError(t, err)
	assert.Equal(t, "h_test_timeline", timeline.ID)
}
