package service

import (
	"context"
	"net/url"

	"github.com/SecurityBrewery/catalyst/generated/restapi/api"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/logs"
)

func (s *Service) GetLogs(ctx context.Context, params *logs.GetLogsParams) *api.Response {
	id, _ := url.QueryUnescape(params.Reference)
	return response(s.database.LogList(ctx, id))
}
