package bus

import (
	"encoding/json"
	"log"

	emitter "github.com/emitter-io/go/v2"

	"github.com/SecurityBrewery/catalyst/generated/model"
)

const channelResult = "result"

type ResultMsg struct {
	Automation string                 `json:"automation"`
	Data       map[string]interface{} `json:"data,omitempty"`
	Target     *model.Origin          `json:"target"`
}

func (b *Bus) PublishResult(automation string, data map[string]interface{}, target *model.Origin) error {
	return b.jsonPublish(&ResultMsg{Automation: automation, Data: data, Target: target}, channelResult, b.config.resultBusKey)
}

func (b *Bus) SubscribeResult(f func(msg *ResultMsg)) error {
	return b.safeSubscribe(b.config.resultBusKey, channelResult, func(c *emitter.Client, m emitter.Message) {
		var msg ResultMsg
		if err := json.Unmarshal(m.Payload(), &msg); err != nil {
			log.Println(err)
			return
		}
		go f(&msg)
	})
}
