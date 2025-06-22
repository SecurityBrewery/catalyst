package data

import (
	"context"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/SecurityBrewery/catalyst/app/auth/password"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

func createTestUser(ctx context.Context, queries *sqlc.Queries) (sqlc.User, error) {
	passwordHash, tokenKey, err := password.Hash("1234567890")
	if err != nil {
		return sqlc.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	return queries.CreateUser(ctx, sqlc.CreateUserParams{
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

func generateDemoUsers(ctx context.Context, queries *sqlc.Queries, count int) ([]sqlc.User, error) {
	records := make([]sqlc.User, 0, count)

	// create the test user
	user, err := queries.GetUser(ctx, "u_test")
	if err != nil {
		newUser, err := createTestUser(ctx, queries)
		if err != nil {
			return nil, err
		}

		records = append(records, newUser)
	} else {
		records = append(records, user)
	}

	for range count - 1 {
		newUser, err := createFakeUser(ctx, queries)
		if err != nil {
			return nil, err
		}

		records = append(records, newUser)
	}

	return records, nil
}

func createFakeUser(ctx context.Context, queries *sqlc.Queries) (sqlc.User, error) {
	id := "u_" + gofakeit.UUID()

	username := gofakeit.Username()

	passwordHash, tokenKey, err := password.Hash(gofakeit.Password(true, true, true, true, false, 16))
	if err != nil {
		return sqlc.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	return queries.CreateUser(ctx, sqlc.CreateUserParams{
		ID:           id,
		Name:         gofakeit.Name(),
		Email:        username + "@catalyst-soar.com",
		Username:     username,
		PasswordHash: passwordHash,
		TokenKey:     tokenKey,
		Verified:     gofakeit.Bool(),
		Created:      dateTime(gofakeit.PastDate()),
		Updated:      dateTime(gofakeit.PastDate()),
	})
}
