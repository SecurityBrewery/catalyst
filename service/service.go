package service

import (
	"context"

	"github.com/arangodb/go-driver"
	maut "github.com/cugu/maut/auth"

	"github.com/SecurityBrewery/catalyst/bus"
	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/storage"
)

type Service struct {
	bus      *bus.Bus
	database *database.Database
	storage  *storage.Storage
	version  string
}

func New(bus *bus.Bus, database *database.Database, storage *storage.Storage, version string) (*Service, error) {
	return &Service{database: database, bus: bus, storage: storage, version: version}, nil
}

func (s *Service) publishRequest(ctx context.Context, err error, function string, ids []driver.DocumentID) {
	if err != nil {
		return
	}
	if ids != nil {
		userID := "unknown"
		user, _, ok := maut.UserFromContext(ctx)
		if ok {
			userID = user.ID
		}

		s.bus.RequestChannel.Publish(&bus.RequestMsg{
			User:     userID,
			Function: function,
			IDs:      ids,
		})
	}
}
