package service

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/app/data"
	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/hook"
	"github.com/SecurityBrewery/catalyst/app/migration"
	"github.com/SecurityBrewery/catalyst/app/openapi"
	"github.com/SecurityBrewery/catalyst/app/upload/uploader"
)

func newTestService(t *testing.T) *Service {
	t.Helper()

	dir := t.TempDir()
	queries := data.NewTestDB(t, dir)
	hooks := hook.NewHooks()
	uploader, err := uploader.New(dir)
	require.NoError(t, err)

	err = migration.Apply(t.Context(), queries, dir, uploader)
	require.NoError(t, err)

	return New(queries, hooks, uploader, nil)
}

func Test_toString(t *testing.T) {
	t.Parallel()

	str := "hello"

	assert.Equal(t, "hello", toString(&str, "default"))
	assert.Equal(t, "default", toString(nil, "default"))
}

func Test_toInt64(t *testing.T) {
	t.Parallel()

	i := 42
	assert.Equal(t, int64(42), toInt64(&i, 0))
	assert.Equal(t, int64(7), toInt64(nil, 7))
}

func Test_marshal_unmarshal(t *testing.T) {
	t.Parallel()

	m := map[string]any{"key": "value"}
	j := marshal(m)
	assert.JSONEq(t, `{"key":"value"}`, string(j))

	decoded := unmarshal(j)
	require.NotNil(t, decoded)
	assert.Equal(t, "value", decoded["key"])

	assert.Nil(t, unmarshal(json.RawMessage("{invalid}")))
}

func Test_marshalPointer(t *testing.T) {
	t.Parallel()

	m := map[string]any{"x": 1}
	assert.JSONEq(t, `{"x":1}`, string(marshalPointer(&m)))
	assert.JSONEq(t, "{}", string(marshalPointer(nil)))
}

func Test_generateID(t *testing.T) {
	t.Parallel()

	id1 := database.GenerateID("p")

	id2 := database.GenerateID("p")

	if id1 == id2 {
		t.Errorf("expected unique ids, got %s and %s", id1, id2)
	}

	if len(id1) == 0 || len(id2) == 0 {
		t.Error("expected non empty ids")
	}

	if id1[:1] != "p" || id2[:1] != "p" {
		t.Errorf("expected ids to start with prefix")
	}
}

func TestService_DownloadFile(t *testing.T) {
	t.Parallel()

	s := newTestService(t)

	resp, err := s.DownloadFile(t.Context(), openapi.DownloadFileRequestObject{Id: "b_test_file"})
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

	s := newTestService(t)

	// invalid format
	_, err := s.queries.CreateFile(t.Context(), sqlc.CreateFileParams{Name: "bad", Blob: "invalid", Size: 1, Ticket: "test-ticket"})
	require.NoError(t, err)

	_, err = s.DownloadFile(t.Context(), openapi.DownloadFileRequestObject{Id: "f_invalid_format"})
	require.Error(t, err)

	// invalid base64
	_, err = s.queries.CreateFile(t.Context(), sqlc.CreateFileParams{Name: "bad", Blob: "data:text/plain;base64,???", Size: 1, Ticket: "test-ticket"})
	require.NoError(t, err)

	_, err = s.DownloadFile(t.Context(), openapi.DownloadFileRequestObject{Id: "f_invalid_base64"})
	require.Error(t, err)
}
