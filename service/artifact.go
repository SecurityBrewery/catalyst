package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/SecurityBrewery/catalyst/generated/models"
	"github.com/SecurityBrewery/catalyst/generated/restapi/api"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/tickets"
)

func (s *Service) RunArtifact(ctx context.Context, params *tickets.RunArtifactParams) (r *api.Response) {
	artifact, err := s.database.ArtifactGet(ctx, params.ID, params.Name)
	if err != nil {
		return s.response(ctx, "RunArtifact", ticketID(params.ID), nil, err)
	}

	jobID := uuid.NewString()
	origin := &models.Origin{ArtifactOrigin: &models.ArtifactOrigin{TicketId: params.ID, Artifact: params.Name}}
	err = s.bus.PublishJob(jobID, params.Automation, params.Name, &models.Context{Artifact: artifact}, origin)
	return s.response(ctx, "RunArtifact", ticketID(params.ID), nil, err)
}

func (s *Service) EnrichArtifact(ctx context.Context, params *tickets.EnrichArtifactParams) *api.Response {
	i, err := s.database.EnrichArtifact(ctx, params.ID, params.Name, params.Data)
	return s.response(ctx, "EnrichArtifact", ticketID(params.ID), i, err)
}

func (s *Service) SetArtifact(ctx context.Context, params *tickets.SetArtifactParams) *api.Response {
	i, err := s.database.ArtifactUpdate(ctx, params.ID, params.Name, params.Artifact)
	return s.response(ctx, "SetArtifact", ticketID(params.ID), i, err)
}

func (s *Service) GetArtifact(ctx context.Context, params *tickets.GetArtifactParams) *api.Response {
	i, err := s.database.ArtifactGet(ctx, params.ID, params.Name)
	return s.response(ctx, "GetArtifact", nil, i, err)
}
