package upgradetest

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/SecurityBrewery/catalyst/app2"
	"github.com/SecurityBrewery/catalyst/app2/fakedata"
)

func TestUpgrades(t *testing.T) {
	t.Parallel()

	dirEntries, err := os.ReadDir("data")
	if err != nil {
		t.Fatal(err)
	}

	for _, entry := range dirEntries {
		if !entry.IsDir() {
			continue
		}

		t.Run(entry.Name(), func(t *testing.T) {
			t.Parallel()

			pb, err := app2.App(t.Context(), filepath.Join("data", entry.Name()), true)
			if err != nil {
				log.Fatal(err)
			}

			if err := fakedata.ValidateDefaultData(t.Context(), t, pb); err != nil {
				log.Fatal(err)
			}
		})
	}
}
