package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/stretchr/testify/require"
)

func Test_staticFiles(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	e := &core.RequestEvent{}
	e.Request = req
	e.Response = rec

	require.NoError(t, staticFiles()(e))
}
