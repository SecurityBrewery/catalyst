package service

import (
	"context"
	"errors"

	maut "github.com/cugu/maut/auth"

	"github.com/SecurityBrewery/catalyst/generated/model"
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
	user, permissions, ok := maut.UserFromContext(ctx)
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

	return &model.SettingsResponse{
		Tier:           model.SettingsResponseTierCommunity,
		Version:        s.version,
		Roles:          permissions,
		TicketTypes:    ticketTypeList,
		ArtifactStates: globalSettings.ArtifactStates,
		ArtifactKinds:  globalSettings.ArtifactKinds,
		Timeformat:     globalSettings.Timeformat,
	}, nil
}
