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

func newUserResponseID(user *model.NewUserResponse) []driver.DocumentID {
	if user == nil {
		return nil
	}

	return userID(user.ID)
}

func userID(id string) []driver.DocumentID {
	return []driver.DocumentID{driver.DocumentID(fmt.Sprintf("%s/%s", database.UserCollectionName, id))}
}

func (s *Service) ListUsers(ctx context.Context) ([]*model.UserResponse, error) {
	return s.database.UserList(ctx)
}

func (s *Service) CreateUser(ctx context.Context, form *model.UserForm) (doc *model.NewUserResponse, err error) {
	defer s.publishRequest(ctx, err, "CreateUser", newUserResponseID(doc))

	return s.database.UserCreate(ctx, form)
}

func (s *Service) GetUser(ctx context.Context, id string) (*model.UserResponse, error) {
	decodedValue, err := url.QueryUnescape(id)
	if err == nil {
		id = decodedValue
	}
	return s.database.UserGet(ctx, id)
}

func (s *Service) UpdateUser(ctx context.Context, id string, form *model.UserForm) (doc *model.UserResponse, err error) {
	decodedValue, err := url.QueryUnescape(id)
	if err == nil {
		id = decodedValue
	}

	defer s.publishRequest(ctx, err, "UpdateUser", userID(id))

	return s.database.UserUpdate(ctx, id, form)
}

func (s *Service) DeleteUser(ctx context.Context, id string) (err error) {
	decodedValue, err := url.QueryUnescape(id)
	if err == nil {
		id = decodedValue
	}

	defer s.publishRequest(ctx, err, "DeleteUser", userID(id))

	return s.database.UserDelete(ctx, id)
}

func (s *Service) CurrentUser(ctx context.Context) (*model.UserResponse, error) {
	user, _, ok := maut.UserFromContext(ctx)
	if !ok {
		return nil, errors.New("no user in context")
	}
	s.publishRequest(ctx, nil, "CurrentUser", userID(user.ID))

	return &model.UserResponse{
		ID:      user.ID,
		Apikey:  user.APIKey,
		Blocked: user.Blocked,
		Roles:   user.Roles,
	}, nil
}
