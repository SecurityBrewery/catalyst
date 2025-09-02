package data

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/pointer"
)

func TestNewTestDB(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	queries := NewTestDB(t, dir)

	user, err := queries.GetUser(t.Context(), "u_bob_analyst")
	require.NoError(t, err)

	assert.Equal(t, "u_bob_analyst", user.ID)
	assert.Equal(t, "Bob Analyst", *user.Name)
	assert.Equal(t, time.Date(2025, time.June, 21, 22, 21, 26, 271000000, time.UTC), user.Created)

	alice, err := queries.InsertUser(t.Context(), sqlc.InsertUserParams{
		ID:           "u_alice_admin",
		Name:         pointer.Pointer("Alice Admin"),
		Username:     "alice_admin",
		PasswordHash: "",
		TokenKey:     "",
		Created:      time.Date(2025, time.June, 21, 22, 21, 26, 0, time.UTC),
		Updated:      time.Date(2025, time.June, 21, 22, 21, 26, 0, time.UTC),
	})
	require.NoError(t, err)

	assert.Equal(t, time.Date(2025, time.June, 21, 22, 21, 26, 0, time.UTC), alice.Created)
}
