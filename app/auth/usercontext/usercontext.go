package usercontext

import (
	"context"
	"net/http"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

type userKey struct{}

func UserRequest(r *http.Request, user *sqlc.User) *http.Request {
	return r.WithContext(UserContext(r.Context(), user))
}

func UserContext(ctx context.Context, user *sqlc.User) context.Context {
	return context.WithValue(ctx, userKey{}, user)
}

func UserFromContext(ctx context.Context) (*sqlc.User, bool) {
	user, ok := ctx.Value(userKey{}).(*sqlc.User)
	if !ok {
		return nil, false
	}

	return user, true
}

type permissionKey struct{}

func PermissionRequest(r *http.Request, permissions []string) *http.Request {
	return r.WithContext(PermissionContext(r.Context(), permissions))
}

func PermissionContext(ctx context.Context, permissions []string) context.Context {
	return context.WithValue(ctx, permissionKey{}, permissions)
}

func PermissionFromContext(ctx context.Context) ([]string, bool) {
	permissions, ok := ctx.Value(permissionKey{}).([]string)
	if !ok {
		return nil, false
	}

	return permissions, true
}
