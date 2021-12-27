package service

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/arangodb/go-driver"
	"github.com/gin-gonic/gin"

	"github.com/SecurityBrewery/catalyst/bus"
	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/models"
	"github.com/SecurityBrewery/catalyst/generated/restapi/api"
	"github.com/SecurityBrewery/catalyst/storage"
)

type Service struct {
	bus      *bus.Bus
	database *database.Database
	settings *models.Settings
	storage  *storage.Storage
}

func New(bus *bus.Bus, database *database.Database, storage *storage.Storage, settings *models.Settings) (*Service, error) {
	return &Service{database: database, bus: bus, settings: settings, storage: storage}, nil
}

func (s *Service) Healthy() bool {
	return true
}

func (s *Service) response(ctx context.Context, function string, ids []driver.DocumentID, v interface{}, err error) *api.Response {
	if err != nil {
		log.Println(err)
		return &api.Response{Code: httpStatus(err), Body: gin.H{"error": err.Error()}}
	}

	if ids != nil {
		userID := "unknown"
		user, ok := busdb.UserFromContext(ctx)
		if ok {
			userID = user.ID
		}

		go s.bus.PublishRequest(userID, function, ids)
	}

	if v == nil {
		return &api.Response{Code: http.StatusNoContent, Body: v}
	}
	return &api.Response{Code: http.StatusOK, Body: v}
}

func httpStatus(err error) int {
	ae := driver.ArangoError{}
	if errors.As(err, &ae) {
		return ae.Code
	}

	return http.StatusInternalServerError
}
