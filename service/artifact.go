package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/SecurityBrewery/catalyst/generated/models"
	"github.com/SecurityBrewery/catalyst/generated/restapi/api"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/tickets"
)

func (s *Service) RunArtifact(ctx context.Context, params *tickets.RunArtifactParams) *api.Response {
	artifact, err := s.database.ArtifactGet(ctx, params.ID, params.Name)
	if err != nil {
		return response(nil, err)
	}

	jobID := uuid.NewString()
	origin := &models.Origin{ArtifactOrigin: &models.ArtifactOrigin{TicketId: params.ID, Artifact: params.Name}}
	return response(nil, s.bus.PublishJob(jobID, params.Automation, params.Name, &models.Context{Artifact: artifact}, origin))
}

func (s *Service) EnrichArtifact(ctx context.Context, params *tickets.EnrichArtifactParams) *api.Response {
	return response(s.database.EnrichArtifact(ctx, params.ID, params.Name, params.Data))
}

func (s *Service) SetArtifact(ctx context.Context, params *tickets.SetArtifactParams) *api.Response {
	return response(s.database.ArtifactUpdate(ctx, params.ID, params.Name, params.Artifact))
}

func (s *Service) GetArtifact(ctx context.Context, params *tickets.GetArtifactParams) *api.Response {
	return response(s.database.ArtifactGet(ctx, params.ID, params.Name))
}
