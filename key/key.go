package key

import (
	"math/rand"

	"github.com/SecurityBrewery/catalyst/generated/time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Generate() string {
	b := make([]rune, 32)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
