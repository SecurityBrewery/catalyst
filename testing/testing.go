package testing

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

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

	baseApp, baseAppCleanup := App(t)
	defer baseAppCleanup()

	recorder := httptest.NewRecorder()
	body := bytes.NewBufferString(baseTest.Body)
	req := httptest.NewRequest(baseTest.Method, baseTest.URL, body)

	for k, v := range baseTest.RequestHeaders {
		req.Header.Set(k, v)
	}

	if userTest.AuthRecord != "" {
		req.Header.Set("Authorization", userTest.AuthRecord)
	}

	if userTest.Admin != "" {
		req.Header.Set("Authorization", userTest.Admin)
	}

	err := baseApp.Server(t.Context())
	require.NoError(t, err)

	baseApp.Router.ServeHTTP(recorder, req)

	res := recorder.Result()
	defer res.Body.Close()

	assert.Equal(t, userTest.ExpectedStatus, res.StatusCode)

	for _, expectedContent := range userTest.ExpectedContent {
		assert.Contains(t, recorder.Body.String(), expectedContent)
	}

	for _, notExpectedContent := range userTest.NotExpectedContent {
		assert.NotContains(t, recorder.Body.String(), notExpectedContent)
	}

	// for event, count := range userTest.ExpectedEvents { // TODO: add event counting
	// 	assert.Equal(t, count, counter.Count(event))
	// }
}

func b(data map[string]any) []byte {
	b, _ := json.Marshal(data) //nolint:errchkjson

	return b
}

func s(data map[string]any) string {
	return string(b(data))
}
