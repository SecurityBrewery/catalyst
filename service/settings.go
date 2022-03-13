package service

import (
	"context"
	"errors"
	"sort"

	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/model"
	"github.com/SecurityBrewery/catalyst/role"
)

func (s *Service) GetSettings(ctx context.Context) (*model.SettingsResponse, error) {
	globalSettings, err := s.database.Settings(ctx)
	if err != nil {
		return nil, err
	}

	return s.settings(ctx, globalSettings)
}

func (s *Service) SaveSettings(ctx context.Context, settings *model.Settings) (*model.SettingsResponse, error) {
	globalSettings, err := s.database.SaveSettings(ctx, settings)
	if err != nil {
		return nil, err
	}

	return s.settings(ctx, globalSettings)
}

func (s *Service) settings(ctx context.Context, globalSettings *model.Settings) (*model.SettingsResponse, error) {
	user, ok := busdb.UserFromContext(ctx)
	if !ok {
		return nil, errors.New("no user in context")
	}

	userData, err := s.database.UserDataGet(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	ticketTypeList, err := s.database.TicketTypeList(ctx)
	if err != nil {
		return nil, err
	}

	if userData.Timeformat != nil {
		globalSettings.Timeformat = *userData.Timeformat
	}
	roles := role.Strings(role.List())
	sort.Strings(roles)

	return &model.SettingsResponse{
		Tier:           model.SettingsResponseTierCommunity,
		Version:        s.version,
		Roles:          roles,
		TicketTypes:    ticketTypeList,
		ArtifactStates: globalSettings.ArtifactStates,
		ArtifactKinds:  globalSettings.ArtifactKinds,
		Timeformat:     globalSettings.Timeformat,
	}, nil
}
