package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/xeipuuv/gojsonschema"
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

func parseURLInt64(r *http.Request, s string) (int64, error) {
	i, err := strconv.ParseInt(chi.URLParam(r, s), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w", &HTTPError{http.StatusUnprocessableEntity, err})
	}
	return i, nil
}

func parseURLInt(r *http.Request, s string) (int, error) {
	i, err := strconv.Atoi(chi.URLParam(r, s))
	if err != nil {
		return 0, fmt.Errorf("%w", &HTTPError{http.StatusUnprocessableEntity, err})
	}
	return i, nil
}

func parseQueryInt(r *http.Request, s string) (int, error) {
	i, err := strconv.Atoi(r.URL.Query().Get(s))
	if err != nil {
		return 0, fmt.Errorf("%w", &HTTPError{http.StatusUnprocessableEntity, err})
	}
	return i, nil
}

func parseQueryBool(r *http.Request, s string) (bool, error) {
	b, err := strconv.ParseBool(r.URL.Query().Get(s))
	if err != nil {
		return false, fmt.Errorf("%w", &HTTPError{http.StatusUnprocessableEntity, err})
	}
	return b, nil
}

func parseQueryStringArray(r *http.Request, key string) ([]string, error) {
	stringArray, ok := r.URL.Query()[key]
	if !ok {
		return nil, nil
	}
	return removeEmpty(stringArray), nil
}

func removeEmpty(l []string) []string {
	var stringArray []string
	for _, s := range l {
		if s == "" {
			continue
		}
		stringArray = append(stringArray, s)
	}

	return stringArray
}

func parseQueryBoolArray(r *http.Request, key string) ([]bool, error) {
	stringArray, ok := r.URL.Query()[key]
	if !ok {
		return nil, nil
	}
	var boolArray []bool
	for _, s := range stringArray {
		if s == "" {
			continue
		}
		b, err := strconv.ParseBool(s)
		if err != nil {
			return nil, fmt.Errorf("%w", &HTTPError{http.StatusUnprocessableEntity, err})
		}
		boolArray = append(boolArray, b)
	}

	return boolArray, nil
}

func parseQueryOptionalInt(r *http.Request, key string) (*int, error) {
	s := r.URL.Query().Get(key)
	if s == "" {
		return nil, nil
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return nil, fmt.Errorf("%w", &HTTPError{http.StatusUnprocessableEntity, err})
	}
	return &i, nil
}

func parseQueryOptionalStringArray(r *http.Request, key string) ([]string, error) {
	return parseQueryStringArray(r, key)
}

func parseQueryOptionalBoolArray(r *http.Request, key string) ([]bool, error) {
	return parseQueryBoolArray(r, key)
}

func parseBody(b []byte, i interface{}) error {
	dec := json.NewDecoder(bytes.NewBuffer(b))
	err := dec.Decode(i)
	if err != nil {
		return fmt.Errorf("%w", &HTTPError{http.StatusUnprocessableEntity, err})
	}
	return nil
}

func JSONError(w http.ResponseWriter, err error) {
	JSONErrorStatus(w, http.StatusInternalServerError, err)
}

func JSONErrorStatus(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	b, _ := json.Marshal(map[string]string{"error": err.Error()})
	w.Write(b)
}

func response(w http.ResponseWriter, v interface{}, err error) {
	if err != nil {
		var httpError *HTTPError
		if errors.As(err, &httpError) {
			JSONErrorStatus(w, httpError.Status, httpError.Internal)
			return
		}
		JSONError(w, err)
		return
	}

	if v == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(v)
	w.Write(b)
}

func validateSchema(body []byte, schema *gojsonschema.Schema, w http.ResponseWriter) bool {
	jl := gojsonschema.NewBytesLoader(body)
	validationResult, err := schema.Validate(jl)
	if err != nil {
		JSONError(w, err)
		return true
	}
	if !validationResult.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)

		var validationErrors []string
		for _, valdiationError := range validationResult.Errors() {
			validationErrors = append(validationErrors, valdiationError.String())
		}

		b, _ := json.Marshal(map[string]interface{}{"error": "wrong input", "errors": validationErrors})
		w.Write(b)
		return true
	}
	return false
}

func NilMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}

func IgnoreRoles(_ []string) func(next http.Handler) http.Handler {
	return NilMiddleware()
}
