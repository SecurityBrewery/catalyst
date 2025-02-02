package upgradetest

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/SecurityBrewery/catalyst/app"
	"github.com/SecurityBrewery/catalyst/fakedata"
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

			pb, err := app.App(filepath.Join("data", entry.Name()), true)
			if err != nil {
				log.Fatal(err)
			}

			if err := pb.Bootstrap(); err != nil {
				t.Fatal(fmt.Errorf("failed to bootstrap: %w", err))
			}

			if err := fakedata.ValidateDefaultData(pb); err != nil {
				log.Fatal(err)
			}
		})
	}
}
