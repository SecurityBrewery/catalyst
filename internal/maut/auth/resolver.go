package auth

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/crypto/argon2"
)

type LoginResolver interface {
	UserCreateIfNotExists(ctx context.Context, user *User, password string) error
	User(ctx context.Context, userID string) (*User, error)
	UserAPIKeyByHash(ctx context.Context, key string) (*User, error)
	UserByIDAndPassword(ctx context.Context, username string, password string) (*User, error)
	Role(ctx context.Context, roleID string) (*Role, error)
}

var _ LoginResolver = &MemoryResolver{}

type MemoryResolver struct {
	users map[string]*User
	roles map[string]*Role
	items map[string][]string
}

func NewMemoryResolver() *MemoryResolver {
	return &MemoryResolver{
		users: make(map[string]*User),
		roles: make(map[string]*Role),
		items: make(map[string][]string),
	}
}

func (t *MemoryResolver) UserCreateIfNotExists(ctx context.Context, user *User, password string) error {
	if _, ok := t.users[user.ID]; ok {
		return nil
	}

	if password != "" {
		salt, err := salt()
		if err != nil {
			return err
		}
		user.Salt = salt
		user.Hash = argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	}

	t.users[user.ID] = user

	return nil
}

func salt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}

func (t *MemoryResolver) User(ctx context.Context, userID string) (*User, error) {
	if user, ok := t.users[userID]; ok {
		return user, nil
	}

	return nil, errors.New("user not found")
}

func (t *MemoryResolver) UserAPIKeyByHash(ctx context.Context, key string) (*User, error) {
	for _, user := range t.users {
		if user.APIKey && user.Hash != nil && bytes.Equal(argon2.IDKey([]byte(key), user.Salt, 1, 64*1024, 4, 32), user.Hash) {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

func (t *MemoryResolver) UserByIDAndPassword(ctx context.Context, username string, password string) (*User, error) {
	if user, ok := t.users[username]; ok {
		if !user.APIKey && user.Hash != nil && bytes.Equal(argon2.IDKey([]byte(password), user.Salt, 1, 64*1024, 4, 32), user.Hash) {
			return user, nil
		}
	}

	return nil, errors.New("wrong username or password")
}

func (t *MemoryResolver) Role(ctx context.Context, roleID string) (*Role, error) {
	if role, ok := t.roles[roleID]; ok {
		return role, nil
	}

	return nil, errors.New("role not found")
}

func (t *MemoryResolver) RoleCreateIfNotExists(ctx context.Context, role *Role) error {
	if _, ok := t.roles[role.Name]; ok {
		return nil
	}

	t.roles[role.Name] = role

	return nil
}
