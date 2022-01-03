package service

import (
	"context"

	"github.com/SecurityBrewery/catalyst/generated/model"
)

func (s *Service) ListTasks(ctx context.Context) ([]*model.TaskWithContext, error) {
	return s.database.TaskList(ctx)
}
