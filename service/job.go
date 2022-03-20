package service

import (
	"context"
	"fmt"

	"github.com/arangodb/go-driver"
	"github.com/google/uuid"

	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/generated/model"
)

func jobResponseID(job *model.JobResponse) []driver.DocumentID {
	if job == nil {
		return nil
	}

	return jobID(job.ID)
}

func jobID(id string) []driver.DocumentID {
	return []driver.DocumentID{driver.DocumentID(fmt.Sprintf("%s/%s", database.JobCollectionName, id))}
}

func (s *Service) ListJobs(ctx context.Context) ([]*model.JobResponse, error) {
	return s.database.JobList(ctx)
}

func (s *Service) RunJob(ctx context.Context, form *model.JobForm) (doc *model.JobResponse, err error) {
	msgContext := &model.Context{}
	newJobID := uuid.NewString()

	defer s.publishRequest(ctx, err, "RunJob", jobID(newJobID))
	err = s.bus.PublishJob(newJobID, form.Automation, form.Payload, msgContext, form.Origin)

	return &model.JobResponse{
		Automation: form.Automation,
		ID:         newJobID,
		Origin:     form.Origin,
		Payload:    form.Payload,
		Status:     "published",
	}, err
}

func (s *Service) GetJob(ctx context.Context, id string) (*model.JobResponse, error) {
	return s.database.JobGet(ctx, id)
}

func (s *Service) UpdateJob(ctx context.Context, id string, job *model.JobUpdate) (doc *model.JobResponse, err error) {
	defer s.publishRequest(ctx, err, "UpdateJob", jobResponseID(doc))

	return s.database.JobUpdate(ctx, id, job)
}
