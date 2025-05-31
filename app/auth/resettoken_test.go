package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

func TestService_createResetToken(t *testing.T) {
	type fields struct {
		config *Config
	}
	type args struct {
		createUser    *sqlc.User
		tokenDuration time.Duration
		waitDuration  time.Duration
		verifyUser    *sqlc.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "valid token",
			fields: fields{
				config: &Config{AppSecret: "testsecret"},
			},
			args: args{
				createUser:    &sqlc.User{ID: "testuser", Tokenkey: "testtoken"},
				tokenDuration: time.Hour,
				waitDuration:  0,
				verifyUser:    &sqlc.User{ID: "testuser", Tokenkey: "testtoken"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "expired token",
			fields: fields{
				config: &Config{AppSecret: "testsecret"},
			},
			args: args{
				createUser:    &sqlc.User{ID: "testuser", Tokenkey: "testtoken"},
				tokenDuration: 0,
				waitDuration:  time.Second,
				verifyUser:    &sqlc.User{ID: "testuser", Tokenkey: "testtoken"},
			},
			wantErr: assert.Error,
		},
		{
			name: "invalid token",
			fields: fields{
				config: &Config{AppSecret: "testsecret"},
			},
			args: args{
				createUser:    &sqlc.User{ID: "testuser", Tokenkey: "testtoken"},
				tokenDuration: time.Hour,
				waitDuration:  0,
				verifyUser:    &sqlc.User{ID: "invaliduser", Tokenkey: "invalidtoken"},
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				config: tt.fields.config,
			}
			got, err := s.createResetToken(tt.args.createUser, tt.args.tokenDuration)
			require.NoError(t, err, "createResetToken()")

			time.Sleep(tt.args.waitDuration)

			err = s.verifyResetToken(got, tt.args.verifyUser)
			tt.wantErr(t, err, "verifyResetToken()")
		})
	}
}
