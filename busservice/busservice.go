package busservice

import (
	"context"
	"log"

	"github.com/SecurityBrewery/catalyst/bus"
	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/models"
	"github.com/SecurityBrewery/catalyst/role"
)

type busService struct {
	db          *database.Database
	apiURL      string
	apiKey      string
	catalystBus *bus.Bus
}

func New(apiurl, apikey string, catalystBus *bus.Bus, db *database.Database) error {

	h := &busService{db: db, apiURL: apiurl, apiKey: apikey, catalystBus: catalystBus}

	if err := catalystBus.SubscribeRequest(h.logRequest); err != nil {
		return err
	}
	if err := catalystBus.SubscribeResult(h.handleResult); err != nil {
		return err
	}
	if err := catalystBus.SubscribeJob(h.handleJob); err != nil {
		return err
	}

	return nil
}

func busContext() context.Context {
	// TODO: change roles?
	bot := &models.UserResponse{ID: "bot", Roles: []string{role.Admin}}
	return busdb.UserContext(context.Background(), bot)
}

func (h *busService) logRequest(msg *bus.RequestMsg) {
	var logEntries []*models.LogEntry
	for _, i := range msg.IDs {
		logEntries = append(logEntries, &models.LogEntry{Reference: i.String()})
	}

	if err := h.db.LogBatchCreate(busContext(), logEntries); err != nil {
		log.Println(err)
	}
}
