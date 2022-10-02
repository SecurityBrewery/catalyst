package service

import (
	"context"

	maut "github.com/jonas-plum/maut/auth"

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
	if ok {
		userData, _ := s.database.UserDataGet(ctx, user.ID)

		if userData != nil && userData.Timeformat != nil {
			globalSettings.Timeformat = *userData.Timeformat
		}
	}

	ticketTypeList, _ := s.database.TicketTypeList(ctx)

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
