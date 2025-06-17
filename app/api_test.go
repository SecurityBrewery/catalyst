package app_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"

	catalystTesting "github.com/SecurityBrewery/catalyst/app/testing"
)

type request struct {
	Method string
	URL    string
	Body   []byte
}

func TestAPI(t *testing.T) { //nolint:cyclop
	t.Parallel()

	requests := [][]request{
		{
			{"GET", "/ui/login", nil},
			{"GET", "/api/config", nil},
		},
		{
			{"POST", "/auth/local/login", []byte(`{"email":"admin@catalyst-soar.com","password":"password123"}`)},
		},
		{
			{"GET", "/auth/user", nil},
			{"GET", "/api/sidebar", nil},
			{"GET", "/ui/groups", nil},
			{"GET", "/api/tickets", nil},
			{"GET", "/api/tickets", nil},
			{"GET", "/api/tasks", nil},
			{"GET", "/auth/user", nil},
			{"GET", "/auth/user", nil},
			{"GET", "/api/sidebar", nil},
			{"GET", "/api/groups", nil},
			{"GET", "/api/config", nil},
		},
		{
			{"POST", "/api/groups", []byte(`{"name":"playwright-59537c9d-772f-45f0-a8f0-bb99656fa7ed","permissions":["ticket:read"]}`)},
		},
		{
			{"GET", "/api/groups/ID", nil},
			{"GET", "/api/groups", nil},
			{"GET", "/api/config", nil},
		},
		{
			{"GET", "/api/groups/ID/parents", nil},
			{"GET", "/api/groups/ID/permissions", nil},
			{"GET", "/api/groups/ID/children", nil},
			{"GET", "/api/groups/ID/users", nil},
			{"GET", "/api/users", nil},
		},
		{
			{"GET", "/ui/login", nil},
			{"GET", "/api/config", nil},
		},
		{
			{"POST", "/auth/local/login", []byte(`{"email":"admin@catalyst-soar.com","password":"password123"}`)},
		},
		{
			{"GET", "/auth/user", nil},
			{"GET", "/api/sidebar", nil},
			{"GET", "/ui/groups", nil},
			{"GET", "/api/tickets", nil},
			{"GET", "/api/tickets", nil},
			{"GET", "/api/tasks", nil},
			{"GET", "/auth/user", nil},
			{"GET", "/auth/user", nil},
			{"GET", "/api/sidebar", nil},
			{"GET", "/api/groups", nil},
			{"GET", "/api/config", nil},
		},
		{
			{"POST", "/api/groups", []byte(`{"name":"playwright-c9bdcbf8-6aba-4974-9c12-f77fa7f15b34","permissions":["ticket:read"]}`)},
		},
	}

	app, cleanup, _ := catalystTesting.App(t)
	defer cleanup()

	id, token := "", ""

	for _, batch := range requests {
		var wg sync.WaitGroup

		wg.Add(len(batch))

		for _, req := range batch {
			go func() {
				defer wg.Done()

				url := req.URL
				if strings.Contains(url, "ID") {
					url = strings.ReplaceAll(url, "ID", id)
				}

				t.Logf("%s %s", req.Method, url)

				ctx, cancel := context.WithTimeout(t.Context(), time.Second*10)
				defer cancel()

				r, err := http.NewRequestWithContext(ctx, req.Method, url, bytes.NewReader(req.Body))
				assert.NoError(t, err)

				if token != "" {
					r.Header.Set("Authorization", "Bearer "+token)
				}

				w := httptest.NewRecorder()

				app.Router.ServeHTTP(w, r)

				assert.Equal(t, 200, w.Code, "expected status code 200: %s for %s %s", w.Body.String(), req.Method, url)

				if w.Code == 200 && req.Method == http.MethodPost && gjson.Get(w.Body.String(), "id").Exists() {
					id = gjson.Get(w.Body.String(), "id").String()
				}

				if w.Code == 200 && req.Method == http.MethodPost && gjson.Get(w.Body.String(), "token").Exists() {
					token = gjson.Get(w.Body.String(), "token").String()
				}
			}()
		}

		wg.Wait()
	}
}
