package service

import (
	"context"

	"github.com/SecurityBrewery/catalyst/generated/restapi/api"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/tickets"
)

func (s *Service) AddArtifact(ctx context.Context, params *tickets.AddArtifactParams) *api.Response {
	i, err := s.database.AddArtifact(ctx, params.ID, params.Artifact)
	return s.response(ctx, "AddArtifact", ticketID(params.ID), i, err)
}

func (s *Service) RemoveArtifact(ctx context.Context, params *tickets.RemoveArtifactParams) *api.Response {
	i, err := s.database.RemoveArtifact(ctx, params.ID, params.Name)
	return s.response(ctx, "RemoveArtifact", ticketID(params.ID), i, err)
}

func (s *Service) SetSchema(ctx context.Context, params *tickets.SetSchemaParams) *api.Response {
	i, err := s.database.SetTemplate(ctx, params.ID, params.Schema)
	return s.response(ctx, "SetSchema", ticketID(params.ID), i, err)
}

func (s *Service) AddComment(ctx context.Context, params *tickets.AddCommentParams) *api.Response {
	i, err := s.database.AddComment(ctx, params.ID, params.Comment)
	return s.response(ctx, "AddComment", ticketID(params.ID), i, err)
}

func (s *Service) RemoveComment(ctx context.Context, params *tickets.RemoveCommentParams) *api.Response {
	i, err := s.database.RemoveComment(ctx, params.ID, params.CommentID)
	return s.response(ctx, "RemoveComment", ticketID(params.ID), i, err)
}

func (s *Service) LinkTicket(ctx context.Context, params *tickets.LinkTicketParams) *api.Response {
	err := s.database.RelatedCreate(ctx, params.ID, params.LinkedID)
	if err != nil {
		return s.response(ctx, "LinkTicket", ticketID(params.ID), nil, err)
	}

	i, err := s.database.TicketGet(ctx, params.ID)
	return s.response(ctx, "LinkTicket", ticketID(params.ID), i, err)
}

func (s *Service) UnlinkTicket(ctx context.Context, params *tickets.UnlinkTicketParams) *api.Response {
	err := s.database.RelatedRemove(ctx, params.ID, params.LinkedID)
	if err != nil {
		return s.response(ctx, "UnlinkTicket", ticketID(params.ID), nil, err)
	}

	i, err := s.database.TicketGet(ctx, params.ID)
	return s.response(ctx, "UnlinkTicket", ticketID(params.ID), i, err)
}

func (s Service) SetReferences(ctx context.Context, params *tickets.SetReferencesParams) *api.Response {
	i, err := s.database.SetReferences(ctx, params.ID, params.References)
	return s.response(ctx, "SetReferences", ticketID(params.ID), i, err)
}

func (s Service) LinkFiles(ctx context.Context, params *tickets.LinkFilesParams) *api.Response {
	i, err := s.database.LinkFiles(ctx, params.ID, params.Files)
	return s.response(ctx, "LinkFiles", ticketID(params.ID), i, err)
}

func (s Service) AddTicketPlaybook(ctx context.Context, params *tickets.AddTicketPlaybookParams) *api.Response {
	i, err := s.database.AddTicketPlaybook(ctx, params.ID, params.Playbook)
	return s.response(ctx, "AddTicketPlaybook", ticketID(params.ID), i, err)
}

func (s Service) RemoveTicketPlaybook(ctx context.Context, params *tickets.RemoveTicketPlaybookParams) *api.Response {
	i, err := s.database.RemoveTicketPlaybook(ctx, params.ID, params.PlaybookID)
	return s.response(ctx, "RemoveTicketPlaybook", ticketID(params.ID), i, err)
}

func (s Service) CompleteTask(ctx context.Context, params *tickets.CompleteTaskParams) *api.Response {
	i, err := s.database.TaskComplete(ctx, params.ID, params.PlaybookID, params.TaskID, params.Data)
	return s.response(ctx, "CompleteTask", ticketID(params.ID), i, err)
}

func (s Service) SetTask(ctx context.Context, params *tickets.SetTaskParams) *api.Response {
	i, err := s.database.TaskUpdate(ctx, params.ID, params.PlaybookID, params.TaskID, params.Task)
	return s.response(ctx, "SetTask", ticketID(params.ID), i, err)
}

func (s *Service) RunTask(ctx context.Context, params *tickets.RunTaskParams) *api.Response {
	err := s.database.TaskRun(ctx, params.ID, params.PlaybookID, params.TaskID)
	return s.response(ctx, "RunTask", ticketID(params.ID), nil, err)
}
