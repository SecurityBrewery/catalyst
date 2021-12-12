package service

import (
	"context"
	"fmt"

	"github.com/SecurityBrewery/catalyst/generated/models"
	"github.com/SecurityBrewery/catalyst/generated/restapi/api"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/tickets"
)

func (s *Service) CreateTicket(ctx context.Context, params *tickets.CreateTicketParams) *api.Response {
	createdTickets, err := s.database.TicketBatchCreate(ctx, []*models.TicketForm{params.Ticket})
	if len(createdTickets) > 0 {
		return response(createdTickets[0], err)
	}
	return response(nil, err)
}

func (s *Service) CreateTicketBatch(ctx context.Context, params *tickets.CreateTicketBatchParams) *api.Response {
	_, err := s.database.TicketBatchCreate(ctx, params.Ticket)
	return response(nil, err)
}

func (s *Service) GetTicket(ctx context.Context, params *tickets.GetTicketParams) *api.Response {
	return response(s.database.TicketGet(ctx, params.ID))
}

func (s *Service) UpdateTicket(ctx context.Context, params *tickets.UpdateTicketParams) *api.Response {
	return response(s.database.TicketUpdate(ctx, params.ID, params.Ticket))
}

func (s *Service) DeleteTicket(ctx context.Context, params *tickets.DeleteTicketParams) *api.Response {
	if err := s.database.TicketDelete(ctx, params.ID); err != nil {
		return response(nil, err)
	}

	_ = s.storage.DeleteBucket(fmt.Sprint(params.ID))
	return response(nil, nil)
}

func (s *Service) ListTickets(ctx context.Context, params *tickets.ListTicketsParams) *api.Response {
	q := ""
	if params.Query != nil && *params.Query != "" {
		q = *params.Query
	}
	t := ""
	if params.Type != nil && *params.Type != "" {
		t = *params.Type
	}
	return response(s.database.TicketList(ctx, t, q, params.Sort, params.Desc, *params.Offset, *params.Count))
}
