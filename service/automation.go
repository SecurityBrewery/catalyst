package service

import (
	"context"

	"github.com/SecurityBrewery/catalyst/generated/restapi/api"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/automations"
)

func (s *Service) CreateAutomation(ctx context.Context, params *automations.CreateAutomationParams) *api.Response {
	return response(s.database.AutomationCreate(ctx, params.Automation))
}

func (s *Service) GetAutomation(ctx context.Context, params *automations.GetAutomationParams) *api.Response {
	return response(s.database.AutomationGet(ctx, params.ID))
}

func (s *Service) UpdateAutomation(ctx context.Context, params *automations.UpdateAutomationParams) *api.Response {
	return response(s.database.AutomationUpdate(ctx, params.ID, params.Automation))
}

func (s *Service) DeleteAutomation(ctx context.Context, params *automations.DeleteAutomationParams) *api.Response {
	return response(nil, s.database.AutomationDelete(ctx, params.ID))
}

func (s *Service) ListAutomations(ctx context.Context) *api.Response {
	return response(s.database.AutomationList(ctx))
}
