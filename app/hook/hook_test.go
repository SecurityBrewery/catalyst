package hook

import (
	"context"
	"testing"
)

func TestHook_Publish(t *testing.T) {
	t.Parallel()

	type fields struct {
		subscribers []func(ctx context.Context, table string, record any)
	}

	type args struct {
		table  string
		record any
	}

	var called bool

	subscriber := func(_ context.Context, _ string, _ any) {
		called = true
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "publish with no subscribers",
			fields: fields{
				subscribers: nil,
			},
			args: args{
				table:  "test_table",
				record: "test_record",
			},
			want: false,
		},
		{
			name: "publish with one subscriber",
			fields: fields{
				subscribers: []func(ctx context.Context, table string, record any){subscriber},
			},
			args: args{
				table:  "test_table",
				record: "test_record",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			called = false
			h := &Hook{
				subscribers: tt.fields.subscribers,
			}
			h.Publish(t.Context(), tt.args.table, tt.args.record)

			if called != tt.want {
				t.Errorf("Hook.Publish() called = %v, want %v", called, tt.want)
			}
		})
	}
}

func TestHook_Subscribe(t *testing.T) {
	t.Parallel()

	type fields struct {
		subscribers []func(ctx context.Context, table string, record any)
	}

	type args struct {
		fn func(ctx context.Context, table string, record any)
	}

	subscriber := func(_ context.Context, _ string, _ any) {}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "subscribe to empty hook",
			fields: fields{
				subscribers: nil,
			},
			args: args{
				fn: subscriber,
			},
			want: 1,
		},
		{
			name: "subscribe to hook with existing subscriber",
			fields: fields{
				subscribers: []func(ctx context.Context, table string, record any){subscriber},
			},
			args: args{
				fn: subscriber,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := &Hook{
				subscribers: tt.fields.subscribers,
			}
			h.Subscribe(tt.args.fn)

			if got := len(h.subscribers); got != tt.want {
				t.Errorf("Hook.Subscribe() subscriber count = %v, want %v", got, tt.want)
			}
		})
	}
}
