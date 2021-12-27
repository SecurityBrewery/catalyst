package bus

import (
	"encoding/json"
	"log"

	"github.com/arangodb/go-driver"
	emitter "github.com/emitter-io/go/v2"
)

const channelRequest = "request"

type RequestMsg struct {
	IDs      []driver.DocumentID `json:"ids"`
	Function string              `json:"function"`
}

func (b *Bus) PublishRequest(f string, ids []driver.DocumentID) error {
	return b.jsonPublish(&RequestMsg{
		Function: f,
		IDs:      ids,
	}, channelRequest, b.config.requestKey)
}

func (b *Bus) SubscribeRequest(f func(msg *RequestMsg)) error {
	return b.safeSubscribe(b.config.requestKey, channelRequest, func(c *emitter.Client, m emitter.Message) {
		var msg RequestMsg
		if err := json.Unmarshal(m.Payload(), &msg); err != nil {
			log.Println(err)
			return
		}
		go f(&msg)
	})
}
