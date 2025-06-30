package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

func TestService_createResetToken(t *testing.T) {
	t.Parallel()

	type args struct {
		createUser    *sqlc.User
		tokenDuration time.Duration
		waitDuration  time.Duration
		verifyUser    *sqlc.User
	}

	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "valid token",
			args: args{
				createUser:    &sqlc.User{ID: "testuser", Tokenkey: "testtoken"},
				tokenDuration: time.Hour,
				waitDuration:  0,
				verifyUser: &sqlc.User{
					ID:       "testuser",
					Tokenkey: "testtoken",
					Updated:  mustParse(t, "2006-01-02 15:04:05Z", "2025-06-02 19:18:06.292Z"),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "expired token",
			args: args{
				createUser:    &sqlc.User{ID: "testuser", Tokenkey: "testtoken"},
				tokenDuration: 0,
				waitDuration:  time.Second,
				verifyUser: &sqlc.User{
					ID:       "testuser",
					Tokenkey: "testtoken",
					Updated:  mustParse(t, "2006-01-02 15:04:05Z", "2025-06-02 19:18:06.292Z"),
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "invalid token",
			args: args{
				createUser:    &sqlc.User{ID: "testuser", Tokenkey: "testtoken"},
				tokenDuration: time.Hour,
				waitDuration:  0,
				verifyUser: &sqlc.User{
					ID:       "invaliduser",
					Tokenkey: "invalidtoken",
					Updated:  mustParse(t, "2006-01-02 15:04:05Z", "2025-06-02 19:18:06.292Z"),
				},
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := createResetTokenWithDuration(tt.args.createUser, "", "", tt.args.tokenDuration)
			require.NoError(t, err, "createResetToken()")

			time.Sleep(tt.args.waitDuration)

			err = verifyResetToken(got, tt.args.verifyUser, "", "")
			tt.wantErr(t, err, "verifyResetToken()")
		})
	}
}

func mustParse(t *testing.T, layout, value string) time.Time {
	t.Helper()

	parsed, err := time.Parse(layout, value)
	require.NoError(t, err, "mustParse()")

	return parsed
}
