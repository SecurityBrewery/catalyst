package python_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SecurityBrewery/catalyst/reaction/action/python"
)

func TestPython_Run(t *testing.T) {
	t.Parallel()

	type fields struct {
		Requirements string
		Script       string
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
		{
			name: "requests",
			fields: fields{
				Requirements: "requests",
				Script:       "import requests\nprint(requests.get('https://xkcd.com/2961/info.0.json').text)",
			},
			args: args{
				payload: "test",
			},
			want:    []byte("{\"month\": \"7\", \"num\": 2961, \"link\": \"\", \"year\": \"2024\", \"news\": \"\", \"safe_title\": \"CrowdStrike\", \"transcript\": \"\", \"alt\": \"We were going to try swordfighting, but all my compiling is on hold.\", \"img\": \"https://imgs.xkcd.com/comics/crowdstrike.png\", \"title\": \"CrowdStrike\", \"day\": \"19\"}\n"),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			a := &python.Python{
				Requirements: tt.fields.Requirements,
				Script:       tt.fields.Script,
			}
			got, err := a.Run(ctx, tt.args.payload)
			tt.wantErr(t, err)

			assert.Equal(t, tt.want, got)
		})
	}
}
