package migration

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/upload/uploader"
)

func TestApply(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	queries := database.TestDB(t, dir)
	uploader, err := uploader.New(dir)
	require.NoError(t, err)

	require.NoError(t, Apply(t.Context(), queries, dir, uploader))
}
