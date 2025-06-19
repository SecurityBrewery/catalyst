package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/urfave/cli/v3"

	"github.com/SecurityBrewery/catalyst/app/auth/password"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/service"
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

	admin, err := catalyst.Queries.CreateUser(ctx, sqlc.CreateUserParams{
		ID:              service.GenerateID("u"),
		Name:            command.Args().Get(0),
		Email:           command.Args().Get(0),
		EmailVisibility: false,
		Username:        "admin",
		PasswordHash:    passwordHash,
		TokenKey:        tokenKey,
		Avatar:          "",
		Verified:        true,
	})
	if err != nil {
		return err
	}

	if err := catalyst.Queries.AssignGroupToUser(ctx, sqlc.AssignGroupToUserParams{
		UserID:  admin.ID,
		GroupID: "admin",
	}); err != nil {
		return err
	}

	slog.InfoContext(ctx, "creating test user", "id", admin.ID)

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

	slog.InfoContext(ctx, "setting password for user", "id", user.ID)

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

	slog.InfoContext(ctx, "deleting user", "id", mail)

	return nil
}
