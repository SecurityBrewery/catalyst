package fakedata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_fakeTicketComment(t *testing.T) {
	t.Parallel()

	assert.NotEmpty(t, fakeTicketComment())
}

func Test_fakeTicketDescription(t *testing.T) {
	t.Parallel()

	assert.NotEmpty(t, fakeTicketDescription())
}

func Test_fakeTicketTask(t *testing.T) {
	t.Parallel()

	assert.NotEmpty(t, fakeTicketTask())
}

func Test_fakeTicketTimelineMessage(t *testing.T) {
	t.Parallel()

	assert.NotEmpty(t, fakeTicketTimelineMessage())
}

func Test_random(t *testing.T) {
	t.Parallel()

	type args[T any] struct {
		e []T
	}

	type testCase[T any] struct {
		name string
		args args[T]
	}

	tests := []testCase[int]{
		{
			name: "Test random",
			args: args[int]{e: []int{1, 2, 3, 4, 5}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := random(tt.args.e)

			assert.Contains(t, tt.args.e, got)
		})
	}
}
