package hook

import "context"

type Hook struct {
	subscribers []func(ctx context.Context, table string, record any)
}

func (h *Hook) Publish(ctx context.Context, table string, record any) {
	for _, subscriber := range h.subscribers {
		subscriber(ctx, table, record)
	}
}

func (h *Hook) Subscribe(fn func(ctx context.Context, table string, record any)) {
	h.subscribers = append(h.subscribers, fn)
}
