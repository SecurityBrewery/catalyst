package permission

import (
	"reflect"
	"testing"
)

func TestFromJSONArray(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       string
		want        []string
		shouldError bool
	}{
		{
			name:        "Valid JSON array",
			input:       `["ticket:read", "ticket:write"]`,
			want:        []string{"ticket:read", "ticket:write"},
			shouldError: false,
		},
		{
			name:        "Empty array",
			input:       "[]",
			want:        []string{},
			shouldError: false,
		},
		{
			name:        "Invalid JSON",
			input:       "not json",
			want:        nil,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := FromJSONArray(t.Context(), tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromJSONArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToJSONArray(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{
			name:  "Valid permissions array",
			input: []string{"ticket:read", "ticket:write"},
			want:  `["ticket:read","ticket:write"]`,
		},
		{
			name:  "Empty array",
			input: []string{},
			want:  "[]",
		},
		{
			name:  "Nil array",
			input: nil,
			want:  "[]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := ToJSONArray(t.Context(), tt.input)
			if got != tt.want {
				t.Errorf("ToJSONArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
