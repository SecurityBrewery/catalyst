package hooks

import (
	"context"

	"github.com/arangodb/go-driver"
	"github.com/gin-gonic/gin"

	"github.com/SecurityBrewery/catalyst/index"
)

type Hooks struct {
	DatabaseAfterConnectFuncs []func(ctx context.Context, client driver.Client, name string)
	IngestionFilterFunc       func(ctx context.Context, index *index.Index) (string, error)
	TicketReadFilterFunc      func(ctx context.Context) (string, map[string]interface{}, error)
	TicketWriteFilterFunc     func(ctx context.Context) (string, map[string]interface{}, error)
	GetGroupsFunc             func(ctx context.Context, username string) ([]string, error)
}

func (h *Hooks) DatabaseAfterConnect(ctx context.Context, client driver.Client, name string) {
	for _, f := range h.DatabaseAfterConnectFuncs {
		f(ctx, client, name)
	}
}

func (h *Hooks) IngestionFilter(ctx context.Context, index *index.Index) (string, error) {
	if h.IngestionFilterFunc != nil {
		return h.IngestionFilterFunc(ctx, index)
	}
	return "[]", nil
}

func (h *Hooks) TicketReadFilter(ctx context.Context) (string, map[string]interface{}, error) {
	if h.TicketReadFilterFunc != nil {
		return h.TicketReadFilterFunc(ctx)
	}
	return "", nil, nil
}

func (h *Hooks) TicketWriteFilter(ctx context.Context) (string, map[string]interface{}, error) {
	if h.TicketWriteFilterFunc != nil {
		return h.TicketWriteFilterFunc(ctx)
	}
	return "", nil, nil
}

func (h *Hooks) GetGroups(ctx *gin.Context, username string) ([]string, error) {
	if h.GetGroupsFunc != nil {
		return h.GetGroupsFunc(ctx, username)
	}
	return nil, nil
}
