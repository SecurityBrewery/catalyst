package service

import (
	"context"

	"github.com/SecurityBrewery/catalyst/generated/restapi/api"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/templates"
)

func (s *Service) CreateTemplate(ctx context.Context, params *templates.CreateTemplateParams) *api.Response {
	return response(s.database.TemplateCreate(ctx, params.Template))
}

func (s *Service) GetTemplate(ctx context.Context, params *templates.GetTemplateParams) *api.Response {
	return response(s.database.TemplateGet(ctx, params.ID))
}

func (s *Service) UpdateTemplate(ctx context.Context, params *templates.UpdateTemplateParams) *api.Response {
	return response(s.database.TemplateUpdate(ctx, params.ID, params.Template))
}

func (s *Service) DeleteTemplate(ctx context.Context, params *templates.DeleteTemplateParams) *api.Response {
	return response(nil, s.database.TemplateDelete(ctx, params.ID))
}

func (s *Service) ListTemplates(ctx context.Context) *api.Response {
	return response(s.database.TemplateList(ctx))
}
