package service

import (
	"context"
	"errors"

	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/restapi/api"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/userdata"
)

func (s *Service) GetUserData(ctx context.Context, params *userdata.GetUserDataParams) *api.Response {
	return response(s.database.UserDataGet(ctx, params.ID))
}

func (s *Service) ListUserData(ctx context.Context) *api.Response {
	return response(s.database.UserDataList(ctx))
}

func (s *Service) UpdateUserData(ctx context.Context, params *userdata.UpdateUserDataParams) *api.Response {
	return response(s.database.UserDataUpdate(ctx, params.ID, params.Userdata))
}

func (s *Service) CurrentUserData(ctx context.Context) *api.Response {
	user, ok := busdb.UserFromContext(ctx)
	if !ok {
		return response(nil, errors.New("no user in context"))
	}
	return s.GetUserData(ctx, &userdata.GetUserDataParams{ID: user.ID})
}

func (s *Service) UpdateCurrentUserData(ctx context.Context, params *userdata.UpdateCurrentUserDataParams) *api.Response {
	user, ok := busdb.UserFromContext(ctx)
	if !ok {
		return response(nil, errors.New("no user in context"))
	}

	return response(s.database.UserDataUpdate(ctx, user.ID, params.Userdata))
}
