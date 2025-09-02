package ui_test

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/ui"
)

func TestUI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		wantFiles []string
	}{
		{
			name: "TestUI",
			wantFiles: []string{
				".keep",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := ui.UI()

			var gotFiles []string

			require.NoError(t, fs.WalkDir(got, ".", func(path string, d fs.DirEntry, _ error) error {
				if !d.IsDir() {
					gotFiles = append(gotFiles, path)
				}

				return nil
			}))

			for _, wantFile := range tt.wantFiles {
				assert.Contains(t, gotFiles, wantFile)
			}
		})
	}
}
