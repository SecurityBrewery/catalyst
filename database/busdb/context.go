package busdb

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/SecurityBrewery/catalyst/generated/models"
	"github.com/SecurityBrewery/catalyst/role"
)

const (
	userContextKey  = "user"
	groupContextKey = "groups"
)

func SetContext(ctx *gin.Context, user *models.UserResponse) {
	user.Roles = role.Strings(role.Explodes(user.Roles))
	ctx.Set(userContextKey, user)
}

func SetGroupContext(ctx *gin.Context, groups []string) {
	ctx.Set(groupContextKey, groups)
}

func UserContext(ctx context.Context, user *models.UserResponse) context.Context {
	user.Roles = role.Strings(role.Explodes(user.Roles))
	return context.WithValue(ctx, userContextKey, user)
}

func UserFromContext(ctx context.Context) (*models.UserResponse, bool) {
	u, ok := ctx.Value(userContextKey).(*models.UserResponse)
	return u, ok
}
