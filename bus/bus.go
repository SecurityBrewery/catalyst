package bus

import (
	"encoding/json"
	"log"

	emitter "github.com/emitter-io/go/v2"
)

type Bus struct {
	config *Config
	client *emitter.Client
}

type Config struct {
	Host                 string
	Key                  string
	databaseUpdateBusKey string
	jobBusKey            string
	resultBusKey         string
	requestKey           string
	APIUrl               string
}

func New(c *Config) (*Bus, error) {
	client, err := emitter.Connect(c.Host, func(_ *emitter.Client, msg emitter.Message) {
		log.Printf("received: '%s' topic: '%s'\n", msg.Payload(), msg.Topic())
	})
	if err != nil {
		return nil, err
	}

	c.databaseUpdateBusKey, err = client.GenerateKey(c.Key, channelDatabaseUpdate+"/", "rwls", 0)
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
	c.requestKey, err = client.GenerateKey(c.Key, ChannelRequest+"/", "rwls", 0)
	if err != nil {
		return nil, err
	}

	return &Bus{config: c, client: client}, err
}

func (b *Bus) jsonPublish(msg any, channel, key string) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return b.client.Publish(key, channel, payload)
}

func (b *Bus) safeSubscribe(key, channel string, handler func(c *emitter.Client, m emitter.Message)) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered %s in channel %s\n", r, channel)
		}
	}()

	return b.client.Subscribe(key, channel, handler)
}
