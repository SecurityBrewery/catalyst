package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/arangodb/go-driver"
	"github.com/google/uuid"

	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/generated/api"
	"github.com/SecurityBrewery/catalyst/generated/model"
)

func ticketWithTicketsID(ticketResponse *model.TicketWithTickets) []driver.DocumentID {
	if ticketResponse == nil {
		return nil
	}
	return ticketID(ticketResponse.ID)
}

func ticketID(ticketID int64) []driver.DocumentID {
	id := fmt.Sprintf("%s/%d", database.TicketCollectionName, ticketID)
	return []driver.DocumentID{driver.DocumentID(id)}
}

func ticketIDs(ticketResponses []*model.TicketResponse) []driver.DocumentID {
	var ids []driver.DocumentID
	for _, ticketResponse := range ticketResponses {
		ids = append(ids, ticketID(ticketResponse.ID)...)
	}
	return ids
}

func (s *Service) ListTickets(ctx context.Context, s3 *string, i *int, i2 *int, strings []string, bools []bool, s2 *string) (*model.TicketList, error) {
	q := ""
	if s2 != nil && *s2 != "" {
		q = *s2
	}
	t := ""
	if s3 != nil && *s3 != "" {
		t = *s3
	}

	offset := int64(0)
	if i != nil {
		offset = int64(*i)
	}

	count := int64(25)
	if i2 != nil {
		count = int64(*i2)
	}

	return s.database.TicketList(ctx, t, q, strings, bools, offset, count)
}

func (s *Service) CreateTicket(ctx context.Context, form *model.TicketForm) (doc *model.TicketResponse, err error) {
	createdTickets, err := s.database.TicketBatchCreate(ctx, []*model.TicketForm{form})
	defer s.publishRequest(ctx, err, "CreateTicket", ticketIDs(createdTickets))
	if len(createdTickets) > 0 {
		return createdTickets[0], err
	}
	return nil, err
}

func (s *Service) CreateTicketBatch(ctx context.Context, ticketFormArray *model.TicketFormArray) error {
	if ticketFormArray == nil {
		return &api.HTTPError{Status: http.StatusUnprocessableEntity, Internal: errors.New("no tickets given")}
	}
	createdTickets, err := s.database.TicketBatchCreate(ctx, *ticketFormArray)
	defer s.publishRequest(ctx, err, "CreateTicket", ticketIDs(createdTickets))
	return err
}

func (s *Service) GetTicket(ctx context.Context, i int64) (*model.TicketWithTickets, error) {
	return s.database.TicketGet(ctx, i)
}

func (s *Service) UpdateTicket(ctx context.Context, i int64, ticket *model.Ticket) (doc *model.TicketWithTickets, err error) {
	defer s.publishRequest(ctx, err, "UpdateTicket", ticketWithTicketsID(doc))
	return s.database.TicketUpdate(ctx, i, ticket)
}

func (s *Service) DeleteTicket(ctx context.Context, i int64) (err error) {
	defer s.publishRequest(ctx, err, "DeleteTicket", ticketID(i))
	return s.database.TicketDelete(ctx, i)
}

func (s *Service) AddArtifact(ctx context.Context, i int64, artifact *model.Artifact) (doc *model.TicketWithTickets, err error) {
	defer s.publishRequest(ctx, err, "AddArtifact", ticketWithTicketsID(doc))
	return s.database.AddArtifact(ctx, i, artifact)
}

func (s *Service) GetArtifact(ctx context.Context, i int64, s2 string) (*model.Artifact, error) {
	return s.database.ArtifactGet(ctx, i, s2)
}

func (s *Service) SetArtifact(ctx context.Context, i int64, s2 string, artifact *model.Artifact) (doc *model.TicketWithTickets, err error) {
	defer s.publishRequest(ctx, err, "SetArtifact", ticketWithTicketsID(doc))
	return s.database.ArtifactUpdate(ctx, i, s2, artifact)
}

func (s *Service) RemoveArtifact(ctx context.Context, i int64, s2 string) (doc *model.TicketWithTickets, err error) {
	defer s.publishRequest(ctx, err, "RemoveArtifact", ticketWithTicketsID(doc))
	return s.database.RemoveArtifact(ctx, i, s2)
}

func (s *Service) EnrichArtifact(ctx context.Context, i int64, s2 string, form *model.EnrichmentForm) (doc *model.TicketWithTickets, err error) {
	defer s.publishRequest(ctx, err, "EnrichArtifact", ticketWithTicketsID(doc))
	return s.database.EnrichArtifact(ctx, i, s2, form)
}

func (s *Service) RunArtifact(ctx context.Context, id int64, name string, automation string) error {
	artifact, err := s.database.ArtifactGet(ctx, id, name)
	if err != nil {
		return err
	}

	defer s.publishRequest(ctx, err, "RunArtifact", ticketID(id))

	jobID := uuid.NewString()
	origin := &model.Origin{ArtifactOrigin: &model.ArtifactOrigin{TicketId: id, Artifact: name}}
	return s.bus.PublishJob(jobID, automation, name, &model.Context{Artifact: artifact}, origin)
}

func (s *Service) AddComment(ctx context.Context, i int64, form *model.CommentForm) (doc *model.TicketWithTickets, err error) {
	defer s.publishRequest(ctx, err, "AddComment", ticketWithTicketsID(doc))
	return s.database.AddComment(ctx, i, form)
}

func (s *Service) RemoveComment(ctx context.Context, i int64, i2 int) (doc *model.TicketWithTickets, err error) {
	defer s.publishRequest(ctx, err, "RemoveComment", ticketWithTicketsID(doc))
	return s.database.RemoveComment(ctx, i, int64(i2))
}

func (s *Service) AddTicketPlaybook(ctx context.Context, i int64, form *model.PlaybookTemplateForm) (doc *model.TicketWithTickets, err error) {
	defer s.publishRequest(ctx, err, "AddTicketPlaybook", ticketWithTicketsID(doc))
	return s.database.AddTicketPlaybook(ctx, i, form)
}

func (s *Service) RemoveTicketPlaybook(ctx context.Context, i int64, s2 string) (doc *model.TicketWithTickets, err error) {
	defer s.publishRequest(ctx, err, "RemoveTicketPlaybook", ticketWithTicketsID(doc))
	return s.database.RemoveTicketPlaybook(ctx, i, s2)
}

func (s *Service) SetTask(ctx context.Context, i int64, s3 string, s2 string, task *model.Task) (doc *model.TicketWithTickets, err error) {
	defer s.publishRequest(ctx, err, "SetTask", ticketWithTicketsID(doc))
	return s.database.TaskUpdate(ctx, i, s3, s2, task)
}

func (s *Service) CompleteTask(ctx context.Context, i int64, s3 string, s2 string, m map[string]interface{}) (doc *model.TicketWithTickets, err error) {
	defer s.publishRequest(ctx, err, "CompleteTask", ticketWithTicketsID(doc))
	return s.database.TaskComplete(ctx, i, s3, s2, m)
}

func (s *Service) RunTask(ctx context.Context, i int64, s3 string, s2 string) (err error) {
	defer s.publishRequest(ctx, err, "RunTask", ticketID(i))
	return s.database.TaskRun(ctx, i, s3, s2)
}

func (s *Service) SetReferences(ctx context.Context, i int64, references *model.ReferenceArray) (doc *model.TicketWithTickets, err error) {
	if references == nil {
		return nil, &api.HTTPError{Status: http.StatusUnprocessableEntity, Internal: errors.New("no references given")}
	}
	defer s.publishRequest(ctx, err, "SetReferences", ticketID(i))
	return s.database.SetReferences(ctx, i, *references)
}

func (s *Service) SetSchema(ctx context.Context, i int64, s2 string) (doc *model.TicketWithTickets, err error) {
	defer s.publishRequest(ctx, err, "SetSchema", ticketID(i))
	return s.database.SetTemplate(ctx, i, s2)
}

func (s *Service) LinkTicket(ctx context.Context, i int64, i2 int64) (*model.TicketWithTickets, error) {
	err := s.database.RelatedCreate(ctx, i, i2)
	if err != nil {
		return nil, err
	}
	defer s.publishRequest(ctx, err, "LinkTicket", ticketID(i))

	return s.GetTicket(ctx, i)
}

func (s *Service) UnlinkTicket(ctx context.Context, i int64, i2 int64) (*model.TicketWithTickets, error) {
	err := s.database.RelatedRemove(ctx, i, i2)
	if err != nil {
		return nil, err
	}
	defer s.publishRequest(ctx, err, "UnlinkTicket", ticketID(i))

	return s.GetTicket(ctx, i)
}
