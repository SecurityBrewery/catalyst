package service

import (
	"context"
	"fmt"

	"github.com/arangodb/go-driver"
	"github.com/google/uuid"

	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/generated/models"
	"github.com/SecurityBrewery/catalyst/generated/restapi/api"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/jobs"
)

func jobID(id string) []driver.DocumentID {
	return []driver.DocumentID{driver.DocumentID(fmt.Sprintf("%s/%s", database.JobCollectionName, id))}
}

func (s *Service) RunJob(ctx context.Context, params *jobs.RunJobParams) *api.Response {
	msgContext := &models.Context{}
	newJobID := uuid.NewString()
	return s.response(ctx, "RunJob", jobID(newJobID), nil, s.bus.PublishJob(newJobID, params.Job.Automation, params.Job.Payload, msgContext, params.Job.Origin))
}

func (s *Service) GetJob(ctx context.Context, params *jobs.GetJobParams) *api.Response {
	i, err := s.database.JobGet(ctx, params.ID)
	return s.response(ctx, "GetJob", nil, i, err)
}

func (s *Service) ListJobs(ctx context.Context) *api.Response {
	i, err := s.database.JobList(ctx)
	return s.response(ctx, "ListJobs", nil, i, err)
}

func (s *Service) UpdateJob(ctx context.Context, params *jobs.UpdateJobParams) *api.Response {
	i, err := s.database.JobUpdate(ctx, params.ID, params.Job)
	return s.response(ctx, "UpdateJob", jobID(i.ID), i, err)
}
