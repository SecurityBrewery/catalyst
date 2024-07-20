package python

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPython_Run(t *testing.T) {
	type fields struct {
		Bootstrap string
		Script    string
	}

	type args struct {
		payload string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "empty",
			fields: fields{
				Script: "pass",
			},
			args: args{
				payload: "test",
			},
			want:    []byte(""),
			wantErr: assert.NoError,
		},
		{
			name: "hello world",
			fields: fields{
				Script: "print('hello world')",
			},
			args: args{
				payload: "test",
			},
			want:    []byte("hello world\n"),
			wantErr: assert.NoError,
		},
		{
			name: "echo",
			fields: fields{
				Script: "import sys; print(sys.argv[1])",
			},
			args: args{
				payload: "test",
			},
			want:    []byte("test\n"),
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				Script: "import sys; sys.exit(1)",
			},
			args: args{
				payload: "test",
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			a := &Python{
				Bootstrap: tt.fields.Bootstrap,
				Script:    tt.fields.Script,
			}
			got, err := a.Run(ctx, tt.args.payload)
			tt.wantErr(t, err)

			assert.Equal(t, tt.want, got)
		})
	}
}
