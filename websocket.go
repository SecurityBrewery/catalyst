package catalyst

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/google/uuid"

	"github.com/SecurityBrewery/catalyst/bus"
	"github.com/SecurityBrewery/catalyst/generated/api"
)

type websocketBroker struct {
	clients map[string]chan []byte
	mu      sync.Mutex
}

func (wb *websocketBroker) Publish(b []byte) {
	for _, channel := range wb.clients {
		channel <- b
	}
}

func (wb *websocketBroker) CloseSocket(id string) {
	wb.mu.Lock()
	if channel, ok := wb.clients[id]; ok {
		close(channel)
		delete(wb.clients, id)
	}
	wb.mu.Unlock()
}

func (wb *websocketBroker) NewWebsocket() (string, chan []byte) {
	id := uuid.New().String()
	channel := make(chan []byte, 10)
	wb.mu.Lock()
	wb.clients[id] = channel
	wb.mu.Unlock()

	return id, channel
}

func handleWebSocket(catalystBus *bus.Bus) http.HandlerFunc {
	broker := websocketBroker{clients: map[string]chan []byte{}}

	// send all messages from bus to websocket
	err := catalystBus.SubscribeDatabaseUpdate(func(msg *bus.DatabaseUpdateMsg) {
		b, err := json.Marshal(map[string]any{
			"action": "update",
			"ids":    msg.IDs,
		})
		if err != nil {
			return
		}

		broker.Publish(b)
	})
	if err != nil {
		log.Println(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			api.JSONError(w, errors.New("upgrade failed"))

			return
		}

		go func() {
			defer conn.Close()

			id, messages := broker.NewWebsocket()
			for msg := range messages {
				if err := wsutil.WriteServerMessage(conn, ws.OpText, msg); err != nil {
					broker.CloseSocket(id)
				}
			}
		}()
	}
}
