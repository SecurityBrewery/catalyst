package password

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHash(t *testing.T) {
	t.Parallel()

	type args struct {
		password string
	}

	tests := []struct {
		name    string
		args    args
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "Hash valid password",
			args: args{
				password: "securePassword123!",
			},
			wantErr: require.NoError,
		},
		{
			name: "Long password",
			args: args{
				password: strings.Repeat("a", 75),
			},
			wantErr: require.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotHashedPassword, gotTokenKey, err := Hash(tt.args.password)
			tt.wantErr(t, err, "Hash() should not return an error")

			if err != nil {
				return
			}

			assert.NotEmpty(t, gotHashedPassword, "Hash() gotHashedPassword should not be empty")
			assert.NotEmpty(t, gotTokenKey, "Hash() gotTokenKey should not be empty")

			require.NoError(t, bcrypt.CompareHashAndPassword([]byte(gotHashedPassword), []byte(tt.args.password)), "Hash() hashed password does not match original password")

			assert.GreaterOrEqual(t, len(gotTokenKey), 43, "Hash() gotTokenKey should be at least 43 characters long")
		})
	}
}
