package service

import (
	"context"
	"fmt"

	"github.com/arangodb/go-driver"

	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/generated/model"
)

func ticketTypeResponseID(ticketType *model.TicketTypeResponse) []driver.DocumentID {
	if ticketType == nil {
		return nil
	}
	return userDataID(ticketType.ID)
}

func ticketTypeID(id string) []driver.DocumentID {
	return []driver.DocumentID{driver.DocumentID(fmt.Sprintf("%s/%s", database.UserDataCollectionName, id))}
}

func (s *Service) ListTicketTypes(ctx context.Context) ([]*model.TicketTypeResponse, error) {
	return s.database.TicketTypeList(ctx)
}

func (s *Service) CreateTicketType(ctx context.Context, form *model.TicketTypeForm) (doc *model.TicketTypeResponse, err error) {
	defer s.publishRequest(ctx, err, "CreateTicketType", ticketTypeResponseID(doc))
	return s.database.TicketTypeCreate(ctx, form)
}

func (s *Service) GetTicketType(ctx context.Context, id string) (*model.TicketTypeResponse, error) {
	return s.database.TicketTypeGet(ctx, id)
}

func (s *Service) UpdateTicketType(ctx context.Context, id string, form *model.TicketTypeForm) (doc *model.TicketTypeResponse, err error) {
	defer s.publishRequest(ctx, err, "UpdateTicketType", ticketTypeResponseID(doc))
	return s.database.TicketTypeUpdate(ctx, id, form)
}

func (s *Service) DeleteTicketType(ctx context.Context, id string) (err error) {
	defer s.publishRequest(ctx, err, "DeleteTicketType", ticketTypeID(id))
	return s.database.TicketTypeDelete(ctx, id)
}
