package oidc

import (
	"context"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

type Service struct {
	oauth2Config oauth2.Config
	verifier     *oidc.IDTokenVerifier
	config       *Config
	queries      *sqlc.Queries
}

type Config struct {
	OIDCAuth     bool   `json:"oidcAuth,omitempty" yaml:"oidcAuth,omitempty"`
	OIDCIssuer   string `json:"oidcIssuer,omitempty" yaml:"oidcIssuer,omitempty"`
	ClientID     string `json:"clientID,omitempty" yaml:"clientID,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty" yaml:"clientSecret,omitempty"`
	RedirectURL  string `json:"redirectURL,omitempty" yaml:"redirectURL,omitempty"`
	AuthURL      string `json:"authURL,omitempty" yaml:"authURL,omitempty"`

	OIDCClaimUsername string `json:"oidcClaimUsername,omitempty" yaml:"oidcClaimUsername,omitempty"`
	OIDCClaimEmail    string `json:"oidcClaimEmail,omitempty" yaml:"oidcClaimEmail,omitempty"`
	OIDCClaimName     string `json:"oidcClaimName,omitempty" yaml:"oidcClaimName,omitempty"`
}

func New(ctx context.Context, queries *sqlc.Queries, config *Config) (*Service, error) {
	service := &Service{
		queries: queries,
		config:  config,
	}

	provider, err := oidc.NewProvider(ctx, config.OIDCIssuer)
	if err != nil {
		return nil, err
	}

	oauth2Config := oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	if config.AuthURL != "" {
		oauth2Config.Endpoint.AuthURL = config.AuthURL
	}

	service.oauth2Config = oauth2Config
	service.verifier = provider.Verifier(&oidc.Config{ClientID: config.ClientID, SkipIssuerCheck: true})

	return service, nil
}

func (s *Service) Login(w http.ResponseWriter, r *http.Request) {
	state, err := randomState()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	http.Redirect(w, r, s.oauth2Config.AuthCodeURL(state), http.StatusFound)
}

func (s *Service) Callback(w http.ResponseWriter, r *http.Request) {
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
