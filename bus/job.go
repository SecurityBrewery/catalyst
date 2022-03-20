package bus

import (
	"encoding/json"
	"log"

	emitter "github.com/emitter-io/go/v2"

	"github.com/SecurityBrewery/catalyst/generated/model"
)

const channelJob = "job"

type JobMsg struct {
	ID         string         `json:"id"`
	Automation string         `json:"automation"`
	Origin     *model.Origin  `json:"origin"`
	Message    *model.Message `json:"message"`
}

func (b *Bus) PublishJob(id, automation string, payload any, context *model.Context, origin *model.Origin) error {
	return b.jsonPublish(&JobMsg{
		ID:         id,
		Automation: automation,
		Origin:     origin,
		Message: &model.Message{
			Context: context,
			Payload: payload,
		},
	}, channelJob, b.config.jobBusKey)
}

func (b *Bus) SubscribeJob(f func(msg *JobMsg)) error {
	return b.safeSubscribe(b.config.jobBusKey, channelJob, func(c *emitter.Client, m emitter.Message) {
		var msg JobMsg
		if err := json.Unmarshal(m.Payload(), &msg); err != nil {
			log.Println(err)

			return
		}
		go f(&msg)
	})
}
