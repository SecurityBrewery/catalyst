package catalyst

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	maut "github.com/jonas-plum/maut/auth"

	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/generated/model"
)

type catalystResolver struct {
	database *database.Database
}

func newCatalystResolver(db *database.Database) *catalystResolver {
	return &catalystResolver{
		database: db,
	}
}

func (c *catalystResolver) UserCreateIfNotExists(ctx context.Context, user *maut.User, password string) (err error) {
	if user != nil {
		if _, err := c.database.UserGet(ctx, user.ID); err == nil {
			return nil
		}
	}

	if user == nil || user.APIKey {
		_, err = c.database.UserCreateSetupAPIKey(ctx, password)
	} else {
		_, err = c.database.UserCreate(ctx, &model.UserForm{
			Apikey:   user.APIKey,
			Blocked:  user.Blocked,
			ID:       user.ID,
			Password: &password,
			Roles:    user.Roles,
		})
	}

	return err
}

func (c *catalystResolver) User(ctx context.Context, userID string) (*maut.User, error) {
	user, err := c.database.UserGet(ctx, userID)
	if err != nil {
		return nil, err
	}

	return mapMautUser(user), nil
}

func (c *catalystResolver) UserAPIKeyByHash(ctx context.Context, key string) (*maut.User, error) {
	sha256Hash := fmt.Sprintf("%x", sha256.Sum256([]byte(key)))
	user, err := c.database.UserAPIKeyByHash(ctx, sha256Hash)
	if err != nil {
		return nil, err
	}

	return mapMautUser(user), nil
}

func (c *catalystResolver) UserByIDAndPassword(ctx context.Context, username string, password string) (*maut.User, error) {
	user, err := c.database.UserByIDAndPassword(ctx, username, password)
	if err != nil {
		return nil, err
	}

	return mapMautUser(user), nil
}

func (c *catalystResolver) Role(ctx context.Context, roleID string) (r *maut.Role, err error) {
	switch roleID {
	case "admin":
		return Admin, nil
	case "engineer":
		return engineer, nil
	case "analyst":
		return analyst, nil
	}

	return nil, errors.New("role not found")
}

func mapMautUser(user *model.UserResponse) *maut.User {
	return &maut.User{
		ID:      user.ID,
		APIKey:  user.Apikey,
		Blocked: user.Blocked,
		// Email:   user.Email, // TODO
		// Groups:  user.Groups, // TODO
		// Name:    user.Name, // TODO
		Roles: user.Roles,
	}
}
