package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"

	"github.com/SecurityBrewery/catalyst/generated/models"
	"github.com/SecurityBrewery/catalyst/generated/restapi/api"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/playbooks"
)

func (s *Service) CreatePlaybook(ctx context.Context, params *playbooks.CreatePlaybookParams) *api.Response {
	return response(s.database.PlaybookCreate(ctx, params.Playbook))
}

func (s *Service) GetPlaybook(ctx context.Context, params *playbooks.GetPlaybookParams) *api.Response {
	return response(s.database.PlaybookGet(ctx, params.ID))
}

func (s *Service) UpdatePlaybook(ctx context.Context, params *playbooks.UpdatePlaybookParams) *api.Response {
	if err := validate(params.Playbook, models.PlaybookTemplateFormSchema); err != nil {
		return response(nil, err)
	}

	return response(s.database.PlaybookUpdate(ctx, params.ID, params.Playbook))
}

func (s *Service) DeletePlaybook(ctx context.Context, params *playbooks.DeletePlaybookParams) *api.Response {
	return response(nil, s.database.PlaybookDelete(ctx, params.ID))
}

func (s *Service) ListPlaybooks(ctx context.Context) *api.Response {
	return response(s.database.PlaybookList(ctx))
}

func validate(e interface{}, schema *gojsonschema.Schema) error {
	res, err := schema.Validate(gojsonschema.NewGoLoader(e))
	if err != nil {
		return err
	}

	if len(res.Errors()) > 0 {
		var l []string
		for _, e := range res.Errors() {
			l = append(l, e.String())
		}
		return fmt.Errorf("validation failed: %v", strings.Join(l, ", "))
	}
	return nil
}
