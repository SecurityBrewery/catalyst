package auth

import (
	"context"

	"golang.org/x/oauth2"
	"golang.org/x/text/language"

	"github.com/cugu/maut/oidctestserver"
	"github.com/cugu/maut/oidctestserver/user"
)

func keyAuthenticator() *Authenticator {
	authenticator, err := NewAuthenticator(
		context.Background(),
		&Config{
			CookieSecret:     []byte("test"),
			SimpleAuthEnable: false,
			APIKeyAuthEnable: true,
			OIDCAuthEnable:   false,
			UserCreateConfig: &UserCreateConfig{
				AuthBlockNew:     false,
				AuthDefaultRoles: []string{analystRole},
				AuthAdminUsers:   []string{"admin"},
			},
		},
		testResolver(),
	)
	if err != nil {
		panic(err)
	}

	return authenticator
}

func simpleAuthenticator() *Authenticator {
	authenticator, err := NewAuthenticator(
		context.Background(),
		&Config{
			CookieSecret:     []byte("test"),
			SimpleAuthEnable: true,
			APIKeyAuthEnable: false,
			OIDCAuthEnable:   false,
			UserCreateConfig: &UserCreateConfig{
				AuthBlockNew:     false,
				AuthDefaultRoles: []string{analystRole},
				AuthAdminUsers:   []string{"admin"},
			},
		},
		testResolver(),
	)
	if err != nil {
		panic(err)
	}

	return authenticator
}

func oidcAuthenticator() *Authenticator {
	config := &oidctestserver.Config{
		StorageKey: testKey(),
		Users: []*user.User{
			{
				ID:                "id1",
				Username:          "alice",
				Password:          "password",
				Firstname:         "Test",
				Lastname:          "User",
				Email:             "test-user@example.com",
				EmailVerified:     true,
				Phone:             "",
				PhoneVerified:     false,
				PreferredLanguage: language.German,
			},
		},
	}
	oidcTestServer := oidctestserver.New(config)
	oidcTestServer.Server.Start()

	authenticator, err := NewAuthenticator(
		context.Background(),
		&Config{
			CookieSecret:     []byte("test"),
			SimpleAuthEnable: false,
			APIKeyAuthEnable: false,
			OIDCAuthEnable:   true,
			OIDCIssuer:       oidcTestServer.Server.URL + "/",
			OAuth2: &oauth2.Config{
				ClientID:     "api",
				ClientSecret: "secret",
				RedirectURL:  "ignore",
			},
			UserCreateConfig: &UserCreateConfig{
				AuthBlockNew:      false,
				AuthDefaultRoles:  []string{analystRole},
				AuthAdminUsers:    []string{"admin"},
				OIDCClaimUsername: "preferred_username",
				OIDCClaimEmail:    "email",
				OIDCClaimName:     "name",
				OIDCClaimGroups:   "groups",
			},
		},
		testResolver(),
	)
	if err != nil {
		panic(err)
	}

	return authenticator
}
