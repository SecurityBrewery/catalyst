package ui

import "testing"

func TestUI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
	}{
		{"index.html", "dist/index.html"},
		{"favicon.ico", "dist/favicon.ico"},
		{"manifest.json", "dist/manifest.json"},
		{"img", "dist/img"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f, err := UI.Open(tt.path)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
		})
	}
}
