package upgradetest

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/SecurityBrewery/catalyst/app"
	"github.com/SecurityBrewery/catalyst/data"
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

			db, err := copyDatabase(t, entry.Name())
			if err != nil {
				log.Fatal(err)
			}

			pb, err := app.New(t.Context(), db)
			if err != nil {
				log.Fatal(err)
			}

			data.ValidateDefaultData(t, pb)
		})
	}
}

func copyDatabase(t *testing.T, directory string) (string, error) {
	t.Helper()

	dest := t.TempDir()

	dstDB, err := os.Create(filepath.Join(dest, "data.db"))
	if err != nil {
		return "", err
	}
	defer dstDB.Close()

	srcDB, err := os.Open(filepath.Join("data", directory, "data.db"))
	if err != nil {
		return "", err
	}
	defer srcDB.Close()

	if _, err := dstDB.ReadFrom(srcDB); err != nil {
		return "", err
	}

	return dest, nil
}
