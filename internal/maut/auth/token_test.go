package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

func testKey() *rsa.PrivateKey {
	block, _ := pem.Decode([]byte(`-----BEGIN PRIVATE KEY-----
MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAKhPSTDs4cpKfnMc
p86fCkpnuER7bGc+mGkhkw6bE+BnROfrDCFBSjrENLS5JcsenANQ1kYGt9iVW2fd
ZAWUdDoj+t7g6+fDpzY1BzPSUls421Dmu7joDPY8jSdMzFCeg7Lyj0I36bJJ7ooD
VPW6Q0XQcb8FfBiFPAKuY4elj/YDAgMBAAECgYBo2GMWmCmbM0aL/KjH/KiTawMN
nfkMY6DbtK9/5LjADHSPKAt5V8ueygSvI7rYSiwToLKqEptJztiO3gnls/GmFzj1
V/QEvFs6Ux3b0hD2SGpGy1m6NWWoAFlMISRkNiAxo+AMdCi4I1hpk4+bHr9VO2Bv
V0zKFxmgn1R8qAR+4QJBANqKxJ/qJ5+lyPuDYf5s+gkZWjCLTC7hPxIJQByDLICw
iEnqcn0n9Gslk5ngJIGQcKBXIp5i0jWSdKN/hLxwgHECQQDFKGmo8niLzEJ5sa1r
spww8Hc2aJM0pBwceshT8ZgVPnpgmITU1ENsKpJ+y1RTjZD6N0aj9gS9UB/UXdTr
HBezAkEAqkDRTYOtusH9AXQpM3zSjaQijw72Gs9/wx1RxOSsFtVwV6U97CLkV1S+
2HG1/vn3w/IeFiYGfZXLKFR/pA5BAQJAbFeu6IaGM9yFUzaOZDZ8mnAqMp349t6Q
DB5045xJxLLWsSpfJE2Y12H1qvO1XUzYNIgXq5ZQOHBFbYA6txBy/QJBAKDRQN47
6YClq9652X+1lYIY/h8MxKiXpVZVncXRgY6pbj4pmWEAM88jra9Wq6R77ocyECzi
XCqi18A/sl6ymWc=
-----END PRIVATE KEY-----`))
	parseResult, _ := x509.ParsePKCS8PrivateKey(block.Bytes)

	if key, ok := parseResult.(*rsa.PrivateKey); ok {
		return key
	}
	panic("failed to parse private key")
}

func createTestToken(iss string) string {
	return createToken(iss, testKey())
}

func createRandomTestToken(iss string) string {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)

	return createToken(iss, key)
}

func createToken(iss string, key *rsa.PrivateKey) string {
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: key}, nil)
	if err != nil {
		panic(err)
	}
	jBuilder := jwt.Signed(signer)
	jBuilder = jBuilder.Claims(map[string]interface{}{
		"sub":            "alice",
		"iss":            iss,
		"aud":            "123",
		"iat":            1703376000,
		"exp":            1703376000,
		"email":          "alice@wonderland.net",
		"email_verified": true,
		"name":           "Alice Adams",

		"preferred_username": "alice",
		"groups":             []string{"admin", "user"},
	})
	t, err := jBuilder.CompactSerialize()
	if err != nil {
		log.Fatal(err)
	}

	return t
}
