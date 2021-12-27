package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/arangodb/go-driver"

	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/restapi/api"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/userdata"
)

func userdataID(id string) []driver.DocumentID {
	return []driver.DocumentID{driver.DocumentID(fmt.Sprintf("%s/%s", database.UserDataCollectionName, id))}
}

func (s *Service) GetUserData(ctx context.Context, params *userdata.GetUserDataParams) *api.Response {
	userData, err := s.database.UserDataGet(ctx, params.ID)
	return s.response("GetUserData", nil, userData, err)
}

func (s *Service) ListUserData(ctx context.Context) *api.Response {
	userData, err := s.database.UserDataList(ctx)
	return s.response("ListUserData", nil, userData, err)
}

func (s *Service) UpdateUserData(ctx context.Context, params *userdata.UpdateUserDataParams) *api.Response {
	userData, err := s.database.UserDataUpdate(ctx, params.ID, params.Userdata)
	return s.response("UpdateUserData", userdataID(userData.ID), userData, err)
}

func (s *Service) CurrentUserData(ctx context.Context) *api.Response {
	user, ok := busdb.UserFromContext(ctx)
	if !ok {
		return s.response("CurrentUserData", userdataID(user.ID), nil, errors.New("no user in context"))
	}
	userData, err := s.database.UserDataGet(ctx, user.ID)
	return s.response("GetUserData", nil, userData, err)
}

func (s *Service) UpdateCurrentUserData(ctx context.Context, params *userdata.UpdateCurrentUserDataParams) *api.Response {
	user, ok := busdb.UserFromContext(ctx)
	if !ok {
		return s.response("UpdateCurrentUserData", userdataID(user.ID), nil, errors.New("no user in context"))
	}

	userData, err := s.database.UserDataUpdate(ctx, user.ID, params.Userdata)
	return s.response("UpdateCurrentUserData", userdataID(user.ID), userData, err)
}
