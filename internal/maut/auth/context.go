package auth

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/exp/slices"
)

type contextKey string

const (
	userContextKey       contextKey = "user"
	permissionContextKey contextKey = "permission"
)

func (a *Authenticator) setUserContext(r *http.Request, user *User) *http.Request {
	var permissions []string
	for _, role := range user.Roles {
		userPermissions, err := a.resolver.Role(r.Context(), role)
		if err != nil {
			continue
		}
		permissions = append(permissions, userPermissions.Permissions...)
	}
	slices.Sort(permissions)
	permissions = slices.Compact(permissions)

	r = r.WithContext(context.WithValue(r.Context(), userContextKey, user))
	r = r.WithContext(context.WithValue(r.Context(), permissionContextKey, permissions))

	return r
}

func UserContext(ctx context.Context, user *User, permissions []string) context.Context {
	ctx = context.WithValue(ctx, userContextKey, user)
	ctx = context.WithValue(ctx, permissionContextKey, permissions)

	return ctx
}

func UserFromContext(ctx context.Context) (*User, []string, bool) {
	u, ok := ctx.Value(userContextKey).(*User)
	if !ok {
		return nil, nil, false
	}

	p, ok := ctx.Value(permissionContextKey).([]string)

	return u, p, ok
}

func (a *Authenticator) setContextClaims(ctx context.Context, r *http.Request, userID string) (*http.Request, error) {
	user, err := a.resolver.User(ctx, userID)
	if err != nil {
		return nil, err
	}

	return a.setUserContext(r, user), nil
}

func mapClaims(claims map[string]any, config *UserCreateConfig) (*User, error) {
	username, err := getString(claims, config.OIDCClaimUsername)
	if err != nil {
		return nil, err
	}

	email, err := getString(claims, config.OIDCClaimEmail)
	if err != nil {
		email = ""
	}

	name, err := getString(claims, config.OIDCClaimName)
	if err != nil {
		name = ""
	}

	groups, err := getStringArray(claims, config.OIDCClaimGroups)
	if err != nil {
		groups = []string{}
	}

	s := config.AuthDefaultRoles
	if slices.Contains(config.AuthAdminUsers, username) {
		s = append(s, "admin")
	}

	return &User{
		ID:      username,
		APIKey:  false,
		Blocked: config.AuthBlockNew,
		Email:   &email,
		Groups:  groups,
		Hash:    nil,
		Name:    &name,
		Roles:   s,
	}, nil
}

func getString(m map[string]any, key string) (string, error) {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s, nil
		}

		return "", fmt.Errorf("mapping of %s failed, wrong type (%T)", key, v)
	}

	return "", fmt.Errorf("mapping of %s failed, missing value", key)
}

func getStringArray(m map[string]any, key string) ([]string, error) {
	if v, ok := m[key]; ok {
		if s, ok := v.([]string); ok {
			return s, nil
		}

		return nil, fmt.Errorf("mapping of %s failed, wrong type (%T)", key, v)
	}

	return nil, fmt.Errorf("mapping of %s failed, missing value", key)
}
