package data

import (
	"context"
	"fmt"
	"time"

	"github.com/SecurityBrewery/catalyst/app/auth/password"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

func createTestUser(ctx context.Context, queries *sqlc.Queries) (sqlc.User, error) {
	passwordHash, tokenKey, err := password.Hash("1234567890")
	if err != nil {
		return sqlc.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	return queries.InsertUser(ctx, sqlc.InsertUserParams{
		ID:           "u_test",
		Username:     "u_test",
		Name:         "Test User",
		Email:        "user@catalyst-soar.com",
		Verified:     true,
		PasswordHash: passwordHash,
		TokenKey:     tokenKey,
		Created:      dateTime(time.Now()),
		Updated:      dateTime(time.Now()),
	})
}
