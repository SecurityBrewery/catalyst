package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/arangodb/go-driver"
	maut "github.com/jonas-plum/maut/auth"

	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/generated/model"
)

func userDataResponseID(userData *model.UserDataResponse) []driver.DocumentID {
	if userData == nil {
		return nil
	}

	return userDataID(userData.ID)
}

func userDataID(id string) []driver.DocumentID {
	return []driver.DocumentID{driver.DocumentID(fmt.Sprintf("%s/%s", database.UserDataCollectionName, id))}
}

func (s *Service) ListUserData(ctx context.Context) (doc []*model.UserDataResponse, err error) {
	return s.database.UserDataList(ctx)
}

func (s *Service) GetUserData(ctx context.Context, id string) (*model.UserDataResponse, error) {
	decodedValue, err := url.QueryUnescape(id)
	if err == nil {
		id = decodedValue
	}

	return s.database.UserDataGet(ctx, id)
}

func (s *Service) UpdateUserData(ctx context.Context, id string, data *model.UserData) (doc *model.UserDataResponse, err error) {
	decodedValue, err := url.QueryUnescape(id)
	if err == nil {
		id = decodedValue
	}

	defer s.publishRequest(ctx, err, "UpdateUserData", userDataResponseID(doc))

	return s.database.UserDataUpdate(ctx, id, data)
}

func (s *Service) CurrentUserData(ctx context.Context) (doc *model.UserDataResponse, err error) {
	user, _, ok := maut.UserFromContext(ctx)
	if !ok {
		return nil, errors.New("no user in context")
	}

	return s.database.UserDataGet(ctx, user.ID)
}

func (s *Service) UpdateCurrentUserData(ctx context.Context, data *model.UserData) (doc *model.UserDataResponse, err error) {
	user, _, ok := maut.UserFromContext(ctx)
	if !ok {
		return nil, errors.New("no user in context")
	}

	defer s.publishRequest(ctx, err, "UpdateCurrentUserData", userDataResponseID(doc))

	return s.database.UserDataUpdate(ctx, user.ID, data)
}
