package service

import (
	"context"
	"errors"
	"sort"

	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/models"
	"github.com/SecurityBrewery/catalyst/generated/restapi/api"
	"github.com/SecurityBrewery/catalyst/role"
)

func (s *Service) GetSettings(ctx context.Context) *api.Response {
	user, ok := busdb.UserFromContext(ctx)
	if !ok {
		return s.response("GetSettings", nil, nil, errors.New("no user in context"))
	}

	setting, err := s.database.UserDataGet(ctx, user.ID)
	if err != nil {
		return s.response("GetSettings", nil, nil, err)
	}

	settings := mergeSettings(s.settings, setting)

	ticketTypeList, err := s.database.TicketTypeList(ctx)
	if err != nil {
		return s.response("GetSettings", nil, nil, err)
	}

	settings.TicketTypes = ticketTypeList

	return s.response("GetSettings", nil, settings, nil)
}

func mergeSettings(globalSettings *models.Settings, user *models.UserDataResponse) *models.Settings {
	if user.Timeformat != nil {
		globalSettings.Timeformat = *user.Timeformat
	}
	roles := role.Strings(role.List())
	sort.Strings(roles)
	globalSettings.Roles = roles

	return globalSettings
}
