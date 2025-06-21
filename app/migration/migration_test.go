package migration

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestApply(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	require.NoError(t, Apply(t.Context(), db))
}

func TestRollback_Success(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	require.NoError(t, Apply(t.Context(), db))

	require.NoError(t, Rollback(t.Context(), db))
}

func TestRollback_WrongVersion(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	// Set an incorrect version to simulate a rollback error
	err = setVersion(t.Context(), db, 999)
	require.NoError(t, err)

	err = Rollback(t.Context(), db)
	require.Error(t, err)
	require.Contains(t, err.Error(), "current version 999 exceeds available migrations 4")
}

func TestRollback_NoMigrations(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	_ = setVersion(t.Context(), db, 0)

	err = Rollback(t.Context(), db)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no migrations to rollback")
}
