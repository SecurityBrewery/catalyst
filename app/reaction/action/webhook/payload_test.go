package webhook_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/SecurityBrewery/catalyst/app/reaction/action/webhook"
)

func TestEncodeBody(t *testing.T) {
	t.Parallel()

	type args struct {
		requestBody io.Reader
	}

	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			name: "utf8",
			args: args{
				requestBody: bytes.NewBufferString("body"),
			},
			want:  "body",
			want1: false,
		},
		{
			name: "non-utf8",
			args: args{
				requestBody: bytes.NewBufferString("body\xe0"),
			},
			want:  "Ym9keeA=",
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, got1 := webhook.EncodeBody(tt.args.requestBody)
			if got != tt.want {
				t.Errorf("EncodeBody() got = %v, want %v", got, tt.want)
			}

			if got1 != tt.want1 {
				t.Errorf("EncodeBody() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
