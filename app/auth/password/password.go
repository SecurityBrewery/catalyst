package password

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) (hashedPassword, tokenKey string, err error) {
	hashedPasswordB, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", fmt.Errorf("failed to hash password: %w", err)
	}

	tokenKey, err = GenerateTokenKey()
	if err != nil {
		return "", "", err
	}

	return string(hashedPasswordB), tokenKey, nil
}

func GenerateTokenKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
