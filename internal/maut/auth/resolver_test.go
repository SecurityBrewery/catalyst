package auth

import (
	"golang.org/x/crypto/argon2"
)

func pointer[T any](s T) *T {
	return &s
}

func testResolver() *MemoryResolver {
	return &MemoryResolver{
		users: map[string]*User{
			"alice": {
				APIKey:  false,
				Blocked: false,
				Email:   pointer("alice@wonderland.net"),
				Hash:    argon2.IDKey([]byte("password"), []byte("salt2"), 1, 64*1024, 4, 32),
				ID:      "alice",
				Name:    pointer("Alice"),
				Roles:   []string{analystRole},
				Salt:    []byte("salt2"),
			},
			"mallory": {
				APIKey:  false,
				Blocked: true,
				Email:   pointer("mallory@wonderland.net"),
				Hash:    nil,
				ID:      "mallory",
				Name:    pointer("Mallory"),
				Roles:   []string{analystRole},
			},
			"siem_service": {
				APIKey:  true,
				Blocked: false,
				Hash:    argon2.IDKey([]byte("valid"), []byte("salt"), 1, 64*1024, 4, 32),
				ID:      "siem_service",
				Roles:   []string{analystRole},
				Salt:    []byte("salt"),
			},
		},
	}
}
