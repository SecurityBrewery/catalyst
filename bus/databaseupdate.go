package bus

import (
	"encoding/json"
	"log"

	"github.com/arangodb/go-driver"
	emitter "github.com/emitter-io/go/v2"
)

const channelDatabaseUpdate = "databaseupdate"

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

func (b *Bus) PublishDatabaseUpdate(ids []driver.DocumentID, databaseUpdateType DatabaseUpdateType) error {
	return b.jsonPublish(&DatabaseUpdateMsg{
		IDs:  ids,
		Type: databaseUpdateType,
	}, channelDatabaseUpdate, b.config.databaseUpdateBusKey)
}

func (b *Bus) SubscribeDatabaseUpdate(f func(msg *DatabaseUpdateMsg)) error {
	return b.safeSubscribe(b.config.databaseUpdateBusKey, channelDatabaseUpdate, func(c *emitter.Client, m emitter.Message) {
		var msg DatabaseUpdateMsg
		if err := json.Unmarshal(m.Payload(), &msg); err != nil {
			log.Println(err)

			return
		}
		go f(&msg)
	})
}
