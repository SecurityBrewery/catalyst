package auth

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	stateSessionSession = "maut_state"
	userSessionSession  = "maut_user"
)

type Jar struct {
	store sessions.Store
}

func NewJar(secret []byte) *Jar {
	return &Jar{
		store: sessions.NewCookieStore(secret),
	}
}

func (j *Jar) setStateSession(r *http.Request, w http.ResponseWriter, state string) {
	session := sessions.NewSession(j.store, stateSessionSession)
	session.Values["state"] = state
	session.Options.MaxAge = 60 * 5 // 5 minutes
	session.Options.Path = "/"
	if err := j.store.Save(r, w, session); err != nil {
		log.Println(err)

		return
	}
}

func (j *Jar) deleteStateSession(r *http.Request, w http.ResponseWriter) {
	session, _ := j.store.Get(r, stateSessionSession)
	session.Options.MaxAge = -1
	err := j.store.Save(r, w, session)
	if err != nil {
		log.Println(err)
	}
}

func (j *Jar) stateSession(r *http.Request) (state string, isNew bool) {
	session, _ := j.store.Get(r, stateSessionSession)
	if session.IsNew {
		return "", true
	}

	if state, ok := session.Values["state"]; ok {
		if state, ok := state.(string); ok {
			return state, false
		}
	}

	return "", true
}

func (j *Jar) setUserSession(r *http.Request, w http.ResponseWriter, userID string) {
	session := sessions.NewSession(j.store, userSessionSession)
	session.Values["id"] = userID
	session.Options.MaxAge = 60 * 60 * 24 // 24 hours
	session.Options.Path = "/"
	if err := j.store.Save(r, w, session); err != nil {
		log.Println(err)

		return
	}
}

func (j *Jar) deleteUserSession(r *http.Request, w http.ResponseWriter) {
	session, _ := j.store.Get(r, userSessionSession)
	session.Options.MaxAge = -1
	err := j.store.Save(r, w, session)
	if err != nil {
		log.Println(err)
	}
}

func (j *Jar) userSession(r *http.Request) (userID string, isNew bool) {
	session, _ := j.store.Get(r, userSessionSession)
	if session.IsNew {
		return "", true
	}

	if userID, ok := session.Values["id"]; ok {
		if userID, ok := userID.(string); ok {
			return userID, false
		}
	}

	return "", true
}
