package migration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMigrations_Success(t *testing.T) {
	t.Parallel()

	migs, err := migrations(0)
	require.NoError(t, err)
	require.Len(t, migs, len(migrationGenerators))
}

func TestMigrations_VersionOffset(t *testing.T) {
	t.Parallel()

	migs, err := migrations(1)
	require.NoError(t, err)
	require.Len(t, migs, len(migrationGenerators)-1)
}

func TestMigrations_Error(t *testing.T) {
	t.Parallel()

	migs, err := migrations(999) // Invalid version
	require.Error(t, err)
	require.Nil(t, migs)
	require.Contains(t, err.Error(), "invalid migration version: 999")
}
