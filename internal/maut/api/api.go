package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPError struct {
	Status   int
	Internal error
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTPError(%d): %s", e.Status, e.Internal)
}

func (e *HTTPError) Unwrap() error {
	return e.Internal
}

func JSONError(w http.ResponseWriter, err error) {
	JSONErrorStatus(w, http.StatusInternalServerError, err)
}

func JSONErrorStatus(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	b, _ := json.Marshal(map[string]string{"error": err.Error()})
	_, _ = w.Write(b)
}
