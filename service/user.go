package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/arangodb/go-driver"
	maut "github.com/cugu/maut/auth"

	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/generated/model"
)

func newUserResponseID(user *model.NewUserResponse) []driver.DocumentID {
	if user == nil {
		return nil
	}

	return userID(user.ID)
}

func userResponseID(user *model.UserResponse) []driver.DocumentID {
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

func (s *Service) GetUser(ctx context.Context, s2 string) (*model.UserResponse, error) {
	return s.database.UserGet(ctx, s2)
}

func (s *Service) UpdateUser(ctx context.Context, s2 string, form *model.UserForm) (doc *model.UserResponse, err error) {
	defer s.publishRequest(ctx, err, "UpdateUser", userID(s2))

	return s.database.UserUpdate(ctx, s2, form)
}

func (s *Service) DeleteUser(ctx context.Context, s2 string) (err error) {
	defer s.publishRequest(ctx, err, "DeleteUser", userID(s2))

	return s.database.UserDelete(ctx, s2)
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
