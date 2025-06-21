package migration

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestSQLMigration_UpAndDown(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	m := sqlMigration{
		sqlName: "test_migration",
		upSQL:   "CREATE TABLE test_table (id INTEGER PRIMARY KEY, name TEXT);",
		downSQL: "DROP TABLE test_table;",
	}

	// Test up
	require.NoError(t, m.up(t.Context(), db))
	// Table should exist
	_, err = db.Exec("INSERT INTO test_table (name) VALUES ('foo')")
	require.NoError(t, err)

	// Test down
	require.NoError(t, m.down(t.Context(), db))
	// Table should not exist
	_, err = db.Exec("INSERT INTO test_table (name) VALUES ('bar')")
	require.Error(t, err)
}

func TestNewSQLMigration_FileNotFound(t *testing.T) {
	t.Parallel()

	f := newSQLMigration("does_not_exist")
	_, err := f()
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to read up migration file")
}
