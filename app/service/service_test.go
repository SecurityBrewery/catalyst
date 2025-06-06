package service

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/hook"
	"github.com/SecurityBrewery/catalyst/app/openapi"
)

func newTestService(t *testing.T) *Service {
	t.Helper()

	queries := database.NewTestDB(t)
	hooks := hook.NewHooks()

	return New(queries, hooks, nil)
}

func Test_toString(t *testing.T) {
	t.Parallel()

	str := "hello"

	assert.Equal(t, "hello", toString(&str, "default"))
	assert.Equal(t, "default", toString(nil, "default"))
}

func Test_toNullString(t *testing.T) {
	t.Parallel()

	s := "value"
	assert.True(t, toNullString(&s).Valid)
	assert.False(t, toNullString(nil).Valid)
}

func Test_toInt64(t *testing.T) {
	t.Parallel()

	i := 42
	assert.Equal(t, int64(42), toInt64(&i, 0))
	assert.Equal(t, int64(7), toInt64(nil, 7))
}

func Test_toNullBool(t *testing.T) {
	t.Parallel()

	b := true
	assert.True(t, toNullBool(&b).Valid)
	assert.False(t, toNullBool(nil).Valid)
}

func Test_marshal_unmarshal(t *testing.T) {
	t.Parallel()

	m := map[string]interface{}{"key": "value"}
	json := marshal(m)
	assert.JSONEq(t, `{"key":"value"}`, json)

	decoded := unmarshal(json)
	require.NotNil(t, decoded)
	assert.Equal(t, "value", decoded["key"])

	assert.Nil(t, unmarshal("{invalid}"))
}

func Test_marshalPointer(t *testing.T) {
	t.Parallel()

	m := map[string]interface{}{"x": 1}
	assert.Equal(t, `{"x":1}`, marshalPointer(&m))
	assert.Equal(t, "{}", marshalPointer(nil))
}

func Test_generateID(t *testing.T) {
	t.Parallel()

	id := generateID("test")
	assert.Greater(t, len(id), len("test")+1)
	assert.True(t, strings.HasPrefix(id, "test-"))
}

func TestService_DownloadFile(t *testing.T) {
	t.Parallel()

	s := newTestService(t)

	resp, err := s.DownloadFile(t.Context(), openapi.DownloadFileRequestObject{Id: "f_test_file"})
	require.NoError(t, err)

	download, ok := resp.(openapi.DownloadFile200ApplicationoctetStreamResponse)
	require.True(t, ok)

	data, err := io.ReadAll(download.Body)
	require.NoError(t, err)

	assert.Equal(t, "hello", string(data))
	assert.Equal(t, int64(len("hello")), download.ContentLength)
	assert.Equal(t, "attachment; filename=\"hello.txt\"", download.Headers.ContentDisposition)
	assert.Equal(t, "text/plain", download.Headers.ContentType)
}

func TestService_DownloadFile_Errors(t *testing.T) {
	t.Parallel()

	queries := database.NewTestDB(t)
	hooks := hook.NewHooks()
	s := New(queries, hooks, nil)

	// invalid format
	_, err := queries.CreateFile(t.Context(), sqlc.CreateFileParams{ID: "f_invalid_format", Name: "bad", Blob: "invalid", Size: 1, Ticket: "test-ticket"})
	require.NoError(t, err)

	_, err = s.DownloadFile(t.Context(), openapi.DownloadFileRequestObject{Id: "f_invalid_format"})
	require.Error(t, err)

	// invalid base64
	_, err = queries.CreateFile(t.Context(), sqlc.CreateFileParams{ID: "f_invalid_base64", Name: "bad", Blob: "data:text/plain;base64,???", Size: 1, Ticket: "test-ticket"})
	require.NoError(t, err)

	_, err = s.DownloadFile(t.Context(), openapi.DownloadFileRequestObject{Id: "f_invalid_base64"})
	require.Error(t, err)
}
