package example

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/Nerzal/gocloak/v11"
)

var _ OIDCProvider = &keycloakServer{}

type keycloakServer struct {
	token  *gocloak.JWT
	client gocloak.GoCloak
}

const (
	keycloakRealm = "maut"
	keycloakURL   = "http://localhost:8081"
)

var keycloakSetupDone sync.Once

func newKeycloakServer() OIDCProvider {
	client := gocloak.NewClient(keycloakURL, gocloak.SetAuthRealms("realms"), gocloak.SetAuthAdminRealms("admin/realms"))
	ctx := context.Background()

	token, err := client.LoginAdmin(ctx, "admin", "admin", "master")
	if err != nil {
		panic(err)
	}

	userID := ""

	keycloakSetupDone.Do(func() {
		_, err := client.GetRealm(ctx, token.AccessToken, keycloakRealm)
		if err != nil {
			_, err := client.CreateRealm(ctx, token.AccessToken, gocloak.RealmRepresentation{
				ID:      gocloak.StringP(keycloakRealm),
				Realm:   gocloak.StringP(keycloakRealm),
				Enabled: gocloak.BoolP(true),
			})
			if err != nil {
				panic(err)
			}
		} else {
			// enable realm
			err := client.UpdateRealm(ctx, token.AccessToken, gocloak.RealmRepresentation{
				ID:      gocloak.StringP(keycloakRealm),
				Realm:   gocloak.StringP(keycloakRealm),
				Enabled: gocloak.BoolP(true),
			})
			if err != nil {
				panic(err)
			}
		}

		// get user by name
		users, err := client.GetUsers(ctx, token.AccessToken, keycloakRealm, gocloak.GetUsersParams{
			Username: gocloak.StringP("alice"),
		})
		if err != nil {
			panic(err)
		}
		if len(users) == 0 {
			userID, err = client.CreateUser(ctx, token.AccessToken, keycloakRealm, gocloak.User{
				Username:      gocloak.StringP("alice"),
				Email:         gocloak.StringP(""), // email is required
				EmailVerified: gocloak.BoolP(true),
				Enabled:       gocloak.BoolP(true),
				FirstName:     gocloak.StringP("Test"),
				LastName:      gocloak.StringP("User"),
			})
			if err != nil {
				panic(err)
			}
		} else {
			userID = *users[0].ID
		}

		// set password
		if err := client.SetPassword(ctx, token.AccessToken, userID, keycloakRealm, "password", false); err != nil {
			panic(err)
		}
	})

	return &keycloakServer{client: client, token: token}
}

func (i *keycloakServer) AddClient(clientID, redirectURL string) string {
	ctx := context.Background()
	clientUID := i.clientUID(ctx, clientID, redirectURL)

	if clientUID == "" {
		var err error
		clientUID, err = i.client.CreateClient(context.Background(), i.token.AccessToken, keycloakRealm, gocloak.Client{
			Enabled:                      gocloak.BoolP(true),
			ClientID:                     gocloak.StringP(clientID),
			Protocol:                     gocloak.StringP("openid-connect"),
			RedirectURIs:                 &[]string{redirectURL},
			AuthorizationServicesEnabled: gocloak.BoolP(true),
			ServiceAccountsEnabled:       gocloak.BoolP(true),
			PublicClient:                 gocloak.BoolP(false),
		})
		if err != nil {
			panic(err)
		}
	}

	cred, err := i.client.GetClientSecret(ctx, i.token.AccessToken, keycloakRealm, clientUID)
	if err != nil {
		panic(err)
	}

	return *cred.Value
}

func (i *keycloakServer) clientUID(ctx context.Context, clientID string, redirectURL string) string {
	clientUID := ""
	clients, err := i.client.GetClients(context.Background(), i.token.AccessToken, keycloakRealm, gocloak.GetClientsParams{
		ClientID: &clientID,
	})
	if err != nil {
		return ""
	}

	if len(clients) > 0 {
		clientUID = *clients[0].ID
		err = i.client.UpdateClient(ctx, i.token.AccessToken, keycloakRealm, gocloak.Client{
			ID:           gocloak.StringP(*clients[0].ID),
			ClientID:     gocloak.StringP(clientID),
			RedirectURIs: &[]string{redirectURL},
		})
		if err != nil {
			return ""
		}
	}

	return clientUID
}

func (i *keycloakServer) URL() string {
	return keycloakURL + "/realms/" + keycloakRealm
}

func (i *keycloakServer) Close() {}

func keycloakLogin(client *http.Client, port string, initialResp *http.Response) (*http.Response, error) {
	b, _ := io.ReadAll(initialResp.Body)

	re := regexp.MustCompile(`action="([^"]+)"`)
	matches := re.FindStringSubmatch(string(b))
	if len(matches) < 2 {
		return nil, errors.New("could not find login form: " + string(b))
	}

	target := strings.ReplaceAll(matches[1], "&amp;", "&")

	return client.PostForm(target, url.Values{
		"credentialId": []string{},
		"username":     []string{"alice"},
		"password":     []string{"password"},
	})
}
