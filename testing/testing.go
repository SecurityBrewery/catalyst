package testing

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type BaseTest struct {
	Name           string
	Method         string
	RequestHeaders map[string]string
	URL            string
	Body           string
}

type UserTest struct {
	Name               string
	AuthRecord         string
	Admin              string
	ExpectedStatus     int
	ExpectedContent    []string
	NotExpectedContent []string
	ExpectedEvents     []string
}

type catalystTest struct {
	baseTest  BaseTest
	userTests []UserTest
}

func runMatrixTest(t *testing.T, baseTest BaseTest, userTest UserTest) {
	t.Helper()

	baseApp, counter, baseAppCleanup := App(t)
	defer baseAppCleanup()

	server, err := apis.NewRouter(baseApp)
	require.NoErrorf(t, err, "failed to create router: %v", err)

	err = baseApp.OnServe().Trigger(&core.ServeEvent{
		App:    baseApp,
		Router: server,
	}, func(event *core.ServeEvent) error {
		recorder := httptest.NewRecorder()
		body := bytes.NewBufferString(baseTest.Body)
		req := httptest.NewRequest(baseTest.Method, baseTest.URL, body)

		for k, v := range baseTest.RequestHeaders {
			req.Header.Set(k, v)
		}

		if userTest.AuthRecord != "" {
			token, err := generateRecordToken(t, baseApp, userTest.AuthRecord)
			require.NoErrorf(t, err, "failed to generate record token: %v", err)

			req.Header.Set("Authorization", token)
		}

		if userTest.Admin != "" {
			token, err := generateAdminToken(t, baseApp, userTest.Admin)
			require.NoErrorf(t, err, "failed to generate admin token: %v", err)

			req.Header.Set("Authorization", token)
		}

		mux, err := event.Router.BuildMux()
		require.NoErrorf(t, err, "failed to build mux: %v", err)

		mux.ServeHTTP(recorder, req)

		res := recorder.Result()
		defer res.Body.Close()

		assert.Equal(t, userTest.ExpectedStatus, res.StatusCode)

		for _, expectedContent := range userTest.ExpectedContent {
			assert.Containsf(t, recorder.Body.String(), expectedContent, "expected content not found: %s", expectedContent)
		}

		for _, notExpectedContent := range userTest.NotExpectedContent {
			assert.NotContainsf(t, recorder.Body.String(), notExpectedContent, "unexpected content found: %s", notExpectedContent)
		}

		var seenEvents []string

		for name, count := range counter.Counts() {
			if count > 0 {
				seenEvents = append(seenEvents, name)
			}
		}

		assert.ElementsMatchf(t, userTest.ExpectedEvents, seenEvents, "expected events not found: %v", userTest.ExpectedEvents)

		return nil
	})

	require.NoErrorf(t, err, "failed to trigger serve event: %v", err)
}

func b(data map[string]any) []byte {
	b, _ := json.Marshal(data) //nolint:errchkjson

	return b
}

func s(data map[string]any) string {
	return string(b(data))
}
