package catalyst

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

const (
	stateSessionCookie = "state"
	userSessionCookie  = "user"
)

func setStateCookie(w http.ResponseWriter, state string) {
	http.SetCookie(w, &http.Cookie{Name: stateSessionCookie, Value: state})
}

func stateCookie(r *http.Request) (string, error) {
	stateCookie, err := r.Cookie(stateSessionCookie)
	if err != nil {
		return "", err
	}

	return stateCookie.Value, nil
}

func setClaimsCookie(w http.ResponseWriter, claims map[string]any) {
	b, err := json.Marshal(claims)
	if err != nil {
		log.Println(err)

		return
	}
	http.SetCookie(w, &http.Cookie{Name: userSessionCookie, Value: base64.StdEncoding.EncodeToString(b)})
}

func deleteClaimsCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: userSessionCookie, Value: "", MaxAge: -1})
}

func claimsCookie(r *http.Request) (map[string]any, bool, error) {
	userCookie, err := r.Cookie(userSessionCookie)
	if err != nil {
		return nil, true, nil
	}

	b, err := base64.StdEncoding.DecodeString(userCookie.Value)
	if err != nil {
		return nil, false, fmt.Errorf("could not decode cookie: %w", err)
	}

	var claims map[string]any
	if err := json.Unmarshal(b, &claims); err != nil {
		return nil, false, errors.New("claims not in session")
	}

	return claims, false, err
}
