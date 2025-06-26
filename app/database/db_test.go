package database_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/app/data"
	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/pointer"
)

func TestDBInitialization(t *testing.T) {
	t.Parallel()

	queries := database.TestDB(t, t.TempDir())

	user, err := queries.SystemUser(t.Context())
	require.NoError(t, err)
	assert.Equal(t, "system", user.ID)

	types, err := queries.ListTypes(t.Context(), sqlc.ListTypesParams{Offset: 0, Limit: 10})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(types), 1)
}

func TestDBForeignKeyConstraints(t *testing.T) {
	t.Parallel()

	queries := database.TestDB(t, t.TempDir())

	assert.Error(t, queries.AssignGroupToUser(t.Context(), sqlc.AssignGroupToUserParams{
		UserID:  "does_not_exist",
		GroupID: "also_missing",
	}))
}

func TestNewTestDBDefaultData(t *testing.T) {
	t.Parallel()

	queries := data.NewTestDB(t, t.TempDir())

	user, err := queries.UserByEmail(t.Context(), pointer.Pointer(data.AdminEmail))
	require.NoError(t, err)
	assert.Equal(t, data.AdminEmail, *user.Email)

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

func TestReadWrite(t *testing.T) {
	t.Parallel()

	queries := data.NewTestDB(t, t.TempDir())

	for range 3 {
		y, err := queries.CreateType(t.Context(), sqlc.CreateTypeParams{
			Singular: "Foo",
			Plural:   "Foos",
			Icon:     pointer.Pointer("Bug"),
			Schema:   json.RawMessage("{}"),
		})
		require.NoError(t, err)

		_, err = queries.GetType(t.Context(), y.ID)
		require.NoError(t, err)

		err = queries.DeleteType(t.Context(), y.ID)
		require.NoError(t, err)
	}
}

func TestRead(t *testing.T) {
	t.Parallel()

	queries := data.NewTestDB(t, t.TempDir())

	// read from a table
	_, err := queries.GetUser(t.Context(), "u_bob_analyst")
	require.NoError(t, err)

	// read from a view
	_, err = queries.GetSidebar(t.Context())
	require.NoError(t, err)
}
