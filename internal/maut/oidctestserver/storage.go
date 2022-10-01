package oidctestserver

import (
	"crypto/rsa"

	"github.com/cugu/maut/oidctestserver/internal"
	"github.com/cugu/maut/oidctestserver/user"
)

// serviceKey1 is a public key which will be used for the JWT Profile Authorization Grant
// the corresponding private key is in the service-key1.json (for demonstration purposes)
// var serviceKey1 = &rsa.PublicKey{
// 	N: func() *big.Int {
// 		n, _ := new(big.Int).SetString("00f6d44fb5f34ac2033a75e73cb65ff24e6181edc58845e75a560ac21378284977bb055b1a75b714874e2a2641806205681c09abec76efd52cf40984edcf4c8ca09717355d11ac338f280d3e4c905b00543bdb8ee5a417496cb50cb0e29afc5a0d0471fd5a2fa625bd5281f61e6b02067d4fe7a5349eeae6d6a4300bcd86eef331", 16)
// 		return n
// 	}(),
// 	E: 65537,
// }

func NewTestStorage(key *rsa.PrivateKey, users []*user.User) *internal.Storage {
	userMap := map[string]*user.User{}
	for _, u := range users {
		userMap[u.ID] = u
	}

	return internal.NewStorage(
		key,
		userMap,
		map[string]*internal.Client{},
		map[string]internal.Service{
			// "service": {
			// 	Keys: map[string]*rsa.PublicKey{
			// 		"key1": serviceKey1,
			// 	},
			// },
		},
	)
}
