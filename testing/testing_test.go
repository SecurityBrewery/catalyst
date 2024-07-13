package testing

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tokens"

	"github.com/SecurityBrewery/catalyst/app"
	"github.com/SecurityBrewery/catalyst/migrations"
)

func baseTestApp(t *testing.T) (core.App, string, string, func()) {
	t.Helper()

	temp, err := os.MkdirTemp("", "catalyst_test_data")
	if err != nil {
		t.Fatal(err)
	}

	baseApp := app.App(temp)

	if err := app.Bootstrap(baseApp); err != nil {
		t.Fatal(err)
	}

	defaultTestData(t, baseApp)

	adminToken, err := generateAdminToken(t, baseApp, adminEmail)
	if err != nil {
		t.Fatal(err)
	}

	analystToken, err := generateRecordToken(t, baseApp, analystEmail)
	if err != nil {
		t.Fatal(err)
	}

	return baseApp, adminToken, analystToken, func() { _ = os.RemoveAll(temp) }
}

func testAppFactory(baseApp core.App) func(t *testing.T) *tests.TestApp {
	return func(t *testing.T) *tests.TestApp {
		t.Helper()

		testApp, err := tests.NewTestApp(baseApp.DataDir())
		if err != nil {
			t.Fatal(err)
		}

		return testApp
	}
}

func generateAdminToken(t *testing.T, baseApp core.App, email string) (string, error) {
	t.Helper()

	app, err := tests.NewTestApp(baseApp.DataDir())
	if err != nil {
		return "", err
	}
	defer app.Cleanup()

	admin, err := app.Dao().FindAdminByEmail(email)
	if err != nil {
		return "", err
	}

	return tokens.NewAdminAuthToken(app, admin)
}

func generateRecordToken(t *testing.T, baseApp core.App, email string) (string, error) {
	t.Helper()

	app, err := tests.NewTestApp(baseApp.DataDir())
	if err != nil {
		t.Fatal(err)
	}
	defer app.Cleanup()

	record, err := app.Dao().FindAuthRecordByEmail(migrations.UserCollectionName, email)
	if err != nil {
		return "", err
	}

	return tokens.NewRecordAuthToken(app, record)
}

type BaseTest struct {
	Name           string
	Method         string
	RequestHeaders map[string]string
	URL            string
	Body           string
	TestAppFactory func(t *testing.T) *tests.TestApp
}

type AuthBasedExpectation struct {
	Name               string
	RequestHeaders     map[string]string
	ExpectedStatus     int
	ExpectedContent    []string
	NotExpectedContent []string
	ExpectedEvents     map[string]int
}

type authMatrixText struct {
	baseTest              BaseTest
	authBasedExpectations []AuthBasedExpectation
}

func mergeScenario(base BaseTest, expectation AuthBasedExpectation) tests.ApiScenario {
	return tests.ApiScenario{
		Name:           expectation.Name,
		Method:         base.Method,
		Url:            base.URL,
		Body:           bytes.NewBufferString(base.Body),
		TestAppFactory: base.TestAppFactory,

		RequestHeaders:     mergeMaps(base.RequestHeaders, expectation.RequestHeaders),
		ExpectedStatus:     expectation.ExpectedStatus,
		ExpectedContent:    expectation.ExpectedContent,
		NotExpectedContent: expectation.NotExpectedContent,
		ExpectedEvents:     expectation.ExpectedEvents,
	}
}

func mergeMaps(a, b map[string]string) map[string]string {
	if a == nil {
		return b
	}

	if b == nil {
		return a
	}

	for k, v := range b {
		a[k] = v
	}

	return a
}

func b(data map[string]any) []byte {
	b, _ := json.Marshal(data) //nolint:errchkjson

	return b
}

func s(data map[string]any) string {
	return string(b(data))
}
