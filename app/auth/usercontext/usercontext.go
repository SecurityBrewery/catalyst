package usercontext

import (
	"context"
	"net/http"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

type Key struct{}

func UserRequest(r *http.Request, user *sqlc.User) *http.Request {
	return r.WithContext(UserContext(r.Context(), user))
}

func UserContext(ctx context.Context, user *sqlc.User) context.Context {
	return context.WithValue(ctx, Key{}, user)
}

func UserFromContext(ctx context.Context) (*sqlc.User, bool) {
	user, ok := ctx.Value(Key{}).(*sqlc.User)
	if !ok {
		return nil, false
	}

	return user, true
}
