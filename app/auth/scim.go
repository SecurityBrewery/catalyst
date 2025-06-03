package auth

import (
	"fmt"
	"net/http"
)

func scimUnauthorized(w http.ResponseWriter, detail string) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Set("Content-Type", "application/scim+json; charset=UTF-8")
	_, _ = fmt.Fprintf(w, `{"schemas":["urn:ietf:params:scim:api:messages:2.0:Error"],"detail":%q,"status":"401"}`, detail)
}

func scimError(w http.ResponseWriter, status int, detail string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/scim+json; charset=UTF-8")
	_, _ = fmt.Fprintf(w, `{"schemas":["urn:ietf:params:scim:api:messages:2.0:Error"],"detail":%q,"status":"%d"}`, detail, status)
}
