package service

import (
	"context"
	"fmt"

	"github.com/arangodb/go-driver"

	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/generated/restapi/api"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/automations"
)

func automationID(id string) []driver.DocumentID {
	return []driver.DocumentID{driver.DocumentID(fmt.Sprintf("%s/%s", database.AutomationCollectionName, id))}
}

func (s *Service) CreateAutomation(ctx context.Context, params *automations.CreateAutomationParams) *api.Response {
	i, err := s.database.AutomationCreate(ctx, params.Automation)
	return s.response("CreateAutomation", automationID(i.ID), i, err)
}

func (s *Service) GetAutomation(ctx context.Context, params *automations.GetAutomationParams) *api.Response {
	i, err := s.database.AutomationGet(ctx, params.ID)
	return s.response("GetAutomation", nil, i, err)
}

func (s *Service) UpdateAutomation(ctx context.Context, params *automations.UpdateAutomationParams) *api.Response {
	i, err := s.database.AutomationUpdate(ctx, params.ID, params.Automation)
	return s.response("UpdateAutomation", automationID(i.ID), i, err)
}

func (s *Service) DeleteAutomation(ctx context.Context, params *automations.DeleteAutomationParams) *api.Response {
	err := s.database.AutomationDelete(ctx, params.ID)
	return s.response("DeleteAutomation", automationID(params.ID), nil, err)
}

func (s *Service) ListAutomations(ctx context.Context) *api.Response {
	i, err := s.database.AutomationList(ctx)
	return s.response("ListAutomations", nil, i, err)
}
