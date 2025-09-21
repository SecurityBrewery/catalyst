package migration

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/upload"
)

func TestSQLMigration_UpAndDown(t *testing.T) {
	t.Parallel()

	m := sqlMigration{
		sqlName: "test_migration",
		upSQL:   "CREATE TABLE test_table (id INTEGER PRIMARY KEY, name TEXT);",
	}

	dir := t.TempDir()
	queries := database.TestDB(t, dir)
	uploader, err := upload.New(dir)
	require.NoError(t, err)

	// Test up
	require.NoError(t, m.up(t.Context(), queries, dir, uploader))

	// Table should exist
	_, err = queries.WriteDB.ExecContext(t.Context(), "INSERT INTO test_table (name) VALUES ('foo')")
	require.NoError(t, err)
}

func TestNewSQLMigration_FileNotFound(t *testing.T) {
	t.Parallel()

	f := newSQLMigration("does_not_exist")
	_, err := f()
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to read up migration file")
}
