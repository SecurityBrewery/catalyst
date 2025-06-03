package usercontext

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

func TestPermissionContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		ctx         context.Context
		user        *sqlc.User
		permissions []string
		wantPerms   []string
		wantOk      bool
	}{
		{
			name:        "Set and get permissions",
			ctx:         context.Background(),
			permissions: []string{"ticket:read", "ticket:write"},
			wantPerms:   []string{"ticket:read", "ticket:write"},
			wantOk:      true,
		},
		{
			name:      "No permissions set",
			ctx:       context.Background(),
			wantPerms: nil,
			wantOk:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Test context functions
			ctx := PermissionContext(tt.ctx, tt.permissions)
			gotPerms, gotOk := PermissionFromContext(ctx)

			if !reflect.DeepEqual(gotPerms, tt.wantPerms) {
				t.Errorf("PermissionFromContext() got perms = %v, want %v", gotPerms, tt.wantPerms)
			}
			if gotOk != tt.wantOk {
				t.Errorf("PermissionFromContext() got ok = %v, want %v", gotOk, tt.wantOk)
			}

			// Test request functions
			req := &http.Request{}
			req = PermissionRequest(req, tt.permissions)
			gotPerms, gotOk = PermissionFromContext(req.Context())

			if !reflect.DeepEqual(gotPerms, tt.wantPerms) {
				t.Errorf("PermissionFromContext() got perms = %v, want %v", gotPerms, tt.wantPerms)
			}
			if gotOk != tt.wantOk {
				t.Errorf("PermissionFromContext() got ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestUserContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		ctx    context.Context
		user   *sqlc.User
		wantOk bool
	}{
		{
			name:   "Set and get user",
			ctx:    context.Background(),
			user:   &sqlc.User{ID: "test-user"},
			wantOk: true,
		},
		{
			name:   "No user set",
			ctx:    context.Background(),
			user:   nil,
			wantOk: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Test context functions
			ctx := UserContext(tt.ctx, tt.user)
			gotUser, gotOk := UserFromContext(ctx)

			if !reflect.DeepEqual(gotUser, tt.user) {
				t.Errorf("UserFromContext() got user = %v, want %v", gotUser, tt.user)
			}
			if gotOk != tt.wantOk {
				t.Errorf("UserFromContext() got ok = %v, want %v", gotOk, tt.wantOk)
			}

			// Test request functions
			req := &http.Request{}
			req = UserRequest(req, tt.user)
			gotUser, gotOk = UserFromContext(req.Context())

			if !reflect.DeepEqual(gotUser, tt.user) {
				t.Errorf("UserFromContext() got user = %v, want %v", gotUser, tt.user)
			}
			if gotOk != tt.wantOk {
				t.Errorf("UserFromContext() got ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
