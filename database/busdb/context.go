package busdb

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/SecurityBrewery/catalyst/generated/model"
	"github.com/SecurityBrewery/catalyst/role"
)

const (
	userContextKey  = "user"
	groupContextKey = "groups"
)

func SetContext(r *http.Request, ctx context.Context, user *model.UserResponse) {
	user.Roles = role.Strings(role.Explodes(user.Roles))
	ctx.Set(userContextKey, user)

	r.WithContext()
}

func SetGroupContext(ctx *gin.Context, groups []string) {
	ctx.Set(groupContextKey, groups)
}

func UserContext(ctx context.Context, user *model.UserResponse) context.Context {
	user.Roles = role.Strings(role.Explodes(user.Roles))
	return context.WithValue(ctx, userContextKey, user)
}

func UserFromContext(ctx context.Context) (*model.UserResponse, bool) {
	u, ok := ctx.Value(userContextKey).(*model.UserResponse)
	return u, ok
}
