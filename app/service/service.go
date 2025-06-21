package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/SecurityBrewery/catalyst/app/auth/password"
	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/hook"
	"github.com/SecurityBrewery/catalyst/app/openapi"
	"github.com/SecurityBrewery/catalyst/app/permission"
	"github.com/SecurityBrewery/catalyst/app/reaction/schedule"
)

const (
	defaultLimit  = 100
	defaultOffset = 0
)

var _ openapi.StrictServerInterface = (*Service)(nil)

type Service struct {
	queries   *sqlc.Queries
	hooks     *hook.Hooks
	scheduler *schedule.Scheduler
}

func New(queries *sqlc.Queries, hooks *hook.Hooks, scheduler *schedule.Scheduler) *Service {
	return &Service{
		queries:   queries,
		hooks:     hooks,
		scheduler: scheduler,
	}
}

func (s *Service) ListComments(ctx context.Context, request openapi.ListCommentsRequestObject) (openapi.ListCommentsResponseObject, error) {
	comments, err := s.queries.ListComments(ctx, sqlc.ListCommentsParams{
		Ticket: toString(request.Params.Ticket, ""),
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
	})
	if err != nil {
		return nil, err
	}

	response := make([]openapi.ExtendedComment, 0, len(comments))
	for _, comment := range comments {
		response = append(response, openapi.ExtendedComment{
			Author:     comment.Author,
			Created:    comment.Created,
			Id:         comment.ID,
			Message:    comment.Message,
			Ticket:     comment.Ticket,
			Updated:    comment.Updated,
			AuthorName: comment.AuthorName.String,
		})
	}

	totalCount := 0
	if len(comments) > 0 {
		totalCount = int(comments[0].TotalCount)
	}

	s.hooks.OnRecordsListRequest.Publish(ctx, permission.CommentsTable.ID, response)

	return openapi.ListComments200JSONResponse{
		Body: response,
		Headers: openapi.ListComments200ResponseHeaders{
			XTotalCount: totalCount,
		},
	}, nil
}

func (s *Service) CreateComment(ctx context.Context, request openapi.CreateCommentRequestObject) (openapi.CreateCommentResponseObject, error) {
	s.hooks.OnRecordBeforeCreateRequest.Publish(ctx, permission.CommentsTable.ID, request.Body)

	comment, err := s.queries.CreateComment(ctx, sqlc.CreateCommentParams{
		ID:      database.GenerateID("c"),
		Author:  request.Body.Author,
		Message: request.Body.Message,
		Ticket:  request.Body.Ticket,
	})
	if err != nil {
		return nil, err
	}

	response := openapi.Comment{
		Author:  comment.Author,
		Created: comment.Created,
		Id:      comment.ID,
		Message: comment.Message,
		Ticket:  comment.Ticket,
		Updated: comment.Updated,
	}

	s.hooks.OnRecordAfterCreateRequest.Publish(ctx, permission.CommentsTable.ID, response)

	return openapi.CreateComment200JSONResponse(response), nil
}

func (s *Service) DeleteComment(ctx context.Context, request openapi.DeleteCommentRequestObject) (openapi.DeleteCommentResponseObject, error) {
	s.hooks.OnRecordBeforeDeleteRequest.Publish(ctx, permission.CommentsTable.ID, request.Id)

	err := s.queries.DeleteComment(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterDeleteRequest.Publish(ctx, permission.CommentsTable.ID, request.Id)

	return openapi.DeleteComment204Response{}, nil
}

func (s *Service) GetComment(ctx context.Context, request openapi.GetCommentRequestObject) (openapi.GetCommentResponseObject, error) {
	comment, err := s.queries.GetComment(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	response := openapi.ExtendedComment{
		Author:     comment.Author,
		AuthorName: comment.AuthorName.String,
		Created:    comment.Created,
		Id:         comment.ID,
		Message:    comment.Message,
		Ticket:     comment.Ticket,
		Updated:    comment.Updated,
	}

	s.hooks.OnRecordViewRequest.Publish(ctx, permission.CommentsTable.ID, response)

	return openapi.GetComment200JSONResponse(response), nil
}

func (s *Service) UpdateComment(ctx context.Context, request openapi.UpdateCommentRequestObject) (openapi.UpdateCommentResponseObject, error) {
	s.hooks.OnRecordBeforeUpdateRequest.Publish(ctx, permission.CommentsTable.ID, request.Body)

	comment, err := s.queries.UpdateComment(ctx, sqlc.UpdateCommentParams{
		Message: toNullString(request.Body.Message),
		ID:      request.Id,
	})
	if err != nil {
		return nil, err
	}

	response := openapi.Comment{
		Author:  comment.Author,
		Created: comment.Created,
		Id:      comment.ID,
		Message: comment.Message,
		Ticket:  comment.Ticket,
		Updated: comment.Updated,
	}

	s.hooks.OnRecordAfterUpdateRequest.Publish(ctx, permission.CommentsTable.ID, response)

	return openapi.UpdateComment200JSONResponse(response), nil
}

func (s *Service) GetDashboardCounts(ctx context.Context, _ openapi.GetDashboardCountsRequestObject) (openapi.GetDashboardCountsResponseObject, error) {
	counts, err := s.queries.GetDashboardCounts(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]openapi.DashboardCounts, 0, len(counts))
	for _, count := range counts {
		response = append(response, openapi.DashboardCounts{
			Id:    count.ID,
			Count: int(count.Count),
		})
	}

	s.hooks.OnRecordsListRequest.Publish(ctx, permission.DashboardCountsTable.ID, response)

	return openapi.GetDashboardCounts200JSONResponse(response), nil
}

func (s *Service) ListFiles(ctx context.Context, request openapi.ListFilesRequestObject) (openapi.ListFilesResponseObject, error) {
	files, err := s.queries.ListFiles(ctx, sqlc.ListFilesParams{
		Ticket: toString(request.Params.Ticket, ""),
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
	})
	if err != nil {
		return nil, err
	}

	response := make([]openapi.File, 0, len(files))
	for _, file := range files {
		response = append(response, openapi.File{
			Created: file.Created,
			Id:      file.ID,
			Name:    file.Name,
			Size:    file.Size,
			Ticket:  file.Ticket,
			Updated: file.Updated,
		})
	}

	s.hooks.OnRecordsListRequest.Publish(ctx, permission.FilesTable.ID, response)

	totalCount := 0
	if len(files) > 0 {
		totalCount = int(files[0].TotalCount)
	}

	return openapi.ListFiles200JSONResponse{
		Body: response,
		Headers: openapi.ListFiles200ResponseHeaders{
			XTotalCount: totalCount,
		},
	}, nil
}

func (s *Service) CreateFile(ctx context.Context, request openapi.CreateFileRequestObject) (openapi.CreateFileResponseObject, error) {
	s.hooks.OnRecordBeforeCreateRequest.Publish(ctx, permission.FilesTable.ID, request.Body)

	file, err := s.queries.CreateFile(ctx, sqlc.CreateFileParams{
		ID:     database.GenerateID("b"),
		Name:   request.Body.Name,
		Blob:   request.Body.Blob,
		Size:   request.Body.Size,
		Ticket: request.Body.Ticket,
	})
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterCreateRequest.Publish(ctx, permission.FilesTable.ID, file)

	return openapi.CreateFile200JSONResponse(openapi.File{
		Created: file.Created,
		Id:      file.ID,
		Name:    file.Name,
		Size:    file.Size,
		Ticket:  file.Ticket,
		Updated: file.Updated,
	}), nil
}

func (s *Service) DeleteFile(ctx context.Context, request openapi.DeleteFileRequestObject) (openapi.DeleteFileResponseObject, error) {
	s.hooks.OnRecordBeforeDeleteRequest.Publish(ctx, permission.FilesTable.ID, request.Id)

	err := s.queries.DeleteFile(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterDeleteRequest.Publish(ctx, permission.FilesTable.ID, request.Id)

	return openapi.DeleteFile204Response{}, nil
}

func (s *Service) GetFile(ctx context.Context, request openapi.GetFileRequestObject) (openapi.GetFileResponseObject, error) {
	file, err := s.queries.GetFile(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	response := openapi.File{
		Created: file.Created,
		Id:      file.ID,
		Name:    file.Name,
		Size:    file.Size,
		Ticket:  file.Ticket,
		Updated: file.Updated,
	}

	s.hooks.OnRecordViewRequest.Publish(ctx, permission.FilesTable.ID, response)

	return openapi.GetFile200JSONResponse(response), nil
}

func (s *Service) UpdateFile(ctx context.Context, request openapi.UpdateFileRequestObject) (openapi.UpdateFileResponseObject, error) {
	s.hooks.OnRecordBeforeUpdateRequest.Publish(ctx, permission.FilesTable.ID, request.Body)

	file, err := s.queries.UpdateFile(ctx, sqlc.UpdateFileParams{
		ID:   request.Id,
		Name: toNullString(request.Body.Name),
	})
	if err != nil {
		return nil, err
	}

	response := openapi.File{
		Created: file.Created,
		Id:      file.ID,
		Name:    file.Name,
		Size:    file.Size,
		Ticket:  file.Ticket,
		Updated: file.Updated,
	}

	s.hooks.OnRecordAfterUpdateRequest.Publish(ctx, permission.FilesTable.ID, response)

	return openapi.UpdateFile200JSONResponse(response), nil
}

func (s *Service) ListLinks(ctx context.Context, request openapi.ListLinksRequestObject) (openapi.ListLinksResponseObject, error) {
	links, err := s.queries.ListLinks(ctx, sqlc.ListLinksParams{
		Ticket: toString(request.Params.Ticket, ""),
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
	})
	if err != nil {
		return nil, err
	}

	response := make([]openapi.Link, 0, len(links))
	for _, link := range links {
		response = append(response, openapi.Link{
			Id:      link.ID,
			Name:    link.Name,
			Created: link.Created,
			Updated: link.Updated,
			Url:     link.Url,
			Ticket:  link.Ticket,
		})
	}

	s.hooks.OnRecordsListRequest.Publish(ctx, permission.LinksTable.ID, response)

	totalCount := 0
	if len(links) > 0 {
		totalCount = int(links[0].TotalCount)
	}

	return openapi.ListLinks200JSONResponse{
		Body: response,
		Headers: openapi.ListLinks200ResponseHeaders{
			XTotalCount: totalCount,
		},
	}, nil
}

func (s *Service) DownloadFile(ctx context.Context, request openapi.DownloadFileRequestObject) (openapi.DownloadFileResponseObject, error) {
	file, err := s.queries.GetFile(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	header, blob, ok := strings.Cut(file.Blob, ";base64,")
	if !ok {
		return nil, fmt.Errorf("invalid file blob format for file ID %s", request.Id)
	}

	data, err := base64.StdEncoding.DecodeString(blob)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file blob: %w", err)
	}

	return openapi.DownloadFile200ApplicationoctetStreamResponse{
		Body:          bytes.NewReader(data),
		ContentLength: int64(len(data)),
		Headers: openapi.DownloadFile200ResponseHeaders{
			ContentDisposition: "attachment; filename=\"" + file.Name + "\"",
			ContentType:        strings.TrimPrefix(header, "data:"),
		},
	}, nil
}

func (s *Service) CreateLink(ctx context.Context, request openapi.CreateLinkRequestObject) (openapi.CreateLinkResponseObject, error) {
	s.hooks.OnRecordBeforeCreateRequest.Publish(ctx, permission.LinksTable.ID, request.Body)

	link, err := s.queries.CreateLink(ctx, sqlc.CreateLinkParams{
		ID:     database.GenerateID("l"),
		Name:   request.Body.Name,
		Url:    request.Body.Url,
		Ticket: request.Body.Ticket,
	})
	if err != nil {
		return nil, err
	}

	response := openapi.Link{
		Created: link.Created,
		Id:      link.ID,
		Name:    link.Name,
		Ticket:  link.Ticket,
		Updated: link.Updated,
		Url:     link.Url,
	}

	s.hooks.OnRecordAfterCreateRequest.Publish(ctx, permission.LinksTable.ID, response)

	return openapi.CreateLink200JSONResponse(response), nil
}

func (s *Service) DeleteLink(ctx context.Context, request openapi.DeleteLinkRequestObject) (openapi.DeleteLinkResponseObject, error) {
	s.hooks.OnRecordBeforeDeleteRequest.Publish(ctx, permission.LinksTable.ID, request.Id)

	err := s.queries.DeleteLink(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterDeleteRequest.Publish(ctx, permission.LinksTable.ID, request.Id)

	return openapi.DeleteLink204Response{}, nil
}

func (s *Service) GetLink(ctx context.Context, request openapi.GetLinkRequestObject) (openapi.GetLinkResponseObject, error) {
	link, err := s.queries.GetLink(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	response := openapi.Link{
		Id:      link.ID,
		Name:    link.Name,
		Created: link.Created,
		Updated: link.Updated,
		Url:     link.Url,
		Ticket:  link.Ticket,
	}

	s.hooks.OnRecordViewRequest.Publish(ctx, permission.LinksTable.ID, response)

	return openapi.GetLink200JSONResponse(response), nil
}

func (s *Service) UpdateLink(ctx context.Context, request openapi.UpdateLinkRequestObject) (openapi.UpdateLinkResponseObject, error) {
	s.hooks.OnRecordBeforeUpdateRequest.Publish(ctx, permission.LinksTable.ID, request.Body)

	link, err := s.queries.UpdateLink(ctx, sqlc.UpdateLinkParams{
		ID:   request.Id,
		Name: toNullString(request.Body.Name),
		Url:  toNullString(request.Body.Url),
	})
	if err != nil {
		return nil, err
	}

	response := openapi.Link{
		Created: link.Created,
		Id:      link.ID,
		Name:    link.Name,
		Ticket:  link.Ticket,
		Updated: link.Updated,
		Url:     link.Url,
	}

	s.hooks.OnRecordAfterUpdateRequest.Publish(ctx, permission.LinksTable.ID, response)

	return openapi.UpdateLink200JSONResponse(response), nil
}

func (s *Service) ListReactions(ctx context.Context, request openapi.ListReactionsRequestObject) (openapi.ListReactionsResponseObject, error) {
	reactions, err := s.queries.ListReactions(ctx, sqlc.ListReactionsParams{
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
	})
	if err != nil {
		return nil, err
	}

	response := make([]openapi.Reaction, 0, len(reactions))
	for _, reaction := range reactions {
		response = append(response, openapi.Reaction{
			Action:      reaction.Action,
			Actiondata:  unmarshal(reaction.Actiondata),
			Created:     reaction.Created,
			Id:          reaction.ID,
			Name:        reaction.Name,
			Trigger:     reaction.Trigger,
			Triggerdata: unmarshal(reaction.Triggerdata),
			Updated:     reaction.Updated,
		})
	}

	s.hooks.OnRecordsListRequest.Publish(ctx, permission.ReactionsTable.ID, response)

	totalCount := 0
	if len(reactions) > 0 {
		totalCount = int(reactions[0].TotalCount)
	}

	return openapi.ListReactions200JSONResponse{
		Body: response,
		Headers: openapi.ListReactions200ResponseHeaders{
			XTotalCount: totalCount,
		},
	}, nil
}

func (s *Service) CreateReaction(ctx context.Context, request openapi.CreateReactionRequestObject) (openapi.CreateReactionResponseObject, error) {
	s.hooks.OnRecordBeforeCreateRequest.Publish(ctx, permission.ReactionsTable.ID, request.Body)

	reaction, err := s.queries.CreateReaction(ctx, sqlc.CreateReactionParams{
		ID:          database.GenerateID("r"),
		Name:        request.Body.Name,
		Action:      request.Body.Action,
		Trigger:     request.Body.Trigger,
		Actiondata:  marshal(request.Body.Actiondata),
		Triggerdata: marshal(request.Body.Triggerdata),
	})
	if err != nil {
		return nil, err
	}

	if err := s.scheduler.AddReaction(&reaction); err != nil {
		slog.ErrorContext(ctx, "Failed to add reaction to scheduler", "error", err)
	}

	response := openapi.Reaction{
		Action:      reaction.Action,
		Actiondata:  unmarshal(reaction.Actiondata),
		Created:     reaction.Created,
		Id:          reaction.ID,
		Name:        reaction.Name,
		Trigger:     reaction.Trigger,
		Triggerdata: unmarshal(reaction.Triggerdata),
		Updated:     reaction.Updated,
	}

	s.hooks.OnRecordAfterCreateRequest.Publish(ctx, permission.ReactionsTable.ID, response)

	return openapi.CreateReaction200JSONResponse(response), nil
}

func (s *Service) DeleteReaction(ctx context.Context, request openapi.DeleteReactionRequestObject) (openapi.DeleteReactionResponseObject, error) {
	s.hooks.OnRecordBeforeDeleteRequest.Publish(ctx, permission.ReactionsTable.ID, request.Id)

	err := s.queries.DeleteReaction(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.scheduler.RemoveReaction(request.Id)

	s.hooks.OnRecordAfterDeleteRequest.Publish(ctx, permission.ReactionsTable.ID, request.Id)

	return openapi.DeleteReaction204Response{}, nil
}

func (s *Service) GetReaction(ctx context.Context, request openapi.GetReactionRequestObject) (openapi.GetReactionResponseObject, error) {
	reaction, err := s.queries.GetReaction(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	response := openapi.Reaction{
		Action:      reaction.Action,
		Actiondata:  unmarshal(reaction.Actiondata),
		Created:     reaction.Created,
		Id:          reaction.ID,
		Name:        reaction.Name,
		Trigger:     reaction.Trigger,
		Triggerdata: unmarshal(reaction.Triggerdata),
		Updated:     reaction.Updated,
	}

	s.hooks.OnRecordViewRequest.Publish(ctx, permission.ReactionsTable.ID, response)

	return openapi.GetReaction200JSONResponse(response), nil
}

func (s *Service) UpdateReaction(ctx context.Context, request openapi.UpdateReactionRequestObject) (openapi.UpdateReactionResponseObject, error) {
	s.hooks.OnRecordBeforeUpdateRequest.Publish(ctx, permission.ReactionsTable.ID, request.Body)

	reaction, err := s.queries.UpdateReaction(ctx, sqlc.UpdateReactionParams{
		ID:          request.Id,
		Name:        toNullString(request.Body.Name),
		Action:      toNullString(request.Body.Action),
		Trigger:     toNullString(request.Body.Trigger),
		Actiondata:  marshalPointer(request.Body.Actiondata),
		Triggerdata: marshalPointer(request.Body.Triggerdata),
	})
	if err != nil {
		return nil, err
	}

	s.scheduler.RemoveReaction(request.Id)

	if err := s.scheduler.AddReaction(&reaction); err != nil {
		slog.ErrorContext(ctx, "Failed to add reaction to scheduler", "error", err)
	}

	response := openapi.Reaction{
		Action:      reaction.Action,
		Actiondata:  unmarshal(reaction.Actiondata),
		Created:     reaction.Created,
		Id:          reaction.ID,
		Name:        reaction.Name,
		Trigger:     reaction.Trigger,
		Triggerdata: unmarshal(reaction.Triggerdata),
		Updated:     reaction.Updated,
	}

	s.hooks.OnRecordAfterUpdateRequest.Publish(ctx, permission.ReactionsTable.ID, response)

	return openapi.UpdateReaction200JSONResponse(response), nil
}

func (s *Service) GetSidebar(ctx context.Context, _ openapi.GetSidebarRequestObject) (openapi.GetSidebarResponseObject, error) {
	sidebar, err := s.queries.GetSidebar(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]openapi.Sidebar, 0, len(sidebar))
	for _, s := range sidebar {
		response = append(response, openapi.Sidebar{
			Id:       s.ID,
			Singular: s.Singular,
			Plural:   s.Plural,
			Icon:     s.Icon,
			Count:    int(s.Count),
		})
	}

	s.hooks.OnRecordsListRequest.Publish(ctx, permission.SidebarTable.ID, response)

	return openapi.GetSidebar200JSONResponse(response), nil
}

func (s *Service) ListTasks(ctx context.Context, request openapi.ListTasksRequestObject) (openapi.ListTasksResponseObject, error) {
	tasks, err := s.queries.ListTasks(ctx, sqlc.ListTasksParams{
		Ticket: toString(request.Params.Ticket, ""),
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
	})
	if err != nil {
		return nil, err
	}

	response := make([]openapi.ExtendedTask, 0, len(tasks))
	for _, task := range tasks {
		response = append(response, openapi.ExtendedTask{
			Id:         task.ID,
			Name:       task.Name,
			Created:    task.Created,
			Updated:    task.Updated,
			Open:       task.Open,
			Owner:      task.Owner,
			Ticket:     task.Ticket,
			OwnerName:  task.OwnerName.String,
			TicketName: task.TicketName.String,
			TicketType: task.TicketType.String,
		})
	}

	s.hooks.OnRecordsListRequest.Publish(ctx, permission.TasksTable.ID, response)

	totalCount := 0
	if len(tasks) > 0 {
		totalCount = int(tasks[0].TotalCount)
	}

	return openapi.ListTasks200JSONResponse{
		Body: response,
		Headers: openapi.ListTasks200ResponseHeaders{
			XTotalCount: totalCount,
		},
	}, nil
}

func (s *Service) CreateTask(ctx context.Context, request openapi.CreateTaskRequestObject) (openapi.CreateTaskResponseObject, error) {
	s.hooks.OnRecordBeforeCreateRequest.Publish(ctx, permission.TasksTable.ID, request.Body)

	task, err := s.queries.CreateTask(ctx, sqlc.CreateTaskParams{
		ID:     database.GenerateID("k"),
		Name:   request.Body.Name,
		Open:   request.Body.Open,
		Owner:  request.Body.Owner,
		Ticket: request.Body.Ticket,
	})
	if err != nil {
		return nil, err
	}

	response := openapi.Task{
		Created: task.Created,
		Id:      task.ID,
		Name:    task.Name,
		Open:    task.Open,
		Owner:   task.Owner,
		Ticket:  task.Ticket,
		Updated: task.Updated,
	}

	s.hooks.OnRecordAfterCreateRequest.Publish(ctx, permission.TasksTable.ID, response)

	return openapi.CreateTask200JSONResponse(response), nil
}

func (s *Service) DeleteTask(ctx context.Context, request openapi.DeleteTaskRequestObject) (openapi.DeleteTaskResponseObject, error) {
	s.hooks.OnRecordBeforeDeleteRequest.Publish(ctx, permission.TasksTable.ID, request.Id)

	err := s.queries.DeleteTask(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterDeleteRequest.Publish(ctx, permission.TasksTable.ID, request.Id)

	return openapi.DeleteTask204Response{}, nil
}

func (s *Service) GetTask(ctx context.Context, request openapi.GetTaskRequestObject) (openapi.GetTaskResponseObject, error) {
	task, err := s.queries.GetTask(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	response := openapi.ExtendedTask{
		Id:         task.ID,
		Name:       task.Name,
		Created:    task.Created,
		Updated:    task.Updated,
		Open:       task.Open,
		Owner:      task.Owner,
		Ticket:     task.Ticket,
		OwnerName:  task.OwnerName.String,
		TicketName: task.TicketName.String,
		TicketType: task.TicketType.String,
	}

	s.hooks.OnRecordViewRequest.Publish(ctx, permission.TasksTable.ID, response)

	return openapi.GetTask200JSONResponse(response), nil
}

func (s *Service) UpdateTask(ctx context.Context, request openapi.UpdateTaskRequestObject) (openapi.UpdateTaskResponseObject, error) {
	s.hooks.OnRecordBeforeUpdateRequest.Publish(ctx, permission.TasksTable.ID, request.Body)

	task, err := s.queries.UpdateTask(ctx, sqlc.UpdateTaskParams{
		ID:    request.Id,
		Name:  toNullString(request.Body.Name),
		Open:  toNullBool(request.Body.Open),
		Owner: toNullString(request.Body.Owner),
	})
	if err != nil {
		return nil, err
	}

	response := openapi.Task{
		Created: task.Created,
		Id:      task.ID,
		Name:    task.Name,
		Open:    task.Open,
		Owner:   task.Owner,
		Ticket:  task.Ticket,
		Updated: task.Updated,
	}

	s.hooks.OnRecordAfterUpdateRequest.Publish(ctx, permission.TasksTable.ID, response)

	return openapi.UpdateTask200JSONResponse(response), nil
}

func (s *Service) SearchTickets(ctx context.Context, request openapi.SearchTicketsRequestObject) (openapi.SearchTicketsResponseObject, error) {
	tickets, err := s.queries.SearchTickets(ctx, sqlc.SearchTicketsParams{
		Query:  toNullString(request.Params.Query),
		Type:   toNullString(request.Params.Type),
		Open:   toNullBool(request.Params.Open),
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
	})
	if err != nil {
		return nil, err
	}

	response := make([]openapi.TicketSearch, 0, len(tickets))

	for _, ticket := range tickets {
		response = append(response, openapi.TicketSearch{
			Created:     ticket.Created,
			Description: ticket.Description,
			Id:          ticket.ID,
			Name:        ticket.Name,
			Open:        ticket.Open,
			OwnerName:   ticket.OwnerName.String,
			State:       unmarshal(ticket.State),
			Type:        ticket.Type,
		})
	}

	s.hooks.OnRecordsListRequest.Publish(ctx, permission.TicketsTable.ID, response)

	totalCount := 0
	if len(tickets) > 0 {
		totalCount = int(tickets[0].TotalCount)
	}

	return openapi.SearchTickets200JSONResponse{
		Body: response,
		Headers: openapi.SearchTickets200ResponseHeaders{
			XTotalCount: totalCount,
		},
	}, nil
}

func (s *Service) ListTickets(ctx context.Context, request openapi.ListTicketsRequestObject) (openapi.ListTicketsResponseObject, error) {
	tickets, err := s.queries.ListTickets(ctx, sqlc.ListTicketsParams{
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
	})
	if err != nil {
		return nil, err
	}

	response := make([]openapi.ExtendedTicket, 0, len(tickets))
	for _, ticket := range tickets {
		response = append(response, openapi.ExtendedTicket{
			Created:      ticket.Created,
			Description:  ticket.Description,
			Id:           ticket.ID,
			Name:         ticket.Name,
			Open:         ticket.Open,
			Owner:        ticket.Owner,
			OwnerName:    ticket.OwnerName.String,
			Resolution:   ticket.Resolution,
			Type:         ticket.Type,
			Schema:       unmarshal(ticket.Schema),
			State:        unmarshal(ticket.State),
			TypePlural:   ticket.TypePlural.String,
			TypeSingular: ticket.TypeSingular.String,
			Updated:      ticket.Updated,
		})
	}

	s.hooks.OnRecordsListRequest.Publish(ctx, "tickets", response)

	totalCount := 0
	if len(tickets) > 0 {
		totalCount = int(tickets[0].TotalCount)
	}

	return openapi.ListTickets200JSONResponse{
		Body: response,
		Headers: openapi.ListTickets200ResponseHeaders{
			XTotalCount: totalCount,
		},
	}, nil
}

func (s *Service) CreateTicket(ctx context.Context, request openapi.CreateTicketRequestObject) (openapi.CreateTicketResponseObject, error) {
	s.hooks.OnRecordBeforeCreateRequest.Publish(ctx, permission.TicketsTable.ID, request.Body)

	ticket, err := s.queries.CreateTicket(ctx, sqlc.CreateTicketParams{
		ID:          database.GenerateID(request.Body.Type),
		Name:        request.Body.Name,
		Description: request.Body.Description,
		Owner:       request.Body.Owner,
		Open:        request.Body.Open,
		Resolution:  request.Body.Resolution,
		Type:        request.Body.Type,
		State:       marshal(request.Body.State),
	})
	if err != nil {
		return nil, err
	}

	response := openapi.Ticket{
		Created:     ticket.Created,
		Description: ticket.Description,
		Id:          ticket.ID,
		Name:        ticket.Name,
		Open:        ticket.Open,
		Owner:       ticket.Owner,
		Resolution:  ticket.Resolution,
		Schema:      unmarshal(ticket.Schema),
		State:       unmarshal(ticket.State),
		Type:        ticket.Type,
		Updated:     ticket.Updated,
	}

	s.hooks.OnRecordAfterCreateRequest.Publish(ctx, permission.TicketsTable.ID, response)

	return openapi.CreateTicket200JSONResponse(response), nil
}

func (s *Service) DeleteTicket(ctx context.Context, request openapi.DeleteTicketRequestObject) (openapi.DeleteTicketResponseObject, error) {
	s.hooks.OnRecordBeforeDeleteRequest.Publish(ctx, permission.TicketsTable.ID, request.Id)

	err := s.queries.DeleteTicket(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterDeleteRequest.Publish(ctx, permission.TicketsTable.ID, request.Id)

	return openapi.DeleteTicket204Response{}, nil
}

func (s *Service) GetTicket(ctx context.Context, request openapi.GetTicketRequestObject) (openapi.GetTicketResponseObject, error) {
	ticket, err := s.queries.Ticket(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	response := openapi.ExtendedTicket{
		Created:      ticket.Created,
		Description:  ticket.Description,
		Id:           ticket.ID,
		Name:         ticket.Name,
		Open:         ticket.Open,
		Owner:        ticket.Owner,
		OwnerName:    ticket.OwnerName.String,
		Resolution:   ticket.Resolution,
		Schema:       unmarshal(ticket.Schema),
		State:        unmarshal(ticket.State),
		Type:         ticket.Type,
		TypePlural:   ticket.TypePlural.String,
		TypeSingular: ticket.TypeSingular.String,
		Updated:      ticket.Updated,
	}

	s.hooks.OnRecordViewRequest.Publish(ctx, permission.TicketsTable.ID, response)

	return openapi.GetTicket200JSONResponse(response), nil
}

func (s *Service) UpdateTicket(ctx context.Context, request openapi.UpdateTicketRequestObject) (openapi.UpdateTicketResponseObject, error) {
	s.hooks.OnRecordBeforeUpdateRequest.Publish(ctx, permission.TicketsTable.ID, request.Body)

	ticket, err := s.queries.UpdateTicket(ctx, sqlc.UpdateTicketParams{
		Name:        toNullString(request.Body.Name),
		Description: toNullString(request.Body.Description),
		Open:        toNullBool(request.Body.Open),
		Owner:       toNullString(request.Body.Owner),
		Resolution:  toNullString(request.Body.Resolution),
		Schema:      marshalPointer(request.Body.Schema),
		State:       marshalPointer(request.Body.State),
		Type:        toNullString(request.Body.Type),
		ID:          request.Id,
	})
	if err != nil {
		return nil, err
	}

	response := openapi.Ticket{
		Created:     ticket.Created,
		Description: ticket.Description,
		Id:          ticket.ID,
		Name:        ticket.Name,
		Open:        ticket.Open,
		Owner:       ticket.Owner,
		Resolution:  ticket.Resolution,
		Schema:      unmarshal(ticket.Schema),
		State:       unmarshal(ticket.State),
		Type:        ticket.Type,
		Updated:     ticket.Updated,
	}

	s.hooks.OnRecordAfterUpdateRequest.Publish(ctx, permission.TicketsTable.ID, response)

	return openapi.UpdateTicket200JSONResponse(response), nil
}

func (s *Service) ListTimeline(ctx context.Context, request openapi.ListTimelineRequestObject) (openapi.ListTimelineResponseObject, error) {
	timeline, err := s.queries.ListTimeline(ctx, sqlc.ListTimelineParams{
		Ticket: toString(request.Params.Ticket, ""),
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
	})
	if err != nil {
		return nil, err
	}

	response := make([]openapi.TimelineEntry, 0, len(timeline))
	for _, timeline := range timeline {
		response = append(response, openapi.TimelineEntry{
			Id:      timeline.ID,
			Message: timeline.Message,
			Created: timeline.Created,
			Updated: timeline.Updated,
			Time:    timeline.Time,
			Ticket:  timeline.Ticket,
		})
	}

	s.hooks.OnRecordsListRequest.Publish(ctx, permission.TimelinesTable.ID, response)

	totalCount := 0
	if len(timeline) > 0 {
		totalCount = int(timeline[0].TotalCount)
	}

	return openapi.ListTimeline200JSONResponse{
		Body: response,
		Headers: openapi.ListTimeline200ResponseHeaders{
			XTotalCount: totalCount,
		},
	}, nil
}

func (s *Service) CreateTimeline(ctx context.Context, request openapi.CreateTimelineRequestObject) (openapi.CreateTimelineResponseObject, error) {
	s.hooks.OnRecordBeforeCreateRequest.Publish(ctx, permission.TimelinesTable.ID, request.Body)

	timeline, err := s.queries.CreateTimeline(ctx, sqlc.CreateTimelineParams{
		ID:      database.GenerateID("h"),
		Message: request.Body.Message,
		Time:    request.Body.Time,
		Ticket:  request.Body.Ticket,
	})
	if err != nil {
		return nil, err
	}

	response := openapi.TimelineEntry{
		Created: timeline.Created,
		Id:      timeline.ID,
		Message: timeline.Message,
		Ticket:  timeline.Ticket,
		Time:    timeline.Time,
		Updated: timeline.Updated,
	}

	s.hooks.OnRecordAfterCreateRequest.Publish(ctx, permission.TimelinesTable.ID, response)

	return openapi.CreateTimeline200JSONResponse(response), nil
}

func (s *Service) DeleteTimeline(ctx context.Context, request openapi.DeleteTimelineRequestObject) (openapi.DeleteTimelineResponseObject, error) {
	s.hooks.OnRecordBeforeDeleteRequest.Publish(ctx, permission.TimelinesTable.ID, request.Id)

	err := s.queries.DeleteTimeline(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterDeleteRequest.Publish(ctx, permission.TimelinesTable.ID, request.Id)

	return openapi.DeleteTimeline204Response{}, nil
}

func (s *Service) GetTimeline(ctx context.Context, request openapi.GetTimelineRequestObject) (openapi.GetTimelineResponseObject, error) {
	timeline, err := s.queries.GetTimeline(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	response := openapi.TimelineEntry{
		Id:      timeline.ID,
		Message: timeline.Message,
		Created: timeline.Created,
		Updated: timeline.Updated,
		Time:    timeline.Time,
		Ticket:  timeline.Ticket,
	}

	s.hooks.OnRecordViewRequest.Publish(ctx, permission.TimelinesTable.ID, response)

	return openapi.GetTimeline200JSONResponse(response), nil
}

func (s *Service) UpdateTimeline(ctx context.Context, request openapi.UpdateTimelineRequestObject) (openapi.UpdateTimelineResponseObject, error) {
	s.hooks.OnRecordBeforeUpdateRequest.Publish(ctx, permission.TimelinesTable.ID, request.Body)

	timeline, err := s.queries.UpdateTimeline(ctx, sqlc.UpdateTimelineParams{
		ID:      request.Id,
		Message: toNullString(request.Body.Message),
		Time:    toNullString(request.Body.Time),
	})
	if err != nil {
		return nil, err
	}

	response := openapi.TimelineEntry{
		Created: timeline.Created,
		Id:      timeline.ID,
		Message: timeline.Message,
		Ticket:  timeline.Ticket,
		Time:    timeline.Time,
		Updated: timeline.Updated,
	}

	s.hooks.OnRecordAfterUpdateRequest.Publish(ctx, permission.TimelinesTable.ID, response)

	return openapi.UpdateTimeline200JSONResponse(response), nil
}

func (s *Service) ListTypes(ctx context.Context, request openapi.ListTypesRequestObject) (openapi.ListTypesResponseObject, error) {
	types, err := s.queries.ListTypes(ctx, sqlc.ListTypesParams{
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
	})
	if err != nil {
		return nil, err
	}

	response := make([]openapi.Type, 0, len(types))
	for _, t := range types {
		response = append(response, openapi.Type{
			Created:  t.Created,
			Icon:     t.Icon,
			Id:       t.ID,
			Plural:   t.Plural,
			Schema:   unmarshal(t.Schema),
			Singular: t.Singular,
			Updated:  t.Updated,
		})
	}

	s.hooks.OnRecordsListRequest.Publish(ctx, permission.TypesTable.ID, response)

	totalCount := 0
	if len(types) > 0 {
		totalCount = int(types[0].TotalCount)
	}

	return openapi.ListTypes200JSONResponse{
		Body: response,
		Headers: openapi.ListTypes200ResponseHeaders{
			XTotalCount: totalCount,
		},
	}, nil
}

func (s *Service) CreateType(ctx context.Context, request openapi.CreateTypeRequestObject) (openapi.CreateTypeResponseObject, error) {
	s.hooks.OnRecordBeforeCreateRequest.Publish(ctx, permission.TypesTable.ID, request.Body)

	t, err := s.queries.CreateType(ctx, sqlc.CreateTypeParams{
		ID:       database.GenerateID("y"),
		Icon:     request.Body.Icon,
		Plural:   request.Body.Plural,
		Singular: request.Body.Singular,
		Schema:   marshal(request.Body.Schema),
	})
	if err != nil {
		return nil, err
	}

	response := openapi.Type{
		Created:  t.Created,
		Icon:     t.Icon,
		Id:       t.ID,
		Plural:   t.Plural,
		Schema:   unmarshal(t.Schema),
		Singular: t.Singular,
		Updated:  t.Updated,
	}

	s.hooks.OnRecordAfterCreateRequest.Publish(ctx, permission.TypesTable.ID, response)

	return openapi.CreateType200JSONResponse(response), nil
}

func (s *Service) DeleteType(ctx context.Context, request openapi.DeleteTypeRequestObject) (openapi.DeleteTypeResponseObject, error) {
	s.hooks.OnRecordBeforeDeleteRequest.Publish(ctx, permission.TypesTable.ID, request.Id)

	if err := s.queries.DeleteType(ctx, request.Id); err != nil {
		return nil, fmt.Errorf("failed to delete type: %w", err)
	}

	s.hooks.OnRecordAfterDeleteRequest.Publish(ctx, permission.TypesTable.ID, request.Id)

	return openapi.DeleteType204Response{}, nil
}

func (s *Service) GetType(ctx context.Context, request openapi.GetTypeRequestObject) (openapi.GetTypeResponseObject, error) {
	t, err := s.queries.GetType(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	response := openapi.Type{
		Created:  t.Created,
		Icon:     t.Icon,
		Id:       t.ID,
		Plural:   t.Plural,
		Schema:   unmarshal(t.Schema),
		Singular: t.Singular,
		Updated:  t.Updated,
	}

	s.hooks.OnRecordViewRequest.Publish(ctx, permission.TypesTable.ID, response)

	return openapi.GetType200JSONResponse(response), nil
}

func (s *Service) UpdateType(ctx context.Context, request openapi.UpdateTypeRequestObject) (openapi.UpdateTypeResponseObject, error) {
	s.hooks.OnRecordBeforeUpdateRequest.Publish(ctx, permission.TypesTable.ID, request.Body)

	t, err := s.queries.UpdateType(ctx, sqlc.UpdateTypeParams{
		ID:       request.Id,
		Icon:     toNullString(request.Body.Icon),
		Plural:   toNullString(request.Body.Plural),
		Singular: toNullString(request.Body.Singular),
		Schema:   marshalPointer(request.Body.Schema),
	})
	if err != nil {
		return nil, err
	}

	response := openapi.Type{
		Created:  t.Created,
		Icon:     t.Icon,
		Id:       t.ID,
		Plural:   t.Plural,
		Schema:   unmarshal(t.Schema),
		Singular: t.Singular,
		Updated:  t.Updated,
	}

	s.hooks.OnRecordAfterUpdateRequest.Publish(ctx, permission.TypesTable.ID, response)

	return openapi.UpdateType200JSONResponse(response), nil
}

func (s *Service) ListUsers(ctx context.Context, request openapi.ListUsersRequestObject) (openapi.ListUsersResponseObject, error) {
	users, err := s.queries.ListUsers(ctx, sqlc.ListUsersParams{
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
	})
	if err != nil {
		return nil, err
	}

	response := make([]openapi.User, 0, len(users))
	for _, user := range users {
		response = append(response, openapi.User{
			Avatar:                 user.Avatar,
			Created:                user.Created,
			Email:                  user.Email,
			Id:                     user.ID,
			LastResetSentAt:        user.Lastresetsentat,
			LastVerificationSentAt: user.Lastverificationsentat,
			Name:                   user.Name,
			Updated:                user.Updated,
			Username:               user.Username,
			Verified:               user.Verified,
		})
	}

	s.hooks.OnRecordsListRequest.Publish(ctx, permission.UsersTable.ID, response)

	totalCount := 0
	if len(users) > 0 {
		totalCount = int(users[0].TotalCount)
	}

	return openapi.ListUsers200JSONResponse{
		Body: response,
		Headers: openapi.ListUsers200ResponseHeaders{
			XTotalCount: totalCount,
		},
	}, nil
}

func (s *Service) CreateUser(ctx context.Context, request openapi.CreateUserRequestObject) (openapi.CreateUserResponseObject, error) {
	s.hooks.OnRecordBeforeCreateRequest.Publish(ctx, permission.UsersTable.ID, request.Body)

	tokenKey, err := password.GenerateTokenKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token key: %w", err)
	}

	user, err := s.queries.CreateUser(ctx, sqlc.CreateUserParams{
		ID:           database.GenerateID("u"),
		Name:         request.Body.Name,
		Email:        request.Body.Email,
		Username:     request.Body.Username,
		PasswordHash: "",
		TokenKey:     tokenKey,
		Avatar:       request.Body.Avatar,
		Verified:     request.Body.Verified,
	})
	if err != nil {
		return nil, err
	}

	response := openapi.User{
		Avatar:                 user.Avatar,
		Created:                user.Created,
		Email:                  user.Email,
		Id:                     user.ID,
		LastResetSentAt:        user.Lastresetsentat,
		LastVerificationSentAt: user.Lastverificationsentat,
		Name:                   user.Name,
		Updated:                user.Updated,
		Username:               user.Username,
		Verified:               user.Verified,
	}

	s.hooks.OnRecordAfterCreateRequest.Publish(ctx, permission.UsersTable.ID, response)

	return openapi.CreateUser200JSONResponse(response), nil
}

func (s *Service) DeleteUser(ctx context.Context, request openapi.DeleteUserRequestObject) (openapi.DeleteUserResponseObject, error) {
	s.hooks.OnRecordBeforeDeleteRequest.Publish(ctx, permission.UsersTable.ID, request.Id)

	err := s.queries.DeleteUser(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterDeleteRequest.Publish(ctx, permission.UsersTable.ID, request.Id)

	return openapi.DeleteUser204Response{}, nil
}

func (s *Service) GetUser(ctx context.Context, request openapi.GetUserRequestObject) (openapi.GetUserResponseObject, error) {
	user, err := s.queries.GetUser(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	response := openapi.User{
		Avatar:                 user.Avatar,
		Created:                user.Created,
		Email:                  user.Email,
		Id:                     user.ID,
		LastResetSentAt:        user.Lastresetsentat,
		LastVerificationSentAt: user.Lastverificationsentat,
		Name:                   user.Name,
		Updated:                user.Updated,
		Username:               user.Username,
		Verified:               user.Verified,
	}

	s.hooks.OnRecordViewRequest.Publish(ctx, permission.UsersTable.ID, response)

	return openapi.GetUser200JSONResponse(response), nil
}

func (s *Service) UpdateUser(ctx context.Context, request openapi.UpdateUserRequestObject) (openapi.UpdateUserResponseObject, error) {
	s.hooks.OnRecordBeforeUpdateRequest.Publish(ctx, permission.UsersTable.ID, request.Body)

	var passwordHash, tokenHash sql.NullString

	switch {
	case request.Body.Password == nil && request.Body.PasswordConfirm == nil:
	case request.Body.Password != nil && request.Body.PasswordConfirm != nil:
		if *request.Body.Password != *request.Body.PasswordConfirm {
			return nil, errors.New("passwords do not match")
		}

		passwordHashS, tokenHashS, err := password.Hash(*request.Body.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}

		passwordHash = sql.NullString{String: passwordHashS, Valid: true}
		tokenHash = sql.NullString{String: tokenHashS, Valid: true}
	default:
		return nil, errors.New("password and password confirm must be provided together")
	}

	user, err := s.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		Name:         toNullString(request.Body.Name),
		Email:        toNullString(request.Body.Email),
		Username:     toNullString(request.Body.Username),
		PasswordHash: passwordHash,
		TokenKey:     tokenHash,
		Avatar:       toNullString(request.Body.Avatar),
		Verified:     toNullBool(request.Body.Verified),
		ID:           request.Id,
	})
	if err != nil {
		return nil, err
	}

	response := openapi.User{
		Avatar:                 user.Avatar,
		Created:                user.Created,
		Email:                  user.Email,
		Id:                     user.ID,
		LastResetSentAt:        user.Lastresetsentat,
		LastVerificationSentAt: user.Lastverificationsentat,
		Name:                   user.Name,
		Updated:                user.Updated,
		Username:               user.Username,
		Verified:               user.Verified,
	}

	s.hooks.OnRecordAfterUpdateRequest.Publish(ctx, "users", response)

	return openapi.UpdateUser200JSONResponse(response), nil
}

func (s *Service) ListGroups(ctx context.Context, request openapi.ListGroupsRequestObject) (openapi.ListGroupsResponseObject, error) {
	groups, err := s.queries.ListGroups(ctx, sqlc.ListGroupsParams{
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
	})
	if err != nil {
		return nil, err
	}

	response := make([]openapi.Group, 0, len(groups))

	for _, group := range groups {
		response = append(response, openapi.Group{
			Created:     group.Created,
			Id:          group.ID,
			Name:        group.Name,
			Permissions: permission.FromJSONArray(ctx, group.Permissions),
			Updated:     group.Updated,
		})
	}

	s.hooks.OnRecordsListRequest.Publish(ctx, permission.GroupsTable.ID, response)

	totalCount := 0
	if len(groups) > 0 {
		totalCount = int(groups[0].TotalCount)
	}

	return openapi.ListGroups200JSONResponse{
		Body: response,
		Headers: openapi.ListGroups200ResponseHeaders{
			XTotalCount: totalCount,
		},
	}, nil
}

func (s *Service) CreateGroup(ctx context.Context, request openapi.CreateGroupRequestObject) (openapi.CreateGroupResponseObject, error) {
	s.hooks.OnRecordBeforeCreateRequest.Publish(ctx, permission.GroupsTable.ID, request.Body)

	group, err := s.queries.CreateGroup(ctx, sqlc.CreateGroupParams{
		ID:          database.GenerateID("g"),
		Name:        request.Body.Name,
		Permissions: permission.ToJSONArray(ctx, request.Body.Permissions),
	})
	if err != nil {
		return nil, err
	}

	response := openapi.Group{
		Created:     group.Created,
		Id:          group.ID,
		Name:        group.Name,
		Permissions: permission.FromJSONArray(ctx, group.Permissions),
		Updated:     group.Updated,
	}

	s.hooks.OnRecordAfterCreateRequest.Publish(ctx, permission.GroupsTable.ID, response)

	return openapi.CreateGroup200JSONResponse(response), nil
}

func (s *Service) DeleteGroup(ctx context.Context, request openapi.DeleteGroupRequestObject) (openapi.DeleteGroupResponseObject, error) {
	s.hooks.OnRecordBeforeDeleteRequest.Publish(ctx, permission.GroupsTable.ID, request.Id)

	if request.Id == "admin" {
		return nil, errors.New("cannot delete the admin group")
	}

	err := s.queries.DeleteGroup(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterDeleteRequest.Publish(ctx, permission.GroupsTable.ID, request.Id)

	return openapi.DeleteGroup204Response{}, nil
}

func (s *Service) GetGroup(ctx context.Context, request openapi.GetGroupRequestObject) (openapi.GetGroupResponseObject, error) {
	group, err := s.queries.GetGroup(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	response := openapi.Group{
		Created:     group.Created,
		Id:          group.ID,
		Name:        group.Name,
		Permissions: permission.FromJSONArray(ctx, group.Permissions),
		Updated:     group.Updated,
	}

	s.hooks.OnRecordViewRequest.Publish(ctx, permission.GroupsTable.ID, response)

	return openapi.GetGroup200JSONResponse(response), nil
}

func (s *Service) UpdateGroup(ctx context.Context, request openapi.UpdateGroupRequestObject) (openapi.UpdateGroupResponseObject, error) {
	s.hooks.OnRecordBeforeUpdateRequest.Publish(ctx, permission.GroupsTable.ID, request.Body)

	if request.Id == "admin" {
		return nil, errors.New("cannot update the admin group")
	}

	var permissions sql.NullString

	if request.Body.Permissions != nil {
		permissions = sql.NullString{String: permission.ToJSONArray(ctx, *request.Body.Permissions), Valid: true}
	}

	group, err := s.queries.UpdateGroup(ctx, sqlc.UpdateGroupParams{
		Name:        toNullString(request.Body.Name),
		Permissions: permissions,
		ID:          request.Id,
	})
	if err != nil {
		return nil, err
	}

	response := openapi.Group{
		Created:     group.Created,
		Id:          group.ID,
		Name:        group.Name,
		Permissions: permission.FromJSONArray(ctx, group.Permissions),
		Updated:     group.Updated,
	}

	s.hooks.OnRecordAfterUpdateRequest.Publish(ctx, permission.GroupsTable.ID, response)

	return openapi.UpdateGroup200JSONResponse(response), nil
}

func (s *Service) AddGroupParent(ctx context.Context, request openapi.AddGroupParentRequestObject) (openapi.AddGroupParentResponseObject, error) {
	s.hooks.OnRecordBeforeCreateRequest.Publish(ctx, permission.GroupParentTable.ID, request.Body)

	err := s.queries.AssignParentGroup(ctx, sqlc.AssignParentGroupParams{
		ChildGroupID:  request.Id,
		ParentGroupID: request.Body.GroupId,
	})
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterCreateRequest.Publish(ctx, permission.GroupParentTable.ID, request.Body)

	return openapi.AddGroupParent201Response{}, nil
}

func (s *Service) RemoveGroupParent(ctx context.Context, request openapi.RemoveGroupParentRequestObject) (openapi.RemoveGroupParentResponseObject, error) {
	s.hooks.OnRecordBeforeDeleteRequest.Publish(ctx, permission.GroupParentTable.ID, request.Id)

	err := s.queries.RemoveParentGroup(ctx, sqlc.RemoveParentGroupParams{
		ChildGroupID:  request.Id,
		ParentGroupID: request.ParentGroupId,
	})
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterDeleteRequest.Publish(ctx, permission.GroupParentTable.ID, request.Id)

	return openapi.RemoveGroupParent204Response{}, nil
}

func (s *Service) AddUserGroup(ctx context.Context, request openapi.AddUserGroupRequestObject) (openapi.AddUserGroupResponseObject, error) {
	s.hooks.OnRecordBeforeCreateRequest.Publish(ctx, permission.UserGroupTable.ID, request.Body)

	err := s.queries.AssignGroupToUser(ctx, sqlc.AssignGroupToUserParams{
		UserID:  request.Id,
		GroupID: request.Body.GroupId,
	})
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterCreateRequest.Publish(ctx, permission.UserGroupTable.ID, request.Body)

	return openapi.AddUserGroup201Response{}, nil
}

func (s *Service) RemoveUserGroup(ctx context.Context, request openapi.RemoveUserGroupRequestObject) (openapi.RemoveUserGroupResponseObject, error) {
	s.hooks.OnRecordBeforeDeleteRequest.Publish(ctx, permission.UserGroupTable.ID, request.Id)

	err := s.queries.RemoveGroupFromUser(ctx, sqlc.RemoveGroupFromUserParams{
		UserID:  request.Id,
		GroupID: request.GroupId,
	})
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterDeleteRequest.Publish(ctx, permission.UserGroupTable.ID, request.Id)

	return openapi.RemoveUserGroup204Response{}, nil
}

func (s *Service) ListUserPermissions(ctx context.Context, request openapi.ListUserPermissionsRequestObject) (openapi.ListUserPermissionsResponseObject, error) {
	permissions, err := s.queries.ListUserPermissions(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordsListRequest.Publish(ctx, permission.UserPermissionTable.ID, permissions)

	return openapi.ListUserPermissions200JSONResponse(permissions), nil
}

func (s *Service) ListUserGroups(ctx context.Context, request openapi.ListUserGroupsRequestObject) (openapi.ListUserGroupsResponseObject, error) {
	groups, err := s.queries.ListUserGroups(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	response := make([]openapi.UserGroup, 0, len(groups))
	for _, group := range groups {
		response = append(response, openapi.UserGroup{
			Created:     group.Created,
			Id:          group.ID,
			Name:        group.Name,
			Permissions: permission.FromJSONArray(ctx, group.Permissions),
			Type:        group.GroupType,
			Updated:     group.Updated,
		})
	}

	s.hooks.OnRecordsListRequest.Publish(ctx, permission.UserGroupTable.ID, response)

	return openapi.ListUserGroups200JSONResponse(response), nil
}

func (s *Service) ListGroupUsers(ctx context.Context, request openapi.ListGroupUsersRequestObject) (openapi.ListGroupUsersResponseObject, error) {
	users, err := s.queries.ListGroupUsers(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	response := make([]openapi.GroupUser, 0, len(users))
	for _, user := range users {
		response = append(response, openapi.GroupUser{
			Avatar:                 user.Avatar,
			Created:                user.Created,
			Email:                  user.Email,
			Id:                     user.ID,
			LastResetSentAt:        user.Lastresetsentat,
			LastVerificationSentAt: user.Lastverificationsentat,
			Name:                   user.Name,
			Updated:                user.Updated,
			Username:               user.Username,
			Verified:               user.Verified,
			Type:                   user.GroupType,
		})
	}

	s.hooks.OnRecordsListRequest.Publish(ctx, permission.GroupUserTable.ID, response)

	return openapi.ListGroupUsers200JSONResponse(response), nil
}

func (s *Service) ListParentPermissions(ctx context.Context, request openapi.ListParentPermissionsRequestObject) (openapi.ListParentPermissionsResponseObject, error) {
	permissions, err := s.queries.ListParentPermissions(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordsListRequest.Publish(ctx, permission.GroupPermissionTable.ID, permissions)

	return openapi.ListParentPermissions200JSONResponse(permissions), nil
}

func (s *Service) ListParentGroups(ctx context.Context, request openapi.ListParentGroupsRequestObject) (openapi.ListParentGroupsResponseObject, error) {
	groups, err := s.queries.ListParentGroups(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	response := make([]openapi.UserGroup, 0, len(groups))
	for _, group := range groups {
		response = append(response, openapi.UserGroup{
			Created:     group.Created,
			Id:          group.ID,
			Name:        group.Name,
			Permissions: permission.FromJSONArray(ctx, group.Permissions),
			Type:        group.GroupType,
			Updated:     group.Updated,
		})
	}

	s.hooks.OnRecordsListRequest.Publish(ctx, permission.GroupParentTable.ID, response)

	return openapi.ListParentGroups200JSONResponse(response), nil
}

func (s *Service) ListChildGroups(ctx context.Context, request openapi.ListChildGroupsRequestObject) (openapi.ListChildGroupsResponseObject, error) {
	groups, err := s.queries.ListChildGroups(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	response := make([]openapi.UserGroup, 0, len(groups))
	for _, group := range groups {
		response = append(response, openapi.UserGroup{
			Created:     group.Created,
			Id:          group.ID,
			Name:        group.Name,
			Permissions: permission.FromJSONArray(ctx, group.Permissions),
			Type:        group.GroupType,
			Updated:     group.Updated,
		})
	}

	s.hooks.OnRecordsListRequest.Publish(ctx, permission.GroupChildTable.ID, response)

	return openapi.ListChildGroups200JSONResponse(response), nil
}

func (s *Service) ListWebhooks(ctx context.Context, request openapi.ListWebhooksRequestObject) (openapi.ListWebhooksResponseObject, error) {
	webhooks, err := s.queries.ListWebhooks(ctx, sqlc.ListWebhooksParams{
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
	})
	if err != nil {
		return nil, err
	}

	response := make([]openapi.Webhook, 0, len(webhooks))
	for _, webhook := range webhooks {
		response = append(response, openapi.Webhook{
			Id:   webhook.ID,
			Name: webhook.Name,
		})
	}

	s.hooks.OnRecordsListRequest.Publish(ctx, permission.WebhooksTable.ID, response)

	totalCount := 0
	if len(webhooks) > 0 {
		totalCount = int(webhooks[0].TotalCount)
	}

	return openapi.ListWebhooks200JSONResponse{
		Body: response,
		Headers: openapi.ListWebhooks200ResponseHeaders{
			XTotalCount: totalCount,
		},
	}, nil
}

func (s *Service) CreateWebhook(ctx context.Context, request openapi.CreateWebhookRequestObject) (openapi.CreateWebhookResponseObject, error) {
	s.hooks.OnRecordBeforeCreateRequest.Publish(ctx, permission.WebhooksTable.ID, request.Body)

	webhook, err := s.queries.CreateWebhook(ctx, sqlc.CreateWebhookParams{
		ID:          database.GenerateID("w"),
		Name:        request.Body.Name,
		Destination: request.Body.Destination,
		Collection:  request.Body.Collection,
	})
	if err != nil {
		return nil, err
	}

	response := openapi.Webhook{
		Id:          webhook.ID,
		Name:        webhook.Name,
		Created:     webhook.Created,
		Updated:     webhook.Updated,
		Destination: webhook.Destination,
		Collection:  webhook.Collection,
	}

	s.hooks.OnRecordAfterCreateRequest.Publish(ctx, permission.WebhooksTable.ID, response)

	return openapi.CreateWebhook200JSONResponse(response), nil
}

func (s *Service) DeleteWebhook(ctx context.Context, request openapi.DeleteWebhookRequestObject) (openapi.DeleteWebhookResponseObject, error) {
	s.hooks.OnRecordBeforeDeleteRequest.Publish(ctx, permission.WebhooksTable.ID, request.Id)

	err := s.queries.DeleteWebhook(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterDeleteRequest.Publish(ctx, permission.WebhooksTable.ID, request.Id)

	return openapi.DeleteWebhook204Response{}, nil
}

func (s *Service) GetWebhook(ctx context.Context, request openapi.GetWebhookRequestObject) (openapi.GetWebhookResponseObject, error) {
	webhook, err := s.queries.GetWebhook(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	response := openapi.Webhook{
		Id:          webhook.ID,
		Name:        webhook.Name,
		Created:     webhook.Created,
		Updated:     webhook.Updated,
		Destination: webhook.Destination,
		Collection:  webhook.Collection,
	}

	s.hooks.OnRecordViewRequest.Publish(ctx, permission.WebhooksTable.ID, response)

	return openapi.GetWebhook200JSONResponse(response), nil
}

func (s *Service) UpdateWebhook(ctx context.Context, request openapi.UpdateWebhookRequestObject) (openapi.UpdateWebhookResponseObject, error) {
	s.hooks.OnRecordBeforeUpdateRequest.Publish(ctx, permission.WebhooksTable.ID, request.Body)

	webhook, err := s.queries.UpdateWebhook(ctx, sqlc.UpdateWebhookParams{
		ID:          request.Id,
		Name:        toNullString(request.Body.Name),
		Destination: toNullString(request.Body.Destination),
		Collection:  toNullString(request.Body.Collection),
	})
	if err != nil {
		return nil, err
	}

	response := openapi.Webhook{
		Id:          webhook.ID,
		Name:        webhook.Name,
		Created:     webhook.Created,
		Updated:     webhook.Updated,
		Destination: webhook.Destination,
		Collection:  webhook.Collection,
	}

	s.hooks.OnRecordAfterUpdateRequest.Publish(ctx, permission.WebhooksTable.ID, response)

	return openapi.UpdateWebhook200JSONResponse(response), nil
}

func (s *Service) GetConfig(ctx context.Context, _ openapi.GetConfigRequestObject) (openapi.GetConfigResponseObject, error) {
	flags := []string{}

	features, err := s.queries.ListFeatures(ctx, sqlc.ListFeaturesParams{
		Offset: 0,
		Limit:  1000,
	})
	if err != nil {
		return nil, err
	}

	for _, feature := range features {
		flags = append(flags, feature.Name)
	}

	tables := []openapi.Table{}
	for _, table := range permission.Tables() {
		tables = append(tables, openapi.Table{
			Id:   table.ID,
			Name: table.Name,
		})
	}

	response := openapi.Config{
		Flags:       flags,
		Permissions: permission.All(),
		Tables:      tables,
	}

	return openapi.GetConfig200JSONResponse(response), nil
}

func (s *Service) GetSettings(ctx context.Context, _ openapi.GetSettingsRequestObject) (openapi.GetSettingsResponseObject, error) {
	settings, err := database.LoadSettings(ctx, s.queries)
	if err != nil {
		return nil, err
	}

	return openapi.GetSettings200JSONResponse(mapSettings(settings)), nil
}

func (s *Service) UpdateSettings(ctx context.Context, request openapi.UpdateSettingsRequestObject) (openapi.UpdateSettingsResponseObject, error) {
	settings, err := database.LoadSettings(ctx, s.queries)
	if err != nil {
		return nil, err
	}

	settings.Meta.AppName = request.Body.Meta.AppName
	settings.Meta.AppURL = request.Body.Meta.AppUrl
	// settings.Meta.HideControls = request.Body.Meta.HideControls
	settings.Meta.SenderAddress = request.Body.Meta.SenderAddress
	settings.Meta.SenderName = request.Body.Meta.SenderName
	// settings.Meta.ConfirmEmailChangeTemplate.ActionURL = request.Body.Meta.ConfirmEmailChangeTemplate.ActionUrl
	// settings.Meta.ConfirmEmailChangeTemplate.Body = request.Body.Meta.ConfirmEmailChangeTemplate.Body
	// settings.Meta.ConfirmEmailChangeTemplate.Hidden = request.Body.Meta.ConfirmEmailChangeTemplate.Hidden
	// settings.Meta.ConfirmEmailChangeTemplate.Subject = request.Body.Meta.ConfirmEmailChangeTemplate.Subject
	// settings.Meta.ResetPasswordTemplate.ActionURL = request.Body.Meta.ResetPasswordTemplate.ActionUrl
	// settings.Meta.ResetPasswordTemplate.Body = request.Body.Meta.ResetPasswordTemplate.Body
	// settings.Meta.ResetPasswordTemplate.Hidden = request.Body.Meta.ResetPasswordTemplate.Hidden
	// settings.Meta.ResetPasswordTemplate.Subject = request.Body.Meta.ResetPasswordTemplate.Subject
	// settings.Meta.VerificationTemplate.ActionURL = request.Body.Meta.VerificationTemplate.ActionUrl
	// settings.Meta.VerificationTemplate.Body = request.Body.Meta.VerificationTemplate.Body
	// settings.Meta.VerificationTemplate.Hidden = request.Body.Meta.VerificationTemplate.Hidden
	// settings.Meta.VerificationTemplate.Subject = request.Body.Meta.VerificationTemplate.Subject
	// settings.Logs.LogIP = request.Body.Logs.LogIp
	// settings.Logs.MaxDays = request.Body.Logs.MaxDays
	// settings.Logs.MinLevel = request.Body.Logs.MinLevel
	settings.SMTP.Enabled = request.Body.Smtp.Enabled
	settings.SMTP.Host = request.Body.Smtp.Host
	settings.SMTP.Port = request.Body.Smtp.Port
	settings.SMTP.Username = request.Body.Smtp.Username
	settings.SMTP.Password = request.Body.Smtp.Password
	settings.SMTP.AuthMethod = request.Body.Smtp.AuthMethod
	settings.SMTP.TLS = request.Body.Smtp.Tls
	settings.SMTP.LocalName = request.Body.Smtp.LocalName
	// settings.S3.Enabled = request.Body.S3.Enabled
	// settings.S3.Bucket = request.Body.S3.Bucket
	// settings.S3.Endpoint = request.Body.S3.Endpoint
	// settings.S3.ForcePathStyle = request.Body.S3.ForcePathStyle
	// settings.S3.Region = request.Body.S3.Region
	// settings.S3.Secret = request.Body.S3.Secret
	// settings.S3.AccessKey = request.Body.S3.AccessKey
	// settings.Backups.CronMaxKeep = request.Body.Backups.CronMaxKeep
	// settings.Backups.Cron = request.Body.Backups.Cron
	// settings.Backups.S3.Enabled = request.Body.Backups.S3.Enabled
	// settings.Backups.S3.Bucket = request.Body.Backups.S3.Bucket
	// settings.Backups.S3.Endpoint = request.Body.Backups.S3.Endpoint
	// settings.Backups.S3.ForcePathStyle = request.Body.Backups.S3.ForcePathStyle
	// settings.Backups.S3.Region = request.Body.Backups.S3.Region
	// settings.Backups.S3.AccessKey = request.Body.Backups.S3.AccessKey
	// settings.Backups.S3.Secret = request.Body.Backups.S3.Secret

	if err := database.SaveSettings(ctx, s.queries, settings); err != nil {
		return nil, fmt.Errorf("failed to save settings: %w", err)
	}

	return openapi.UpdateSettings200JSONResponse(mapSettings(settings)), err
}

func toString(value *string, defaultValue string) string {
	if value == nil {
		return defaultValue
	}

	return *value
}

func toNullString(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}

	return sql.NullString{String: *value, Valid: true}
}

func toInt64(value *int, defaultValue int64) int64 {
	if value == nil {
		return defaultValue
	}

	return int64(*value)
}

func toNullBool(value *bool) sql.NullBool {
	if value == nil {
		return sql.NullBool{}
	}

	return sql.NullBool{Bool: *value, Valid: true}
}

func marshal(state map[string]interface{}) string {
	b, _ := json.Marshal(state) //nolint:errchkjson

	return string(b)
}

func marshalPointer(state *map[string]interface{}) string {
	if state == nil {
		return "{}"
	}

	b, _ := json.Marshal(*state) //nolint:errchkjson

	return string(b)
}

func unmarshal(data string) map[string]interface{} {
	var m map[string]interface{}

	if err := json.Unmarshal([]byte(data), &m); err != nil {
		return nil
	}

	return m
}

func mapSettings(settings *database.Settings) openapi.Settings {
	return openapi.Settings{
		Meta: openapi.SettingsMeta{
			AppName:       settings.Meta.AppName,
			AppUrl:        settings.Meta.AppURL,
			HideControls:  settings.Meta.HideControls,
			SenderAddress: settings.Meta.SenderAddress,
			SenderName:    settings.Meta.SenderName,
			ConfirmEmailChangeTemplate: openapi.EmailTemplate{
				ActionUrl: settings.Meta.ConfirmEmailChangeTemplate.ActionURL,
				Body:      settings.Meta.ConfirmEmailChangeTemplate.Body,
				Hidden:    settings.Meta.ConfirmEmailChangeTemplate.Hidden,
				Subject:   settings.Meta.ConfirmEmailChangeTemplate.Subject,
			},
			ResetPasswordTemplate: openapi.EmailTemplate{
				ActionUrl: settings.Meta.ResetPasswordTemplate.ActionURL,
				Body:      settings.Meta.ResetPasswordTemplate.Body,
				Hidden:    settings.Meta.ResetPasswordTemplate.Hidden,
				Subject:   settings.Meta.ResetPasswordTemplate.Subject,
			},
			VerificationTemplate: openapi.EmailTemplate{
				ActionUrl: settings.Meta.VerificationTemplate.ActionURL,
				Body:      settings.Meta.VerificationTemplate.Body,
				Hidden:    settings.Meta.VerificationTemplate.Hidden,
				Subject:   settings.Meta.VerificationTemplate.Subject,
			},
		},
		Backups: openapi.SettingsBackups{
			Cron:        settings.Backups.Cron,
			CronMaxKeep: settings.Backups.CronMaxKeep,
			S3: openapi.S3Config{
				AccessKey:      settings.Backups.S3.AccessKey,
				Bucket:         settings.Backups.S3.Bucket,
				Enabled:        settings.Backups.S3.Enabled,
				Endpoint:       settings.Backups.S3.Endpoint,
				ForcePathStyle: settings.Backups.S3.ForcePathStyle,
				Region:         settings.Backups.S3.Region,
				Secret:         settings.Backups.S3.Secret,
			},
		},
		Logs: openapi.SettingsLogs{
			LogIp:    settings.Logs.LogIP,
			MaxDays:  settings.Logs.MaxDays,
			MinLevel: settings.Logs.MinLevel,
		},
		S3: openapi.S3Config{
			AccessKey:      settings.S3.AccessKey,
			Bucket:         settings.S3.Bucket,
			Enabled:        settings.S3.Enabled,
			Endpoint:       settings.S3.Endpoint,
			ForcePathStyle: settings.S3.ForcePathStyle,
			Region:         settings.S3.Region,
			Secret:         settings.S3.Secret,
		},
		Smtp: openapi.SettingsSmtp{
			AuthMethod: settings.SMTP.AuthMethod,
			Enabled:    settings.SMTP.Enabled,
			Host:       settings.SMTP.Host,
			LocalName:  settings.SMTP.LocalName,
			Password:   settings.SMTP.Password,
			Port:       settings.SMTP.Port,
			Tls:        settings.SMTP.TLS,
			Username:   settings.SMTP.Username,
		},
	}
}
