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

	require.NoError(t, Apply(t.Context(), db, nil, "", nil))
}
