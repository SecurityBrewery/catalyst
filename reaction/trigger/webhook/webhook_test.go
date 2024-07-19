package webhook

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_bearerToken(t *testing.T) {
	type args struct {
		r *http.Request
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no token",
			args: args{r: &http.Request{}},
			want: "",
		},
		{
			name: "no bearer token",
			args: args{r: &http.Request{Header: map[string][]string{"Authorization": {"xxx"}}}},
			want: "",
		},
		{
			name: "token in header",
			args: args{r: &http.Request{Header: map[string][]string{"Authorization": {"Bearer token"}}}},
			want: "token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, bearerToken(tt.args.r), "bearerToken(%v)", tt.args.r)
		})
	}
}
