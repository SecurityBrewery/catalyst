package service

import (
	"context"
	"fmt"

	"github.com/arangodb/go-driver"

	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/generated/model"
)

func automationResponseID(automation *model.AutomationResponse) []driver.DocumentID {
	if automation == nil {
		return nil
	}
	return automationID(automation.ID)
}

func automationID(id string) []driver.DocumentID {
	return []driver.DocumentID{driver.DocumentID(fmt.Sprintf("%s/%s", database.AutomationCollectionName, id))}
}

func (s *Service) ListAutomations(ctx context.Context) ([]*model.AutomationResponse, error) {
	return s.database.AutomationList(ctx)
}

func (s *Service) CreateAutomation(ctx context.Context, form *model.AutomationForm) (doc *model.AutomationResponse, err error) {
	defer s.publishRequest(ctx, err, "CreateAutomation", automationResponseID(doc))
	return s.database.AutomationCreate(ctx, form)
}

func (s *Service) GetAutomation(ctx context.Context, id string) (*model.AutomationResponse, error) {
	return s.database.AutomationGet(ctx, id)
}

func (s *Service) UpdateAutomation(ctx context.Context, id string, form *model.AutomationForm) (doc *model.AutomationResponse, err error) {
	defer s.publishRequest(ctx, err, "UpdateAutomation", automationResponseID(doc))
	return s.database.AutomationUpdate(ctx, id, form)
}

func (s *Service) DeleteAutomation(ctx context.Context, id string) (err error) {
	defer s.publishRequest(ctx, err, "DeleteAutomation", automationID(id))
	return s.database.AutomationDelete(ctx, id)
}
