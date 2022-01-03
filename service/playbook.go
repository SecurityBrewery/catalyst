package service

import (
	"context"
	"fmt"

	"github.com/arangodb/go-driver"

	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/generated/model"
)

func playbookResponseID(playbook *model.PlaybookTemplateResponse) []driver.DocumentID {
	if playbook == nil {
		return nil
	}
	return playbookID(playbook.ID)
}

func playbookID(id string) []driver.DocumentID {
	return []driver.DocumentID{driver.DocumentID(fmt.Sprintf("%s/%s", database.PlaybookCollectionName, id))}
}

func (s *Service) ListPlaybooks(ctx context.Context) ([]*model.PlaybookTemplateResponse, error) {
	return s.database.PlaybookList(ctx)
}

func (s *Service) CreatePlaybook(ctx context.Context, form *model.PlaybookTemplateForm) (doc *model.PlaybookTemplateResponse, err error) {
	defer s.publishRequest(ctx, err, "CreatePlaybook", playbookResponseID(doc))
	return s.database.PlaybookCreate(ctx, form)
}

func (s *Service) GetPlaybook(ctx context.Context, id string) (*model.PlaybookTemplateResponse, error) {
	return s.database.PlaybookGet(ctx, id)
}

func (s *Service) UpdatePlaybook(ctx context.Context, id string, form *model.PlaybookTemplateForm) (doc *model.PlaybookTemplateResponse, err error) {
	defer s.publishRequest(ctx, err, "UpdatePlaybook", playbookResponseID(doc))
	return s.database.PlaybookUpdate(ctx, id, form)
}

func (s *Service) DeletePlaybook(ctx context.Context, id string) (err error) {
	defer s.publishRequest(ctx, err, "DeletePlaybook", playbookID(id))
	return s.database.PlaybookDelete(ctx, id)
}
