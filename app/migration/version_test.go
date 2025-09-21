package migration

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestVersionAndSetVersion(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err, "failed to open in-memory db")

	defer db.Close()

	ver, err := version(t.Context(), db)
	require.NoError(t, err, "failed to get version")
	require.Equal(t, 0, ver, "expected version 0")

	err = setVersion(t.Context(), db, 2)
	require.NoError(t, err, "failed to set version")

	ver, err = version(t.Context(), db)
	require.NoError(t, err, "failed to get version after set")
	require.Equal(t, 2, ver, "expected version 2")
}
