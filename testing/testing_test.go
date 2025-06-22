package testing

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type baseTest struct {
	Name           string
	Method         string
	RequestHeaders map[string]string
	URL            string
	Body           string
}

type userTest struct {
	Name               string
	AuthRecord         string
	Admin              string
	ExpectedStatus     int
	ExpectedHeaders    map[string]string
	ExpectedContent    []string
	NotExpectedContent []string
	ExpectedEvents     map[string]int
}

type catalystTest struct {
	baseTest  baseTest
	userTests []userTest
}

func runMatrixTest(t *testing.T, baseTest baseTest, userTest userTest) {
	t.Helper()

	baseApp, cleanup, counter := App(t)

	t.Cleanup(cleanup)

	recorder := httptest.NewRecorder()
	body := bytes.NewBufferString(baseTest.Body)
	req := httptest.NewRequest(baseTest.Method, baseTest.URL, body)

	for k, v := range baseTest.RequestHeaders {
		req.Header.Set(k, v)
	}

	if userTest.AuthRecord != "" {
		user, err := baseApp.Queries.UserByEmail(t.Context(), userTest.AuthRecord)
		require.NoError(t, err)

		permissions, err := baseApp.Queries.ListUserPermissions(t.Context(), user.ID)
		require.NoError(t, err)

		loginToken, err := baseApp.Auth.CreateAccessToken(t.Context(), &user, permissions, time.Hour)
		require.NoError(t, err)

		req.Header.Set("Authorization", "Bearer "+loginToken)
	}

	if userTest.Admin != "" {
		user, err := baseApp.Queries.UserByEmail(t.Context(), userTest.Admin)
		require.NoError(t, err)

		permissions, err := baseApp.Queries.ListUserPermissions(t.Context(), user.ID)
		require.NoError(t, err)

		loginToken, err := baseApp.Auth.CreateAccessToken(t.Context(), &user, permissions, time.Hour)
		require.NoError(t, err)

		req.Header.Set("Authorization", "Bearer "+loginToken)
	}

	baseApp.Router.ServeHTTP(recorder, req)

	res := recorder.Result()
	defer res.Body.Close()

	assert.Equal(t, userTest.ExpectedStatus, res.StatusCode)

	for k, v := range userTest.ExpectedHeaders {
		assert.Equal(t, v, res.Header.Get(k))
	}

	for _, expectedContent := range userTest.ExpectedContent {
		assert.Contains(t, recorder.Body.String(), expectedContent)
	}

	for _, notExpectedContent := range userTest.NotExpectedContent {
		assert.NotContains(t, recorder.Body.String(), notExpectedContent)
	}

	for event, count := range userTest.ExpectedEvents {
		assert.Equalf(t, count, counter.Count(event), "expected %d events for %s, got %d", count, event, counter.Count(event))
	}
}

func b(data map[string]any) []byte {
	b, _ := json.Marshal(data) //nolint:errchkjson

	return b
}

func s(data map[string]any) string {
	return string(b(data))
}
