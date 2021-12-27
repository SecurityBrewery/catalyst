package service

import (
	"context"

	"github.com/SecurityBrewery/catalyst/generated/restapi/api"
)

func (s *Service) GetStatistics(ctx context.Context) *api.Response {
	i, err := s.database.Statistics(ctx)
	return s.response("GetStatistics", nil, i, err)
}
