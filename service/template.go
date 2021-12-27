package service

import (
	"context"
	"fmt"

	"github.com/arangodb/go-driver"

	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/generated/restapi/api"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/templates"
)

func templateID(s string) []driver.DocumentID {
	return []driver.DocumentID{driver.DocumentID(fmt.Sprintf("%s/%s", database.TemplateCollectionName, s))}
}

func (s *Service) CreateTemplate(ctx context.Context, params *templates.CreateTemplateParams) *api.Response {
	i, err := s.database.TemplateCreate(ctx, params.Template)
	return s.response(ctx, "CreateTemplate", templateID(i.ID), i, err)
}

func (s *Service) GetTemplate(ctx context.Context, params *templates.GetTemplateParams) *api.Response {
	i, err := s.database.TemplateGet(ctx, params.ID)
	return s.response(ctx, "GetTemplate", nil, i, err)
}

func (s *Service) UpdateTemplate(ctx context.Context, params *templates.UpdateTemplateParams) *api.Response {
	i, err := s.database.TemplateUpdate(ctx, params.ID, params.Template)
	return s.response(ctx, "UpdateTemplate", templateID(i.ID), i, err)
}

func (s *Service) DeleteTemplate(ctx context.Context, params *templates.DeleteTemplateParams) *api.Response {
	err := s.database.TemplateDelete(ctx, params.ID)
	return s.response(ctx, "DeleteTemplate", templateID(params.ID), nil, err)
}

func (s *Service) ListTemplates(ctx context.Context) *api.Response {
	i, err := s.database.TemplateList(ctx)
	return s.response(ctx, "ListTemplates", nil, i, err)
}
