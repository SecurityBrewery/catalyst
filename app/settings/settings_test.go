package settings_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/app/data"
	"github.com/SecurityBrewery/catalyst/app/settings"
)

func TestUpdateSettings(t *testing.T) {
	t.Parallel()

	queries := data.NewTestDB(t, t.TempDir())

	_, err := settings.Update(t.Context(), queries, func(settings *settings.Settings) {
		settings.Meta.AppName = "UpdatedApp"
		settings.Meta.AppURL = "https://example.com"
	})
	require.NoError(t, err, "Save should not return an error")

	got, err := settings.Load(t.Context(), queries)
	require.NoError(t, err, "Load should not return an error")

	require.Equal(t, "UpdatedApp", got.Meta.AppName, "AppName should match after saving and loading settings")
	require.Equal(t, "https://example.com", got.Meta.AppURL, "AppURL should match after saving and loading settings")
}
