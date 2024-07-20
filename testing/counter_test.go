package testing

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCounter(t *testing.T) {
	type args struct {
		name   string
		repeat int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test Counter",
			args: args{name: "test", repeat: 5},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCounter()

			for range tt.args.repeat {
				c.Increment(tt.args.name)
			}

			assert.Equal(t, tt.want, c.Count(tt.args.name))
		})
	}
}
