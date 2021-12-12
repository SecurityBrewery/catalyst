package bus

import (
	"encoding/json"
	"log"

	"github.com/arangodb/go-driver"
	emitter "github.com/emitter-io/go/v2"

	"github.com/SecurityBrewery/catalyst/generated/models"
)

const (
	channelUpdate = "data"
	channelJob    = "job"
	channelResult = "result"
)

type Bus struct {
	config *Config
	client *emitter.Client
}

type Config struct {
	Host         string
	Key          string
	resultBusKey string
	jobBusKey    string
	dataBusKey   string
	APIUrl       string
}

type JobMsg struct {
	ID         string          `json:"id"`
	Automation string          `json:"automation"`
	Origin     *models.Origin  `json:"origin"`
	Message    *models.Message `json:"message"`
}

type ResultMsg struct {
	Automation string                 `json:"automation"`
	Data       map[string]interface{} `json:"data,omitempty"`
	Target     *models.Origin         `json:"target"`
}

func New(c *Config) (*Bus, error) {
	client, err := emitter.Connect(c.Host, func(_ *emitter.Client, msg emitter.Message) {
		log.Printf("received: '%s' topic: '%s'\n", msg.Payload(), msg.Topic())
	})
	if err != nil {
		return nil, err
	}

	c.dataBusKey, err = client.GenerateKey(c.Key, channelUpdate+"/", "rwls", 0)
	if err != nil {
		return nil, err
	}
	c.jobBusKey, err = client.GenerateKey(c.Key, channelJob+"/", "rwls", 0)
	if err != nil {
		return nil, err
	}
	c.resultBusKey, err = client.GenerateKey(c.Key, channelResult+"/", "rwls", 0)
	if err != nil {
		return nil, err
	}

	return &Bus{config: c, client: client}, err
}

func (b *Bus) PublishUpdate(ids []driver.DocumentID) error {
	return b.jsonPublish(ids, channelUpdate, b.config.dataBusKey)
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

func (b *Bus) PublishResult(automation string, data map[string]interface{}, target *models.Origin) error {
	return b.jsonPublish(&ResultMsg{Automation: automation, Data: data, Target: target}, channelResult, b.config.resultBusKey)
}

func (b *Bus) jsonPublish(msg interface{}, channel, key string) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return b.client.Publish(key, channel, payload)
}

func (b *Bus) SubscribeUpdate(f func(ids []driver.DocumentID)) error {
	return b.safeSubscribe(b.config.dataBusKey, channelUpdate, func(c *emitter.Client, m emitter.Message) {
		var msg []driver.DocumentID
		if err := json.Unmarshal(m.Payload(), &msg); err != nil {
			log.Println(err)
			return
		}
		go f(msg)
	})
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

func (b *Bus) safeSubscribe(key, channel string, handler func(c *emitter.Client, m emitter.Message)) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered %s in channel %s\n", r, channel)
		}
	}()
	return b.client.Subscribe(key, channel, handler)
}
