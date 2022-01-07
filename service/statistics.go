package service

import (
	"context"

	"github.com/SecurityBrewery/catalyst/generated/model"
)

func (s *Service) GetStatistics(ctx context.Context) (*model.Statistics, error) {
	return s.database.Statistics(ctx)
}
