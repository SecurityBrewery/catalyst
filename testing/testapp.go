package testing

import (
	"os"
	"testing"

	"github.com/SecurityBrewery/catalyst/app2"
)

func App(t *testing.T) (*app2.App2, func()) {
	t.Helper()

	temp, err := os.MkdirTemp("", "catalyst_test_data")
	if err != nil {
		t.Fatal(err)
	}

	baseApp, err := app2.App(temp, true)
	if err != nil {
		t.Fatal(err)
	}

	defaultTestData(t, baseApp)

	return baseApp, func() { _ = os.RemoveAll(temp) }
}
