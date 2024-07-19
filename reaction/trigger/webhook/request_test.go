package webhook

import "testing"

func Test_isJSON(t *testing.T) {
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
			if got := isJSON(tt.args.data); got != tt.want {
				t.Errorf("isJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
