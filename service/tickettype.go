package service

import (
	"context"
	"fmt"

	"github.com/arangodb/go-driver"

	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/generated/restapi/api"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/tickettypes"
)

func ticketTypeID(id string) []driver.DocumentID {
	return []driver.DocumentID{driver.DocumentID(fmt.Sprintf("%s/%s", database.TicketTypeCollectionName, id))}
}

func (s *Service) CreateTicketType(ctx context.Context, params *tickettypes.CreateTicketTypeParams) *api.Response {
	ticketType, err := s.database.TicketTypeCreate(ctx, params.Tickettype)
	return s.response("CreateTicketType", ticketTypeID(ticketType.ID), ticketType, err)
}

func (s *Service) GetTicketType(ctx context.Context, params *tickettypes.GetTicketTypeParams) *api.Response {
	ticketType, err := s.database.TicketTypeGet(ctx, params.ID)
	return s.response("GetTicketType", nil, ticketType, err)
}

func (s *Service) UpdateTicketType(ctx context.Context, params *tickettypes.UpdateTicketTypeParams) *api.Response {
	ticketType, err := s.database.TicketTypeUpdate(ctx, params.ID, params.Tickettype)
	return s.response("UpdateTicketType", ticketTypeID(ticketType.ID), ticketType, err)
}

func (s *Service) DeleteTicketType(ctx context.Context, params *tickettypes.DeleteTicketTypeParams) *api.Response {
	err := s.database.TicketTypeDelete(ctx, params.ID)
	return s.response("DeleteTicketType", ticketTypeID(params.ID), nil, err)
}

func (s *Service) ListTicketTypes(ctx context.Context) *api.Response {
	ticketTypes, err := s.database.TicketTypeList(ctx)
	return s.response("ListTicketTypes", nil, ticketTypes, err)
}
