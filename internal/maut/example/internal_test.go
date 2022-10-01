package example

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/url"

	"golang.org/x/text/language"

	"github.com/cugu/maut/oidctestserver"
	"github.com/cugu/maut/oidctestserver/user"
)

var _ OIDCProvider = &internalOIDCServer{}

type internalOIDCServer struct {
	server *oidctestserver.OIDCServer
}

func newInternalOIDCServer() OIDCProvider {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	server := oidctestserver.New(&oidctestserver.Config{
		StorageKey: key,
		Users: []*user.User{{
			ID:                "id1",
			Username:          "alice",
			Password:          "password",
			Firstname:         "Test",
			Lastname:          "User",
			Email:             "test-user@zitadel.ch",
			EmailVerified:     true,
			Phone:             "",
			PhoneVerified:     false,
			PreferredLanguage: language.German,
		}},
	})
	server.Server.Start()

	return &internalOIDCServer{server}
}

func (i *internalOIDCServer) AddClient(id, redirectURL string) string {
	secret := "secret"
	i.server.RegisterClients(
		oidctestserver.WebClient(id, secret, redirectURL),
	)

	return "secret"
}

func (i *internalOIDCServer) URL() string {
	return i.server.Server.URL + "/"
}

func (i *internalOIDCServer) Close() {
	i.server.Server.Close()
}

func internalLogin(client *http.Client, port string, initialResp *http.Response) (*http.Response, error) {
	return client.PostForm(initialResp.Request.URL.String(), url.Values{
		"id":       initialResp.Request.URL.Query()["authRequestID"],
		"username": []string{"alice"},
		"password": []string{"password"},
	})
}
