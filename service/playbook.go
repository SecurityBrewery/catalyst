package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/arangodb/go-driver"
	"github.com/xeipuuv/gojsonschema"

	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/generated/models"
	"github.com/SecurityBrewery/catalyst/generated/restapi/api"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/playbooks"
)

func playbookID(id string) []driver.DocumentID {
	return []driver.DocumentID{driver.DocumentID(fmt.Sprintf("%s/%s", database.PlaybookCollectionName, id))}
}

func (s *Service) CreatePlaybook(ctx context.Context, params *playbooks.CreatePlaybookParams) *api.Response {
	i, err := s.database.PlaybookCreate(ctx, params.Playbook)
	return s.response("CreatePlaybook", playbookID(i.ID), i, err)
}

func (s *Service) GetPlaybook(ctx context.Context, params *playbooks.GetPlaybookParams) *api.Response {
	i, err := s.database.PlaybookGet(ctx, params.ID)
	return s.response("GetPlaybook", nil, i, err)
}

func (s *Service) UpdatePlaybook(ctx context.Context, params *playbooks.UpdatePlaybookParams) *api.Response {
	if err := validate(params.Playbook, models.PlaybookTemplateFormSchema); err != nil {
		return s.response("UpdatePlaybook", nil, nil, err)
	}

	i, err := s.database.PlaybookUpdate(ctx, params.ID, params.Playbook)
	return s.response("UpdatePlaybook", playbookID(i.ID), i, err)
}

func (s *Service) DeletePlaybook(ctx context.Context, params *playbooks.DeletePlaybookParams) *api.Response {
	err := s.database.PlaybookDelete(ctx, params.ID)
	return s.response("DeletePlaybook", playbookID(params.ID), nil, err)
}

func (s *Service) ListPlaybooks(ctx context.Context) *api.Response {
	i, err := s.database.PlaybookList(ctx)
	return s.response("ListPlaybooks", nil, i, err)
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
