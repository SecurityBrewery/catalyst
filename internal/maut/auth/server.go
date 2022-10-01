package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/cugu/maut/api"
)

func (a *Authenticator) Server() *chi.Mux {
	server := chi.NewRouter()

	server.Get("/config", a.authConfig())

	if a.config.OIDCAuthEnable {
		server.Get("/callback", a.Callback())
		server.Get("/oidclogin", a.redirectToOIDCLogin())
	}
	if a.config.SimpleAuthEnable {
		server.Post("/login", a.login())
	}
	server.Post("/logout", a.logout())

	return server
}

func (a *Authenticator) authConfig() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		b, _ := json.Marshal(map[string]any{
			"simple": a.config.SimpleAuthEnable,
			"oidc":   a.config.OIDCAuthEnable,
		})

		_, _ = writer.Write(b)
	}
}

func (a *Authenticator) Callback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state, isNew := a.jar.stateSession(r)
		if isNew || state == "" {
			api.JSONErrorStatus(w, http.StatusInternalServerError, errors.New("state missing"))

			return
		}

		if state != r.URL.Query().Get("state") {
			api.JSONErrorStatus(w, http.StatusInternalServerError, errors.New("state mismatch"))

			return
		}

		oauth2Token, err := a.config.OAuth2.Exchange(r.Context(), r.URL.Query().Get("code"))
		if err != nil {
			api.JSONErrorStatus(w, http.StatusInternalServerError, fmt.Errorf("oauth2 exchange failed: %w", err))

			return
		}

		// Extract the ID Token from OAuth2 token.
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			api.JSONErrorStatus(w, http.StatusInternalServerError, errors.New("missing id token"))

			return
		}

		username, _, apiError := a.verifyClaims(r, rawIDToken)
		if apiError != nil {
			api.JSONErrorStatus(w, apiError.Status, apiError.Internal)

			return
		}

		a.jar.setUserSession(r, w, username)
		a.jar.deleteStateSession(r, w)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (a *Authenticator) redirectToOIDCLogin() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		state, err := state()
		if err != nil {
			api.JSONErrorStatus(w, http.StatusInternalServerError, errors.New("generating state failed"))

			return
		}

		a.jar.setStateSession(r, w, state)

		http.Redirect(w, r, a.config.OAuth2.AuthCodeURL(state), http.StatusFound)
	}
}

func (a *Authenticator) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type credentials struct {
			Username string
			Password string
		}
		cr := credentials{}

		if err := json.NewDecoder(r.Body).Decode(&cr); err != nil {
			api.JSONErrorStatus(w, http.StatusUnauthorized, errors.New("wrong username or password"))

			return
		}

		user, err := a.resolver.UserByIDAndPassword(r.Context(), cr.Username, cr.Password)
		if err != nil {
			api.JSONErrorStatus(w, http.StatusUnauthorized, errors.New("wrong username or password"))

			return
		}

		a.jar.setUserSession(r, w, user.ID)

		b, _ := json.Marshal(map[string]string{"login": "successful"})
		_, _ = w.Write(b)
	}
}

func (a *Authenticator) logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.jar.deleteUserSession(r, w)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
