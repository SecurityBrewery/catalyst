package service

import (
	"context"
	"errors"

	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/restapi/api"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/users"
)

func (s *Service) GetUser(ctx context.Context, params *users.GetUserParams) *api.Response {
	return response(s.database.UserGet(ctx, params.ID))
}

func (s *Service) ListUsers(ctx context.Context) *api.Response {
	return response(s.database.UserList(ctx))
}

func (s *Service) CreateUser(ctx context.Context, params *users.CreateUserParams) *api.Response {
	return response(s.database.UserCreate(ctx, params.User))
}

func (s *Service) DeleteUser(ctx context.Context, params *users.DeleteUserParams) *api.Response {
	return response(nil, s.database.UserDelete(ctx, params.ID))
}

func (s *Service) CurrentUser(ctx context.Context) *api.Response {
	user, ok := busdb.UserFromContext(ctx)
	if !ok {
		return response(nil, errors.New("no user in context"))
	}
	return response(user, nil)
}

func (s *Service) UpdateUser(ctx context.Context, params *users.UpdateUserParams) *api.Response {
	return response(s.database.UserUpdate(ctx, params.ID, params.User))
}
