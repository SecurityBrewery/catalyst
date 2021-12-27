package service

import (
	"context"

	"github.com/SecurityBrewery/catalyst/generated/restapi/api"
)

func (s *Service) ListTasks(ctx context.Context) *api.Response {
	i, err := s.database.TaskList(ctx)
	return s.response("ListTasks", nil, i, err)
}
