package service

import (
	"context"

	"github.com/SecurityBrewery/catalyst/generated/restapi/api"
	"github.com/SecurityBrewery/catalyst/generated/restapi/operations/tickets"
)

func (s *Service) AddArtifact(ctx context.Context, params *tickets.AddArtifactParams) *api.Response {
	return response(s.database.AddArtifact(ctx, params.ID, params.Artifact))
}

func (s *Service) RemoveArtifact(ctx context.Context, params *tickets.RemoveArtifactParams) *api.Response {
	return response(s.database.RemoveArtifact(ctx, params.ID, params.Name))
}

func (s *Service) SetSchema(ctx context.Context, params *tickets.SetSchemaParams) *api.Response {
	return response(s.database.SetTemplate(ctx, params.ID, params.Schema))
}

func (s *Service) AddComment(ctx context.Context, params *tickets.AddCommentParams) *api.Response {
	return response(s.database.AddComment(ctx, params.ID, params.Comment))
}

func (s *Service) RemoveComment(ctx context.Context, params *tickets.RemoveCommentParams) *api.Response {
	return response(s.database.RemoveComment(ctx, params.ID, params.CommentID))
}

func (s *Service) LinkTicket(ctx context.Context, params *tickets.LinkTicketParams) *api.Response {
	err := s.database.RelatedCreate(ctx, params.ID, params.LinkedID)
	if err != nil {
		return response(nil, err)
	}

	return s.GetTicket(ctx, &tickets.GetTicketParams{ID: params.ID})
}

func (s *Service) UnlinkTicket(ctx context.Context, params *tickets.UnlinkTicketParams) *api.Response {
	err := s.database.RelatedRemove(ctx, params.ID, params.LinkedID)
	if err != nil {
		return response(nil, err)
	}

	return s.GetTicket(ctx, &tickets.GetTicketParams{ID: params.ID})
}

func (s Service) SetReferences(ctx context.Context, params *tickets.SetReferencesParams) *api.Response {
	return response(s.database.SetReferences(ctx, params.ID, params.References))
}

func (s Service) LinkFiles(ctx context.Context, params *tickets.LinkFilesParams) *api.Response {
	return response(s.database.LinkFiles(ctx, params.ID, params.Files))
}

func (s Service) AddTicketPlaybook(ctx context.Context, params *tickets.AddTicketPlaybookParams) *api.Response {
	return response(s.database.AddTicketPlaybook(ctx, params.ID, params.Playbook))
}

func (s Service) RemoveTicketPlaybook(ctx context.Context, params *tickets.RemoveTicketPlaybookParams) *api.Response {
	return response(s.database.RemoveTicketPlaybook(ctx, params.ID, params.PlaybookID))
}

func (s Service) CompleteTask(ctx context.Context, params *tickets.CompleteTaskParams) *api.Response {
	return response(s.database.TaskComplete(ctx, params.ID, params.PlaybookID, params.TaskID, params.Data))
}

func (s Service) SetTask(ctx context.Context, params *tickets.SetTaskParams) *api.Response {
	return response(s.database.TaskUpdate(ctx, params.ID, params.PlaybookID, params.TaskID, params.Task))
}

func (s *Service) RunTask(ctx context.Context, params *tickets.RunTaskParams) *api.Response {
	return response(nil, s.database.TaskRun(ctx, params.ID, params.PlaybookID, params.TaskID))
}
