package auth

import (
	"log"
	"net/http"

	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/argon2"

	"github.com/SecurityBrewery/catalyst/generated/time"
)

const (
	stateSessionCookie = "state"
	userSessionCookie  = "user"
)

type Jar struct {
	store *securecookie.SecureCookie
}

func NewJar(secret []byte) *Jar {
	hashSalt := securecookie.GenerateRandomKey(64)
	blockSalt := securecookie.GenerateRandomKey(64)

	return &Jar{
		store: securecookie.New(
			argon2.IDKey(secret, hashSalt, 1, 64*1024, 4, 64),
			argon2.IDKey(secret, blockSalt, 1, 64*1024, 4, 32),
		),
	}
}

func (j *Jar) setStateCookie(w http.ResponseWriter, state string) {
	encoded, err := j.store.Encode(userSessionCookie, state)
	if err != nil {
		log.Println(err)

		return
	}

	http.SetCookie(w, &http.Cookie{Name: stateSessionCookie, Value: encoded, Path: "/", Expires: time.Now().AddDate(0, 0, 1)})
}

func (j *Jar) stateCookie(r *http.Request) (string, error) {
	stateCookie, err := r.Cookie(stateSessionCookie)
	if err != nil {
		return "", err
	}

	var state string
	err = j.store.Decode(userSessionCookie, stateCookie.Value, &state)

	return state, err
}

func (j *Jar) setClaimsCookie(w http.ResponseWriter, claims map[string]any) {
	encoded, err := j.store.Encode(userSessionCookie, claims)
	if err != nil {
		log.Println(err)

		return
	}

	http.SetCookie(w, &http.Cookie{Name: userSessionCookie, Value: encoded, Path: "/", Expires: time.Now().AddDate(0, 0, 1)})
}

func deleteClaimsCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: userSessionCookie, Value: "", MaxAge: -1})
}

func (j *Jar) claimsCookie(r *http.Request) (map[string]any, bool, error) {
	userCookie, err := r.Cookie(userSessionCookie)
	if err != nil {
		return nil, true, nil
	}

	var claims map[string]any
	err = j.store.Decode(userSessionCookie, userCookie.Value, &claims)

	log.Println("claims:", claims)

	return claims, false, err
}
