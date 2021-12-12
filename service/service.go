package service

import (
	"errors"
	"log"
	"net/http"

	"github.com/arangodb/go-driver"
	"github.com/gin-gonic/gin"

	"github.com/SecurityBrewery/catalyst/bus"
	"github.com/SecurityBrewery/catalyst/database"
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

func response(v interface{}, err error) *api.Response {
	if err != nil {
		log.Println(err)
		return &api.Response{Code: httpStatus(err), Body: gin.H{"error": err.Error()}}
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
