package busservice

import (
	"context"
	"log"

	maut "github.com/cugu/maut/auth"

	"github.com/SecurityBrewery/catalyst/bus"
	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/generated/model"
	"github.com/SecurityBrewery/catalyst/generated/time"
)

type busService struct {
	db          *database.Database
	apiURL      string
	apiKey      string
	catalystBus *bus.Bus
	network     string
}

func New(apiURL, apikey, network string, catalystBus *bus.Bus, db *database.Database) {
	h := &busService{db: db, apiURL: apiURL, apiKey: apikey, network: network, catalystBus: catalystBus}

	catalystBus.RequestChannel.Subscribe(h.logRequest)
	catalystBus.ResultChannel.Subscribe(h.handleResult)
	catalystBus.JobChannel.Subscribe(h.handleJob)
}

func busContext() context.Context {
	// TODO: change roles?
	bot := &maut.User{ID: "bot", Roles: []string{maut.AdminRole}}

	return maut.UserContext(context.Background(), bot, nil) // TODO add permissions ?
}

func (h *busService) logRequest(msg *bus.RequestMsg) {
	var logEntries []*model.LogEntry
	for _, i := range msg.IDs {
		logEntries = append(logEntries, &model.LogEntry{
			Type:      "request",
			Reference: i.String(),
			Creator:   msg.User,
			Message:   msg.Function,
			Created:   time.Now().UTC(),
		})
	}

	if err := h.db.LogBatchCreate(busContext(), logEntries); err != nil {
		log.Println(err)
	}
}
