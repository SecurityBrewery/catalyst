package app

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isIgnored(t *testing.T) {
	t.Parallel()

	type args struct {
		err error
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "error is ignored",
			args: args{err: errors.New("1673167670_multi_match_migrate")},
			want: true,
		},
		{
			name: "error is not ignored",
			args: args{err: errors.New("1673167670_multi_match")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equalf(t, tt.want, isIgnored(tt.args.err), "isIgnored(%v)", tt.args.err)
		})
	}
}
