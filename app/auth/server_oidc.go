package auth

import "net/http"

func (s *Service) oidcLogin(w http.ResponseWriter, r *http.Request) {
	state, err := randomState()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	http.Redirect(w, r, s.oauth2Config.AuthCodeURL(state), http.StatusFound)
}

func (s *Service) oidcCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	oauth2Token, err := s.oauth2Config.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token", http.StatusInternalServerError)

		return
	}

	_, _, apiError := s.verifyClaims(r, rawIDToken)
	if apiError != nil {
		http.Error(w, apiError.Error(), http.StatusInternalServerError)

		return
	}

	// TODO: login
	// s.SessionManager.Put(r.Context(), userID)

	http.Redirect(w, r, "/", http.StatusFound)
}
