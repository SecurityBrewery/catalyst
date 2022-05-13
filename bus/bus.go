package bus

import (
	"github.com/arangodb/go-driver"

	"github.com/SecurityBrewery/catalyst/generated/model"
)

type ResultMsg struct {
	Automation string         `json:"automation"`
	Data       map[string]any `json:"data,omitempty"`
	Target     *model.Origin  `json:"target"`
}

type RequestMsg struct {
	IDs      []driver.DocumentID `json:"ids"`
	Function string              `json:"function"`
	User     string              `json:"user"`
}

type JobMsg struct {
	ID         string         `json:"id"`
	Automation string         `json:"automation"`
	Origin     *model.Origin  `json:"origin"`
	Message    *model.Message `json:"message"`
}

type DatabaseUpdateType string

const (
	DatabaseEntryRead    DatabaseUpdateType = "read"
	DatabaseEntryCreated DatabaseUpdateType = "created"
	DatabaseEntryUpdated DatabaseUpdateType = "updated"
)

type DatabaseUpdateMsg struct {
	IDs  []driver.DocumentID `json:"ids"`
	Type DatabaseUpdateType  `json:"type"`
}

type Bus struct {
	ResultChannel   *Channel[*ResultMsg]
	RequestChannel  *Channel[*RequestMsg]
	JobChannel      *Channel[*JobMsg]
	DatabaseChannel *Channel[*DatabaseUpdateMsg]
}

func New() *Bus {
	return &Bus{
		ResultChannel:   &Channel[*ResultMsg]{},
		RequestChannel:  &Channel[*RequestMsg]{},
		JobChannel:      &Channel[*JobMsg]{},
		DatabaseChannel: &Channel[*DatabaseUpdateMsg]{},
	}
}

type Channel[T any] struct {
	Subscriber []func(T)
}

func (c *Channel[T]) Publish(msg T) {
	for _, s := range c.Subscriber {
		go s(msg)
	}
}

func (c *Channel[T]) Subscribe(handler func(T)) {
	c.Subscriber = append(c.Subscriber, handler)
}
