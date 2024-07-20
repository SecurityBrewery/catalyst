package testing

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	ExpectedEvents     map[string]int
}

type catalystTest struct {
	baseTest  BaseTest
	userTests []UserTest
}

func runMatrixTest(t *testing.T, baseTest BaseTest, userTest UserTest) {
	t.Helper()

	baseApp, counter, baseAppCleanup := App(t)
	defer baseAppCleanup()

	server, err := apis.InitApi(baseApp)
	require.NoError(t, err)

	if err := baseApp.OnBeforeServe().Trigger(&core.ServeEvent{
		App:    baseApp,
		Router: server,
	}); err != nil {
		t.Fatal(fmt.Errorf("failed to trigger OnBeforeServe: %w", err))
	}

	recorder := httptest.NewRecorder()
	body := bytes.NewBufferString(baseTest.Body)
	req := httptest.NewRequest(baseTest.Method, baseTest.URL, body)

	for k, v := range baseTest.RequestHeaders {
		req.Header.Set(k, v)
	}

	if userTest.AuthRecord != "" {
		token, err := generateRecordToken(t, baseApp, userTest.AuthRecord)
		require.NoError(t, err)

		req.Header.Set("Authorization", token)
	}

	if userTest.Admin != "" {
		token, err := generateAdminToken(t, baseApp, userTest.Admin)
		require.NoError(t, err)

		req.Header.Set("Authorization", token)
	}

	server.ServeHTTP(recorder, req)

	res := recorder.Result()
	defer res.Body.Close()

	assert.Equal(t, userTest.ExpectedStatus, res.StatusCode)

	for _, expectedContent := range userTest.ExpectedContent {
		assert.Contains(t, recorder.Body.String(), expectedContent)
	}

	for _, notExpectedContent := range userTest.NotExpectedContent {
		assert.NotContains(t, recorder.Body.String(), notExpectedContent)
	}

	for event, count := range userTest.ExpectedEvents {
		assert.Equal(t, count, counter.Count(event))
	}
}

func b(data map[string]any) []byte {
	b, _ := json.Marshal(data) //nolint:errchkjson

	return b
}

func s(data map[string]any) string {
	return string(b(data))
}
