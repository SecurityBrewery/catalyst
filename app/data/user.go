package data

import (
	"context"
	"fmt"
	"time"

	"github.com/SecurityBrewery/catalyst/app/auth/password"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/pointer"
)

func createTestUser(ctx context.Context, queries *sqlc.Queries) (sqlc.User, error) {
	passwordHash, tokenKey, err := password.Hash("1234567890")
	if err != nil {
		return sqlc.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	return queries.InsertUser(ctx, sqlc.InsertUserParams{
		ID:           "u_test",
		Username:     "u_test",
		Name:         pointer.Pointer("Test User"),
		Email:        pointer.Pointer("user@catalyst-soar.com"),
		Active:       true,
		PasswordHash: passwordHash,
		TokenKey:     tokenKey,
		Created:      time.Now(),
		Updated:      time.Now(),
	})
}
