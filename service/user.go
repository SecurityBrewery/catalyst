package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/arangodb/go-driver"

	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/restapi/api"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/users"
)

func userID(id string) []driver.DocumentID {
	return []driver.DocumentID{driver.DocumentID(fmt.Sprintf("%s/%s", database.UserCollectionName, id))}
}

func (s *Service) GetUser(ctx context.Context, params *users.GetUserParams) *api.Response {
	i, err := s.database.UserGet(ctx, params.ID)
	return s.response(ctx, "GetUser", nil, i, err)
}

func (s *Service) ListUsers(ctx context.Context) *api.Response {
	i, err := s.database.UserList(ctx)
	return s.response(ctx, "ListUsers", nil, i, err)
}

func (s *Service) CreateUser(ctx context.Context, params *users.CreateUserParams) *api.Response {
	i, err := s.database.UserCreate(ctx, params.User)
	return s.response(ctx, "CreateUser", userID(i.ID), i, err)
}

func (s *Service) DeleteUser(ctx context.Context, params *users.DeleteUserParams) *api.Response {
	err := s.database.UserDelete(ctx, params.ID)
	return s.response(ctx, "DeleteUser", userID(params.ID), nil, err)
}

func (s *Service) CurrentUser(ctx context.Context) *api.Response {
	user, ok := busdb.UserFromContext(ctx)
	if !ok {
		return s.response(ctx, "CurrentUser", nil, nil, errors.New("no user in context"))
	}
	return s.response(ctx, "CurrentUser", nil, user, nil)
}

func (s *Service) UpdateUser(ctx context.Context, params *users.UpdateUserParams) *api.Response {
	i, err := s.database.UserUpdate(ctx, params.ID, params.User)
	return s.response(ctx, "UpdateUser", userID(i.ID), i, err)
}
