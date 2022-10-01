package auth

import (
	"context"
	"reflect"
	"testing"
)

func TestUserContext(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx         context.Context
		user        *User
		permissions []string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{"nil context", args{context.Background(), &User{}, nil}, context.WithValue(context.WithValue(context.Background(), userContextKey, &User{}), permissionContextKey, []string(nil))},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := UserContext(tt.args.ctx, tt.args.user, tt.args.permissions); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getString(t *testing.T) {
	t.Parallel()
	type args struct {
		m   map[string]any
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"nil map", args{nil, "key"}, "", true},
		{"empty map", args{map[string]any{}, "key"}, "", true},
		{"key not found", args{map[string]any{"key": "value"}, "key2"}, "", true},
		{"key found", args{map[string]any{"key": "value"}, "key"}, "value", false},
		{"key found, wrong type", args{map[string]any{"key": 1}, "key"}, "", true},
		{"key found, wrong type", args{map[string]any{"key": true}, "key"}, "", true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := getString(tt.args.m, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("getString() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("getString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getStringArray(t *testing.T) {
	t.Parallel()
	type args struct {
		m   map[string]any
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"nil map", args{nil, "key"}, nil, true},
		{"empty map", args{map[string]any{}, "key"}, nil, true},
		{"key not found", args{map[string]any{"key": "value"}, "key2"}, nil, true},
		{"key found", args{map[string]any{"key": []string{"value"}}, "key"}, []string{"value"}, false},
		{"key found, wrong type", args{map[string]any{"key": 1}, "key"}, nil, true},
		{"key found, wrong type", args{map[string]any{"key": true}, "key"}, nil, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := getStringArray(tt.args.m, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("getStringArray() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getStringArray() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mapClaims(t *testing.T) {
	t.Parallel()
	config := &UserCreateConfig{
		OIDCClaimUsername: "preferred_username",
		OIDCClaimEmail:    "email",
		OIDCClaimName:     "Name",
		OIDCClaimGroups:   "groups",
		AuthAdminUsers:    []string{"test"},
	}

	testClaims := map[string]interface{}{
		"preferred_username": "test",
		"email":              "test@example.com",
		"Name":               "Test",
		"groups":             []string{"group1", "group2"},
	}
	testUser := &User{
		ID:      "test",
		APIKey:  false,
		Blocked: false,
		Email:   pointer("test@example.com"),
		Groups:  []string{"group1", "group2"},
		Name:    pointer("Test"),
		Roles:   []string{AdminRole},
	}

	type args struct {
		claims map[string]any
		config *UserCreateConfig
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
	}{
		{"nil claims", args{nil, config}, nil, true},
		{"empty claims", args{map[string]any{}, config}, nil, true},
		{"only email", args{map[string]any{"email": ""}, config}, nil, true},
		{"only Name", args{map[string]any{"Name": ""}, config}, nil, true},
		{"only groups", args{map[string]any{"groups": ""}, config}, nil, true},
		{"only groups", args{map[string]any{"groups": []string{}}, config}, nil, true},
		{"valid", args{testClaims, config}, testUser, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := mapClaims(tt.args.claims, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("mapClaims() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapClaims() got = %v, want %v", got, tt.want)
			}
		})
	}
}
