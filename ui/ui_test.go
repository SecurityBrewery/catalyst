package ui

import "testing"

func TestUI(t *testing.T) {
	requiredFiles := []string{
		"dist/index.html",
		"dist/favicon.ico",
		"dist/manifest.json",
		"dist/img",
	}
	for _, requiredFile := range requiredFiles {
		t.Run("Require "+requiredFile, func(t *testing.T) {
			f, err := UI.Open(requiredFile)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
		})
	}
}
