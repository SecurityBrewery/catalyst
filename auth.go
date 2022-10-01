package catalyst

import (
	"context"

	maut "github.com/cugu/maut/auth"

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

func (c *catalystResolver) UserCreateIfNotExists(ctx context.Context, user *maut.User, password string) error {
	_, err := c.database.UserGet(ctx, user.ID)
	if err != nil {
		_, err := c.database.UserCreate(ctx, &model.UserForm{
			Apikey:   user.APIKey,
			Blocked:  user.Blocked,
			ID:       user.ID,
			Password: &password,
			Roles:    user.Roles,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *catalystResolver) User(ctx context.Context, userID string) (*maut.User, error) {
	user, err := c.database.UserGet(ctx, userID)
	if err != nil {
		return nil, err
	}

	return mapMautUser(user), nil
}

func (c *catalystResolver) UserAPIKeyByHash(ctx context.Context, hash string) (*maut.User, error) {
	user, err := c.database.UserAPIKeyByHash(ctx, hash)
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
		r = &maut.Role{
			Name: "admin",
			// Permissions: role.Strings(role.Explode(role.Admin)), TODO
		}
	case "engineer":
		r = &maut.Role{
			Name: "engineer",
			// Permissions: role.Strings(role.Explode(role.Engineer)), TODO
		}
	case "analyst":
		r = &maut.Role{
			Name: "analyst",
			// Permissions: role.Strings(role.Explode(role.Analyst)), TODO
		}
	}

	return r, nil
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
