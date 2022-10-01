package auth

import "testing"

func TestRole_Contains(t *testing.T) {
	t.Parallel()
	type fields struct {
		Name        string
		Permissions []string
	}
	type args struct {
		p string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"admin", fields{"admin", nil}, args{"test"}, true},
		{"user", fields{"user", []string{"test"}}, args{"test"}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := Role{
				Name:        tt.fields.Name,
				Permissions: tt.fields.Permissions,
			}
			if got := r.Contains(tt.args.p); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
