package auth

import (
	"fmt"
	"net/http"
)

func unauthorizedJSON(w http.ResponseWriter, msg string) {
	errorJSON(w, http.StatusUnauthorized, msg)
}

func errorJSON(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, _ = fmt.Fprintf(w, `{"status": %d, "error": %q, "message": %q}`, status, http.StatusText(status), msg)
}
