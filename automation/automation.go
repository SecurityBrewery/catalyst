package automation

import (
	"context"
	"log"

	"github.com/SecurityBrewery/catalyst/bus"
	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/models"
	"github.com/SecurityBrewery/catalyst/role"
)

func New(apiurl, apikey string, bus *bus.Bus, db *database.Database) error {
	if err := jobAutomation(jobContext(), apiurl, apikey, bus, db); err != nil {
		log.Fatal(err)
	}

	return resultAutomation(bus, db)
}

func jobContext() context.Context {
	// TODO: change roles?
	bot := &models.UserResponse{ID: "bot", Roles: []string{role.Admin}}
	return busdb.UserContext(context.Background(), bot)
}
