package bus

import (
	"encoding/json"
	"log"

	emitter "github.com/emitter-io/go/v2"

	"github.com/SecurityBrewery/catalyst/generated/models"
)

const channelJob = "job"

type JobMsg struct {
	ID         string          `json:"id"`
	Automation string          `json:"automation"`
	Origin     *models.Origin  `json:"origin"`
	Message    *models.Message `json:"message"`
}

func (b *Bus) PublishJob(id, automation string, payload interface{}, context *models.Context, origin *models.Origin) error {
	return b.jsonPublish(&JobMsg{
		ID:         id,
		Automation: automation,
		Origin:     origin,
		Message: &models.Message{
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
