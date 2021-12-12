package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/SecurityBrewery/catalyst/generated/models"
	"github.com/SecurityBrewery/catalyst/generated/restapi/api"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/jobs"
)

func (s *Service) RunJob(_ context.Context, params *jobs.RunJobParams) *api.Response {
	msgContext := &models.Context{}
	jobID := uuid.NewString()
	return response(nil, s.bus.PublishJob(jobID, params.Job.Automation, params.Job.Payload, msgContext, params.Job.Origin))
}

func (s *Service) GetJob(ctx context.Context, params *jobs.GetJobParams) *api.Response {
	return response(s.database.JobGet(ctx, params.ID))
}

func (s *Service) ListJobs(ctx context.Context) *api.Response {
	return response(s.database.JobList(ctx))
}

func (s *Service) UpdateJob(ctx context.Context, params *jobs.UpdateJobParams) *api.Response {
	return response(s.database.JobUpdate(ctx, params.ID, params.Job))
}
