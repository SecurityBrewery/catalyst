package internal

import (
	"crypto/rsa"
)

type Service struct {
	Keys map[string]*rsa.PublicKey
}
