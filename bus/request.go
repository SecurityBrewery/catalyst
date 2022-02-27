package bus

import (
	"encoding/json"
	"log"

	"github.com/arangodb/go-driver"
	emitter "github.com/emitter-io/go/v2"
)

const ChannelRequest = "request"

type RequestMsg struct {
	IDs      []driver.DocumentID `json:"ids"`
	Function string              `json:"function"`
	User     string              `json:"user"`
}

func (b *Bus) PublishRequest(user, f string, ids []driver.DocumentID) error {
	return b.jsonPublish(&RequestMsg{
		User:     user,
		Function: f,
		IDs:      ids,
	}, ChannelRequest, b.config.requestKey)
}

func (b *Bus) SubscribeRequest(f func(msg *RequestMsg)) error {
	return b.safeSubscribe(b.config.requestKey, ChannelRequest, func(c *emitter.Client, m emitter.Message) {
		msg := &RequestMsg{}
		if err := json.Unmarshal(m.Payload(), msg); err != nil {
			log.Println(err)
			return
		}
		go f(msg)
	})
}
