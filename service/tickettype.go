package service

import (
	"context"

	"github.com/SecurityBrewery/catalyst/generated/restapi/api"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/tickettypes"
)

func (s *Service) CreateTicketType(ctx context.Context, params *tickettypes.CreateTicketTypeParams) *api.Response {
	return response(s.database.TicketTypeCreate(ctx, params.Tickettype))
}

func (s *Service) GetTicketType(ctx context.Context, params *tickettypes.GetTicketTypeParams) *api.Response {
	return response(s.database.TicketTypeGet(ctx, params.ID))
}

func (s *Service) UpdateTicketType(ctx context.Context, params *tickettypes.UpdateTicketTypeParams) *api.Response {
	return response(s.database.TicketTypeUpdate(ctx, params.ID, params.Tickettype))
}

func (s *Service) DeleteTicketType(ctx context.Context, params *tickettypes.DeleteTicketTypeParams) *api.Response {
	return response(nil, s.database.TicketTypeDelete(ctx, params.ID))
}

func (s *Service) ListTicketTypes(ctx context.Context) *api.Response {
	return response(s.database.TicketTypeList(ctx))
}
