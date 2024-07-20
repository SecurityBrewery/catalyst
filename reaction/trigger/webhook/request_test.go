package webhook_test

import (
	"testing"

	"github.com/SecurityBrewery/catalyst/reaction/trigger/webhook"
)

func Test_isJSON(t *testing.T) {
	t.Parallel()

	type args struct {
		data []byte
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid JSON",
			args: args{
				data: []byte(`{"key": "value"}`),
			},
			want: true,
		},
		{
			name: "invalid JSON",
			args: args{
				data: []byte(`{"key": "value"`),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := webhook.IsJSON(tt.args.data); got != tt.want {
				t.Errorf("isJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
