package service

import (
	"context"

	"github.com/SecurityBrewery/catalyst/generated/restapi/api"
)

func (s *Service) ListTasks(ctx context.Context) *api.Response {
	return response(s.database.TaskList(ctx))
}
