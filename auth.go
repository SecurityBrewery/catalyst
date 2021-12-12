package catalyst

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/models"
	"github.com/SecurityBrewery/catalyst/hooks"
	"github.com/SecurityBrewery/catalyst/role"
)

type AuthConfig struct {
	OIDCIssuer string
	OAuth2     *oauth2.Config

	OIDCClaimUsername string
	OIDCClaimEmail    string
	// OIDCClaimGroups   string
	OIDCClaimName    string
	AuthBlockNew     bool
	AuthDefaultRoles []role.Role

	provider *oidc.Provider
}

func (c *AuthConfig) Verifier(ctx context.Context) (*oidc.IDTokenVerifier, error) {
	if c.provider == nil {
		err := c.Load(ctx)
		if err != nil {
			return nil, err
		}
	}
	return c.provider.Verifier(&oidc.Config{SkipClientIDCheck: true}), nil
}

func (c *AuthConfig) Load(ctx context.Context) error {
	provider, err := oidc.NewProvider(ctx, c.OIDCIssuer)
	if err != nil {
		return err
	}
	c.provider = provider
	c.OAuth2.Endpoint = provider.Endpoint()

	return nil
}

const (
	SessionName  = "catalyst-session"
	stateSession = "state"
	userSession  = "user"
)

func Authenticate(db *database.Database, config *AuthConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		iss := config.OIDCIssuer

		keyHeader := ctx.Request.Header.Get("PRIVATE-TOKEN")
		if keyHeader != "" {
			keyAuth(db, keyHeader)(ctx)
			return
		}

		authHeader := ctx.Request.Header.Get("User")

		if authHeader != "" {
			bearerAuth(db, authHeader, iss, config)(ctx)
			return
		}
		sessionAuth(db, config)(ctx)
	}
}

func oidcCtx(ctx *gin.Context) (context.Context, context.CancelFunc) {
	/*
		if config.TLSCertFile != "" && config.TLSKeyFile != "" {
			cert, err := tls.LoadX509KeyPair(config.TLSCertFile, config.TLSKeyFile)
			if err != nil {
				return nil, err
			}

			rootCAs, _ := x509.SystemCertPool()
			if rootCAs == nil {
				rootCAs = x509.NewCertPool()
			}
			for _, c := range cert.Certificate {
				rootCAs.AppendCertsFromPEM(c)
			}

			return oidc.ClientContext(ctx, &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						RootCAs:            rootCAs,
						InsecureSkipVerify: true,
					},
				},
			}), nil
		}
	*/
	cctx, cancel := context.WithTimeout(ctx, time.Minute)
	return cctx, cancel
}

func bearerAuth(db *database.Database, authHeader string, iss string, config *AuthConfig) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no bearer token"})
			return
		}

		oidcCtx, cancel := oidcCtx(ctx)
		defer cancel()

		verifier, err := config.Verifier(oidcCtx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "could not verify: " + err.Error()})
			return
		}
		authToken, err := verifier.Verify(oidcCtx, authHeader[7:])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not verify bearer token: %v", err)})
			return
		}

		var claims map[string]interface{}
		if err := authToken.Claims(&claims); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to parse claims: %v", err)})
			return
		}

		// if claims.Iss != iss {
		// 	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "wrong issuer"})
		// 	return
		// }

		session := sessions.Default(ctx)
		session.Set(userSession, claims)
		if err = session.Save(); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprintf("could not set session: %v", err))
			return
		}

		if err = setContextClaims(ctx, db, claims, config); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not load user: %s", err)})
			return
		}
		ctx.Next()
	}
}

func keyAuth(db *database.Database, keyHeader string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		h := fmt.Sprintf("%x", sha256.Sum256([]byte(keyHeader)))

		key, err := db.UserByHash(ctx, h)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not verify private token: %v", err)})
			return
		}

		setContextUser(ctx, key, db.Hooks)

		ctx.Next()
	}
}

func sessionAuth(db *database.Database, config *AuthConfig) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		user := session.Get(userSession)
		if user == nil {
			redirectToLogin(ctx, session, config.OAuth2)

			return
		}

		claims, ok := user.(map[string]interface{})
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "claims not in session"})
			return
		}

		if err := setContextClaims(ctx, db, claims, config); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not load user: %s", err)})
			return
		}

		ctx.Next()
	}
}

func setContextClaims(ctx *gin.Context, db *database.Database, claims map[string]interface{}, config *AuthConfig) error {
	newUser, newSetting, err := mapUserAndSettings(claims, config)
	if err != nil {
		return err
	}

	if _, ok := busdb.UserFromContext(ctx); !ok {
		busdb.SetContext(ctx, &models.UserResponse{ID: "auth", Roles: []string{role.Admin}, Apikey: false, Blocked: false})
	}

	user, err := db.UserGetOrCreate(ctx, newUser)
	if err != nil {
		return err
	}

	_, err = db.UserDataGetOrCreate(ctx, newUser.ID, newSetting)
	if err != nil {
		return err
	}

	setContextUser(ctx, user, db.Hooks)
	return nil
}

func setContextUser(ctx *gin.Context, user *models.UserResponse, hooks *hooks.Hooks) {
	groups, err := hooks.GetGroups(ctx, user.ID)
	if err == nil {
		busdb.SetGroupContext(ctx, groups)
	}

	busdb.SetContext(ctx, user)
}

func mapUserAndSettings(claims map[string]interface{}, config *AuthConfig) (*models.UserForm, *models.UserData, error) {
	// handle Bearer tokens
	// if typ, ok := claims["typ"]; ok && typ == "Bearer" {
	// 	return &models.User{
	// 		Username: "bot",
	// 		Blocked:  false,
	// 		Email:    pointer.String("bot@example.org"),
	// 		Roles:    []string{"user:read", "settings:read", "ticket", "backup:read", "backup:restore"},
	// 		Name:     pointer.String("Bot"),
	// 	}, nil
	// }

	username, err := getString(claims, config.OIDCClaimUsername)
	if err != nil {
		return nil, nil, err
	}
	email, err := getString(claims, config.OIDCClaimEmail)
	if err != nil {
		email = ""
	}

	name, err := getString(claims, config.OIDCClaimName)
	if err != nil {
		name = ""
	}

	return &models.UserForm{
			ID:      username,
			Blocked: config.AuthBlockNew,
			Roles:   role.Strings(config.AuthDefaultRoles),
		}, &models.UserData{
			Email: &email,
			Name:  &name,
		}, nil
}

func getString(m map[string]interface{}, key string) (string, error) {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s, nil
		}
		return "", fmt.Errorf("mapping of %s failed, wrong type (%T)", key, v)
	}

	return "", fmt.Errorf("mapping of %s failed, missing value", key)
}

func redirectToLogin(ctx *gin.Context, session sessions.Session, oauth2Config *oauth2.Config) {
	state, err := state()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "generating state failed"})
		return
	}
	session.Set(stateSession, state)
	err = session.Save()
	if err != nil {
		log.Println(err)
	}

	ctx.Redirect(http.StatusFound, oauth2Config.AuthCodeURL(state))
	log.Println("abort", ctx.Request.URL.String())
	ctx.Abort()
}

func AuthorizeBlockedUser(ctx *gin.Context) {
	user, ok := busdb.UserFromContext(ctx)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "no user in context"})
		return
	}

	if user.Blocked {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user is blocked"})
		return
	}

	ctx.Next()
}

func AuthorizeRole(roles []role.Role) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, ok := busdb.UserFromContext(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "no user in context"})
			return
		}

		if !role.UserHasRoles(user, roles) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("missing role %s has %s", roles, user.Roles)})
			return
		}

		ctx.Next()
	}
}

func callback(config *AuthConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		state := session.Get(stateSession)
		if state == "" {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "state missing"})
			return
		}

		if state != ctx.Query("state") {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "state mismatch"})
			return
		}

		oauth2Token, err := config.OAuth2.Exchange(ctx, ctx.Query("code"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": gin.H{"error": fmt.Sprintf("oauth2 exchange failed: %s", err)}})
			return
		}

		// Extract the ID Token from OAuth2 token.
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "missing id token"})
			return
		}

		oidcCtx, cancel := oidcCtx(ctx)
		defer cancel()

		verifier, err := config.Verifier(oidcCtx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "could not verify: " + err.Error()})
			return
		}

		// Parse and verify ID Token payload.
		idToken, err := verifier.Verify(oidcCtx, rawIDToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "token verification failed: " + err.Error()})
			return
		}

		// Extract custom claims
		var claims map[string]interface{}
		if err := idToken.Claims(&claims); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "claim extraction failed"})
			return
		}

		session.Set(userSession, claims)
		err = session.Save()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not save session %s", err)})
			return
		}

		ctx.Redirect(http.StatusFound, "/")
	}
}

func state() (string, error) {
	rnd := make([]byte, 32)
	if _, err := rand.Read(rnd); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(rnd), nil
}
