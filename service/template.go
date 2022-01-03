package service

import (
	"context"
	"fmt"

	"github.com/arangodb/go-driver"

	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/generated/model"
)

func templateResponseID(template *model.TicketTemplateResponse) []driver.DocumentID {
	if template == nil {
		return nil
	}
	return templateID(template.ID)
}

func templateID(id string) []driver.DocumentID {
	return []driver.DocumentID{driver.DocumentID(fmt.Sprintf("%s/%s", database.TemplateCollectionName, id))}
}

func (s *Service) ListTemplates(ctx context.Context) ([]*model.TicketTemplateResponse, error) {
	return s.database.TemplateList(ctx)
}

func (s *Service) CreateTemplate(ctx context.Context, form *model.TicketTemplateForm) (doc *model.TicketTemplateResponse, err error) {
	defer s.publishRequest(ctx, err, "CreateTemplate", templateResponseID(doc))
	return s.database.TemplateCreate(ctx, form)
}

func (s *Service) GetTemplate(ctx context.Context, id string) (*model.TicketTemplateResponse, error) {
	return s.database.TemplateGet(ctx, id)
}

func (s *Service) UpdateTemplate(ctx context.Context, id string, form *model.TicketTemplateForm) (doc *model.TicketTemplateResponse, err error) {
	defer s.publishRequest(ctx, err, "UpdateTemplate", templateResponseID(doc))
	return s.database.TemplateUpdate(ctx, id, form)
}

func (s *Service) DeleteTemplate(ctx context.Context, id string) (err error) {
	defer s.publishRequest(ctx, err, "DeleteTemplate", templateID(id))
	return s.database.TemplateDelete(ctx, id)
}
