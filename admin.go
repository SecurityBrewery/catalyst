package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/urfave/cli/v3"

	"github.com/SecurityBrewery/catalyst/app/auth/password"
	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

func adminCreate(ctx context.Context, command *cli.Command) error {
	catalyst, cleanup, err := setup(ctx, command)
	if err != nil {
		return fmt.Errorf("failed to setup catalyst: %w", err)
	}

	defer cleanup()

	if command.Args().Len() != 2 {
		return errors.New("usage: catalyst admin create <email> <password>")
	}

	passwordHash, tokenKey, err := password.Hash(command.Args().Get(1))
	if err != nil {
		return errors.New("failed to hash password: " + err.Error())
	}

	user, err := catalyst.Queries.CreateUser(ctx, sqlc.CreateUserParams{
		ID:           database.GenerateID("u"),
		Name:         command.Args().Get(0),
		Email:        command.Args().Get(0),
		Username:     "admin",
		PasswordHash: passwordHash,
		TokenKey:     tokenKey,
		Avatar:       "",
		Verified:     true,
		Created:      time.Now().UTC().Format(time.RFC3339),
		Updated:      time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	if err := catalyst.Queries.AssignGroupToUser(ctx, sqlc.AssignGroupToUserParams{
		UserID:  user.ID,
		GroupID: "admin",
	}); err != nil {
		return err
	}

	slog.InfoContext(ctx, "Creating admin", "id", user.ID, "email", user.Email)

	return nil
}

func adminSetPassword(ctx context.Context, command *cli.Command) error {
	catalyst, cleanup, err := setup(ctx, command)
	if err != nil {
		return fmt.Errorf("failed to setup catalyst: %w", err)
	}

	defer cleanup()

	if command.Args().Len() != 2 {
		return errors.New("usage: catalyst admin set-password <email> <password>")
	}

	user, err := catalyst.Queries.UserByEmail(ctx, command.Args().Get(0))
	if err != nil {
		return err
	}

	passwordHash, tokenKey, err := password.Hash(command.Args().Get(1))
	if err != nil {
		return errors.New("failed to hash password: " + err.Error())
	}

	if _, err := catalyst.Queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:           user.ID,
		PasswordHash: sql.NullString{String: passwordHash, Valid: true},
		TokenKey:     sql.NullString{String: tokenKey, Valid: true},
	}); err != nil {
		return err
	}

	slog.InfoContext(ctx, "Setting password for admin", "id", user.ID, "email", user.Email)

	return nil
}

func adminDelete(ctx context.Context, command *cli.Command) error {
	catalyst, cleanup, err := setup(ctx, command)
	if err != nil {
		return fmt.Errorf("failed to setup catalyst: %w", err)
	}

	defer cleanup()

	if command.Args().Len() != 1 {
		return errors.New("usage: catalyst admin delete <email>")
	}

	mail := command.Args().Get(0)

	user, err := catalyst.Queries.UserByEmail(ctx, mail)
	if err != nil {
		return err
	}

	if err := catalyst.Queries.DeleteUser(ctx, user.ID); err != nil {
		return err
	}

	slog.InfoContext(ctx, "Deleted admin", "id", user.ID, "email", mail)

	return nil
}
