package busdb

import (
	"context"
	"net/http"

	"github.com/SecurityBrewery/catalyst/generated/model"
	"github.com/SecurityBrewery/catalyst/role"
)

const (
	userContextKey  = "user"
	groupContextKey = "groups"
)

func SetContext(r *http.Request, user *model.UserResponse) *http.Request {
	user.Roles = role.Strings(role.Explodes(user.Roles))

	return r.WithContext(context.WithValue(r.Context(), userContextKey, user))
}

func SetGroupContext(r *http.Request, groups []string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), groupContextKey, groups))
}

func UserContext(ctx context.Context, user *model.UserResponse) context.Context {
	user.Roles = role.Strings(role.Explodes(user.Roles))
	return context.WithValue(ctx, userContextKey, user)
}

func UserFromContext(ctx context.Context) (*model.UserResponse, bool) {
	u, ok := ctx.Value(userContextKey).(*model.UserResponse)
	return u, ok
}
