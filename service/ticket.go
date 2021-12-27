package service

import (
	"context"
	"fmt"

	"github.com/arangodb/go-driver"

	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/generated/models"
	"github.com/SecurityBrewery/catalyst/generated/restapi/api"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/tickets"
)

func ticketID(id int64) []driver.DocumentID {
	return []driver.DocumentID{driver.DocumentID(fmt.Sprintf("%s/%d", database.TicketCollectionName, id))}
}

func ticketIDs(ticketResponses []*models.TicketResponse) []driver.DocumentID {
	var ids []driver.DocumentID
	for _, ticketResponse := range ticketResponses {
		ids = append(ids, ticketID(ticketResponse.ID)...)
	}
	return ids
}

func (s *Service) CreateTicket(ctx context.Context, params *tickets.CreateTicketParams) *api.Response {
	createdTickets, err := s.database.TicketBatchCreate(ctx, []*models.TicketForm{params.Ticket})
	if len(createdTickets) > 0 {
		return s.response("CreateTicket", ticketIDs(createdTickets), createdTickets[0], err)
	}
	return s.response("CreateTicket", ticketIDs(createdTickets), nil, err)
}

func (s *Service) CreateTicketBatch(ctx context.Context, params *tickets.CreateTicketBatchParams) *api.Response {
	ticketBatch, err := s.database.TicketBatchCreate(ctx, params.Ticket)
	return s.response("CreateTicketBatch", ticketIDs(ticketBatch), nil, err)
}

func (s *Service) GetTicket(ctx context.Context, params *tickets.GetTicketParams) *api.Response {
	ticket, err := s.database.TicketGet(ctx, params.ID)
	return s.response("GetTicket", nil, ticket, err)
}

func (s *Service) UpdateTicket(ctx context.Context, params *tickets.UpdateTicketParams) *api.Response {
	ticket, err := s.database.TicketUpdate(ctx, params.ID, params.Ticket)
	return s.response("UpdateTicket", ticketID(ticket.ID), ticket, err)
}

func (s *Service) DeleteTicket(ctx context.Context, params *tickets.DeleteTicketParams) *api.Response {
	if err := s.database.TicketDelete(ctx, params.ID); err != nil {
		return s.response("DeleteTicket", ticketID(params.ID), nil, err)
	}

	_ = s.storage.DeleteBucket(fmt.Sprint(params.ID))
	return s.response("DeleteTicket", ticketID(params.ID), nil, nil)
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

	ticketList, err := s.database.TicketList(ctx, t, q, params.Sort, params.Desc, *params.Offset, *params.Count)
	return s.response("ListTickets", nil, ticketList, err)
}
