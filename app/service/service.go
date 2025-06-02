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

	"github.com/google/uuid"

	"github.com/SecurityBrewery/catalyst/app/auth"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/hook"
	"github.com/SecurityBrewery/catalyst/app/openapi"
	"github.com/SecurityBrewery/catalyst/permission"
	"github.com/SecurityBrewery/catalyst/reaction/schedule"
)

const (
	defaultLimit  = 100
	defaultOffset = 0
)

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

	s.hooks.OnRecordsListRequest.Publish(ctx, "comments", response)

	return openapi.ListComments200JSONResponse{
		Body: response,
		Headers: openapi.ListComments200ResponseHeaders{
			XTotalCount: totalCount,
		},
	}, nil
}

func (s *Service) CreateComment(ctx context.Context, request openapi.CreateCommentRequestObject) (openapi.CreateCommentResponseObject, error) {
	s.hooks.OnRecordBeforeCreateRequest.Publish(ctx, "comments", request.Body)

	comment, err := s.queries.CreateComment(ctx, sqlc.CreateCommentParams{
		ID:      generateID("c"),
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

	s.hooks.OnRecordAfterCreateRequest.Publish(ctx, "comments", response)

	return openapi.CreateComment200JSONResponse(response), nil
}

func (s *Service) DeleteComment(ctx context.Context, request openapi.DeleteCommentRequestObject) (openapi.DeleteCommentResponseObject, error) {
	s.hooks.OnRecordBeforeDeleteRequest.Publish(ctx, "comments", request.Id)

	err := s.queries.DeleteComment(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterDeleteRequest.Publish(ctx, "comments", request.Id)

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

	s.hooks.OnRecordViewRequest.Publish(ctx, "comments", response)

	return openapi.GetComment200JSONResponse(response), nil
}

func (s *Service) UpdateComment(ctx context.Context, request openapi.UpdateCommentRequestObject) (openapi.UpdateCommentResponseObject, error) {
	s.hooks.OnRecordBeforeUpdateRequest.Publish(ctx, "comments", request.Body)

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

	s.hooks.OnRecordAfterUpdateRequest.Publish(ctx, "comments", response)

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

	s.hooks.OnRecordsListRequest.Publish(ctx, "dashboard_counts", response)

	return openapi.GetDashboardCounts200JSONResponse(response), nil
}

func (s *Service) ListFeatures(ctx context.Context, request openapi.ListFeaturesRequestObject) (openapi.ListFeaturesResponseObject, error) {
	features, err := s.queries.ListFeatures(ctx, sqlc.ListFeaturesParams{
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
	})
	if err != nil {
		return nil, err
	}

	response := make([]openapi.Feature, 0, len(features))
	for _, feature := range features {
		response = append(response, openapi.Feature{
			Id:      feature.ID,
			Name:    feature.Name,
			Created: feature.Created,
			Updated: feature.Updated,
		})
	}

	totalCount := 0
	if len(features) > 0 {
		totalCount = int(features[0].TotalCount)
	}

	s.hooks.OnRecordsListRequest.Publish(ctx, "features", response)

	return openapi.ListFeatures200JSONResponse{
		Body: response,
		Headers: openapi.ListFeatures200ResponseHeaders{
			XTotalCount: totalCount,
		},
	}, nil
}

func (s *Service) CreateFeature(ctx context.Context, request openapi.CreateFeatureRequestObject) (openapi.CreateFeatureResponseObject, error) {
	s.hooks.OnRecordBeforeCreateRequest.Publish(ctx, "features", request.Body)

	feature, err := s.queries.CreateFeature(ctx, request.Body.Name)
	if err != nil {
		return nil, err
	}

	response := openapi.Feature{
		Id:      feature.ID,
		Name:    feature.Name,
		Created: feature.Created,
		Updated: feature.Updated,
	}

	s.hooks.OnRecordAfterCreateRequest.Publish(ctx, "features", response)

	return openapi.CreateFeature200JSONResponse(response), nil
}

func (s *Service) DeleteFeature(ctx context.Context, request openapi.DeleteFeatureRequestObject) (openapi.DeleteFeatureResponseObject, error) {
	s.hooks.OnRecordBeforeDeleteRequest.Publish(ctx, "features", request.Id)

	err := s.queries.DeleteFeature(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterDeleteRequest.Publish(ctx, "features", request.Id)

	return openapi.DeleteFeature204Response{}, nil
}

func (s *Service) GetFeature(ctx context.Context, request openapi.GetFeatureRequestObject) (openapi.GetFeatureResponseObject, error) {
	feature, err := s.queries.GetFeature(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	response := openapi.Feature{
		Id:      feature.ID,
		Name:    feature.Name,
		Created: feature.Created,
		Updated: feature.Updated,
	}

	s.hooks.OnRecordViewRequest.Publish(ctx, "features", response)

	return openapi.GetFeature200JSONResponse(response), nil
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
	s.hooks.OnRecordBeforeCreateRequest.Publish(ctx, "files", request.Body)

	file, err := s.queries.CreateFile(ctx, sqlc.CreateFileParams{
		ID:     generateID("f"),
		Name:   request.Body.Name,
		Blob:   request.Body.Blob,
		Size:   request.Body.Size,
		Ticket: request.Body.Ticket,
	})
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterCreateRequest.Publish(ctx, "files", file)

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
	s.hooks.OnRecordBeforeDeleteRequest.Publish(ctx, "files", request.Id)

	err := s.queries.DeleteFile(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterDeleteRequest.Publish(ctx, "files", request.Id)

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

	s.hooks.OnRecordViewRequest.Publish(ctx, "files", response)

	return openapi.GetFile200JSONResponse(response), nil
}

func (s *Service) UpdateFile(ctx context.Context, request openapi.UpdateFileRequestObject) (openapi.UpdateFileResponseObject, error) {
	s.hooks.OnRecordBeforeUpdateRequest.Publish(ctx, "files", request.Body)

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

	s.hooks.OnRecordAfterUpdateRequest.Publish(ctx, "files", response)

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
	s.hooks.OnRecordBeforeCreateRequest.Publish(ctx, "links", request.Body)

	link, err := s.queries.CreateLink(ctx, sqlc.CreateLinkParams{
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

	s.hooks.OnRecordAfterCreateRequest.Publish(ctx, "links", response)

	return openapi.CreateLink200JSONResponse(response), nil
}

func (s *Service) DeleteLink(ctx context.Context, request openapi.DeleteLinkRequestObject) (openapi.DeleteLinkResponseObject, error) {
	s.hooks.OnRecordBeforeDeleteRequest.Publish(ctx, "links", request.Id)

	err := s.queries.DeleteLink(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterDeleteRequest.Publish(ctx, "links", request.Id)

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

	s.hooks.OnRecordViewRequest.Publish(ctx, "links", response)

	return openapi.GetLink200JSONResponse(response), nil
}

func (s *Service) UpdateLink(ctx context.Context, request openapi.UpdateLinkRequestObject) (openapi.UpdateLinkResponseObject, error) {
	s.hooks.OnRecordBeforeUpdateRequest.Publish(ctx, "links", request.Body)

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

	s.hooks.OnRecordAfterUpdateRequest.Publish(ctx, "links", response)

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

	s.hooks.OnRecordsListRequest.Publish(ctx, "reactions", response)

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
	s.hooks.OnRecordBeforeCreateRequest.Publish(ctx, "reactions", request.Body)

	reaction, err := s.queries.CreateReaction(ctx, sqlc.CreateReactionParams{
		ID:      generateID("r"),
		Name:    request.Body.Name,
		Action:  request.Body.Action,
		Trigger: request.Body.Trigger,
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

	s.hooks.OnRecordAfterCreateRequest.Publish(ctx, "reactions", response)

	return openapi.CreateReaction200JSONResponse(response), nil
}

func (s *Service) DeleteReaction(ctx context.Context, request openapi.DeleteReactionRequestObject) (openapi.DeleteReactionResponseObject, error) {
	s.hooks.OnRecordBeforeDeleteRequest.Publish(ctx, "reactions", request.Id)

	err := s.queries.DeleteReaction(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.scheduler.RemoveReaction(request.Id)

	s.hooks.OnRecordAfterDeleteRequest.Publish(ctx, "reactions", request.Id)

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

	s.hooks.OnRecordViewRequest.Publish(ctx, "reactions", response)

	return openapi.GetReaction200JSONResponse(response), nil
}

func (s *Service) UpdateReaction(ctx context.Context, request openapi.UpdateReactionRequestObject) (openapi.UpdateReactionResponseObject, error) {
	s.hooks.OnRecordBeforeUpdateRequest.Publish(ctx, "reactions", request.Body)

	reaction, err := s.queries.UpdateReaction(ctx, sqlc.UpdateReactionParams{
		ID:      request.Id,
		Name:    toNullString(request.Body.Name),
		Action:  toNullString(request.Body.Action),
		Trigger: toNullString(request.Body.Trigger),
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

	s.hooks.OnRecordAfterUpdateRequest.Publish(ctx, "reactions", response)

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

	s.hooks.OnRecordsListRequest.Publish(ctx, "sidebar", response)

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

	s.hooks.OnRecordsListRequest.Publish(ctx, "tasks", response)

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
	s.hooks.OnRecordBeforeCreateRequest.Publish(ctx, "tasks", request.Body)

	task, err := s.queries.CreateTask(ctx, sqlc.CreateTaskParams{
		ID:     generateID("t"),
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

	s.hooks.OnRecordAfterCreateRequest.Publish(ctx, "tasks", response)

	return openapi.CreateTask200JSONResponse(response), nil
}

func (s *Service) DeleteTask(ctx context.Context, request openapi.DeleteTaskRequestObject) (openapi.DeleteTaskResponseObject, error) {
	s.hooks.OnRecordBeforeDeleteRequest.Publish(ctx, "tasks", request.Id)

	err := s.queries.DeleteTask(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterDeleteRequest.Publish(ctx, "tasks", request.Id)

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

	s.hooks.OnRecordViewRequest.Publish(ctx, "tasks", response)

	return openapi.GetTask200JSONResponse(response), nil
}

func (s *Service) UpdateTask(ctx context.Context, request openapi.UpdateTaskRequestObject) (openapi.UpdateTaskResponseObject, error) {
	s.hooks.OnRecordBeforeUpdateRequest.Publish(ctx, "tasks", request.Body)

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

	s.hooks.OnRecordAfterUpdateRequest.Publish(ctx, "tasks", response)

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

	s.hooks.OnRecordsListRequest.Publish(ctx, "tickets", response)

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
	s.hooks.OnRecordBeforeCreateRequest.Publish(ctx, "tickets", request.Body)

	ticket, err := s.queries.CreateTicket(ctx, sqlc.CreateTicketParams{
		ID:          generateID(request.Body.Type),
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

	s.hooks.OnRecordAfterCreateRequest.Publish(ctx, "tickets", response)

	return openapi.CreateTicket200JSONResponse(response), nil
}

func (s *Service) DeleteTicket(ctx context.Context, request openapi.DeleteTicketRequestObject) (openapi.DeleteTicketResponseObject, error) {
	s.hooks.OnRecordBeforeDeleteRequest.Publish(ctx, "tickets", request.Id)

	err := s.queries.DeleteTicket(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterDeleteRequest.Publish(ctx, "tickets", request.Id)

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

	s.hooks.OnRecordViewRequest.Publish(ctx, "tickets", response)

	return openapi.GetTicket200JSONResponse(response), nil
}

func (s *Service) UpdateTicket(ctx context.Context, request openapi.UpdateTicketRequestObject) (openapi.UpdateTicketResponseObject, error) {
	s.hooks.OnRecordBeforeUpdateRequest.Publish(ctx, "tickets", request.Body)

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

	s.hooks.OnRecordAfterUpdateRequest.Publish(ctx, "tickets", response)

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

	s.hooks.OnRecordsListRequest.Publish(ctx, "timeline", response)

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
	s.hooks.OnRecordBeforeCreateRequest.Publish(ctx, "timeline", request.Body)

	timeline, err := s.queries.CreateTimeline(ctx, sqlc.CreateTimelineParams{
		ID:      generateID("h"),
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

	s.hooks.OnRecordAfterCreateRequest.Publish(ctx, "timeline", response)

	return openapi.CreateTimeline200JSONResponse(response), nil
}

func (s *Service) DeleteTimeline(ctx context.Context, request openapi.DeleteTimelineRequestObject) (openapi.DeleteTimelineResponseObject, error) {
	s.hooks.OnRecordBeforeDeleteRequest.Publish(ctx, "timeline", request.Id)

	err := s.queries.DeleteTimeline(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterDeleteRequest.Publish(ctx, "timeline", request.Id)

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

	s.hooks.OnRecordViewRequest.Publish(ctx, "timeline", response)

	return openapi.GetTimeline200JSONResponse(response), nil
}

func (s *Service) UpdateTimeline(ctx context.Context, request openapi.UpdateTimelineRequestObject) (openapi.UpdateTimelineResponseObject, error) {
	s.hooks.OnRecordBeforeUpdateRequest.Publish(ctx, "timeline", request.Body)

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

	s.hooks.OnRecordAfterUpdateRequest.Publish(ctx, "timeline", response)

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

	s.hooks.OnRecordsListRequest.Publish(ctx, "types", response)

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
	s.hooks.OnRecordBeforeCreateRequest.Publish(ctx, "types", request.Body)

	t, err := s.queries.CreateType(ctx, sqlc.CreateTypeParams{
		ID:       generateID("t"),
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

	s.hooks.OnRecordAfterCreateRequest.Publish(ctx, "types", response)

	return openapi.CreateType200JSONResponse(response), nil
}

func (s *Service) DeleteType(ctx context.Context, request openapi.DeleteTypeRequestObject) (openapi.DeleteTypeResponseObject, error) {
	s.hooks.OnRecordBeforeDeleteRequest.Publish(ctx, "types", request.Id)

	err := s.queries.DeleteType(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterDeleteRequest.Publish(ctx, "types", request.Id)

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

	s.hooks.OnRecordViewRequest.Publish(ctx, "types", response)

	return openapi.GetType200JSONResponse(response), nil
}

func (s *Service) UpdateType(ctx context.Context, request openapi.UpdateTypeRequestObject) (openapi.UpdateTypeResponseObject, error) {
	s.hooks.OnRecordBeforeUpdateRequest.Publish(ctx, "types", request.Body)

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

	s.hooks.OnRecordAfterUpdateRequest.Publish(ctx, "types", response)

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
			EmailVisibility:        user.Emailvisibility,
			Id:                     user.ID,
			LastLoginAlertSentAt:   user.Lastloginalertsentat,
			LastResetSentAt:        user.Lastresetsentat,
			LastVerificationSentAt: user.Lastverificationsentat,
			Name:                   user.Name,
			Updated:                user.Updated,
			Username:               user.Username,
			Verified:               user.Verified,
		})
	}

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
	s.hooks.OnRecordBeforeCreateRequest.Publish(ctx, "users", request.Body)

	if request.Body.Password != request.Body.PasswordConfirm {
		return nil, errors.New("passwords do not match")
	}

	passwordHash, tokenKey, err := auth.HashPassword(request.Body.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.queries.CreateUser(ctx, sqlc.CreateUserParams{
		ID:              generateID("u"),
		Name:            request.Body.Name,
		Email:           request.Body.Email,
		EmailVisibility: request.Body.EmailVisibility,
		Username:        request.Body.Username,
		PasswordHash:    passwordHash,
		TokenKey:        tokenKey,
		Avatar:          request.Body.Avatar,
		Verified:        request.Body.Verified,
	})
	if err != nil {
		return nil, err
	}

	response := openapi.User{
		Avatar:                 user.Avatar,
		Created:                user.Created,
		Email:                  user.Email,
		EmailVisibility:        user.Emailvisibility,
		Id:                     user.ID,
		LastLoginAlertSentAt:   user.Lastloginalertsentat,
		LastResetSentAt:        user.Lastresetsentat,
		LastVerificationSentAt: user.Lastverificationsentat,
		Name:                   user.Name,
		Updated:                user.Updated,
		Username:               user.Username,
		Verified:               user.Verified,
	}

	s.hooks.OnRecordAfterCreateRequest.Publish(ctx, "users", response)

	return openapi.CreateUser200JSONResponse(response), nil
}

func (s *Service) DeleteUser(ctx context.Context, request openapi.DeleteUserRequestObject) (openapi.DeleteUserResponseObject, error) {
	s.hooks.OnRecordBeforeDeleteRequest.Publish(ctx, "users", request.Id)

	err := s.queries.DeleteUser(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterDeleteRequest.Publish(ctx, "users", request.Id)

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
		EmailVisibility:        user.Emailvisibility,
		Id:                     user.ID,
		LastLoginAlertSentAt:   user.Lastloginalertsentat,
		LastResetSentAt:        user.Lastresetsentat,
		LastVerificationSentAt: user.Lastverificationsentat,
		Name:                   user.Name,
		Updated:                user.Updated,
		Username:               user.Username,
		Verified:               user.Verified,
	}

	s.hooks.OnRecordViewRequest.Publish(ctx, "users", response)

	return openapi.GetUser200JSONResponse(response), nil
}

func (s *Service) UpdateUser(ctx context.Context, request openapi.UpdateUserRequestObject) (openapi.UpdateUserResponseObject, error) {
	s.hooks.OnRecordBeforeUpdateRequest.Publish(ctx, "users", request.Body)

	var passwordHash, tokenHash sql.NullString

	switch {
	case request.Body.Password == nil && request.Body.PasswordConfirm == nil:
	case request.Body.Password != nil && request.Body.PasswordConfirm != nil:
		if request.Body.Password != request.Body.PasswordConfirm {
			return nil, errors.New("passwords do not match")
		}

		passwordHashS, tokenHashS, err := auth.HashPassword(*request.Body.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}

		passwordHash = sql.NullString{String: passwordHashS, Valid: true}
		tokenHash = sql.NullString{String: tokenHashS, Valid: true}
	default:
		return nil, errors.New("password and password confirm must be provided together")
	}

	user, err := s.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		Name:            toNullString(request.Body.Name),
		Email:           toNullString(request.Body.Email),
		EmailVisibility: toNullBool(request.Body.EmailVisibility),
		Username:        toNullString(request.Body.Username),
		PasswordHash:    passwordHash,
		TokenKey:        tokenHash,
		Avatar:          toNullString(request.Body.Avatar),
		Verified:        toNullBool(request.Body.Verified),
		ID:              request.Id,
	})
	if err != nil {
		return nil, err
	}

	response := openapi.User{
		Avatar:                 user.Avatar,
		Created:                user.Created,
		Email:                  user.Email,
		EmailVisibility:        user.Emailvisibility,
		Id:                     user.ID,
		LastLoginAlertSentAt:   user.Lastloginalertsentat,
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

func (s *Service) ListRoles(ctx context.Context, request openapi.ListRolesRequestObject) (openapi.ListRolesResponseObject, error) {
	roles, err := s.queries.ListRoles(ctx, sqlc.ListRolesParams{
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
	})
	if err != nil {
		return nil, err
	}

	response := make([]openapi.Role, 0, len(roles))

	for _, role := range roles {
		response = append(response, openapi.Role{
			Created:     role.Created,
			Id:          role.ID,
			Name:        role.Name,
			Permissions: permission.FromJSONArray(ctx, role.Permissions),
			Updated:     role.Updated,
		})
	}

	s.hooks.OnRecordsListRequest.Publish(ctx, "roles", response)

	totalCount := 0
	if len(roles) > 0 {
		totalCount = int(roles[0].TotalCount)
	}

	return openapi.ListRoles200JSONResponse{
		Body: response,
		Headers: openapi.ListRoles200ResponseHeaders{
			XTotalCount: totalCount,
		},
	}, nil
}

func (s *Service) CreateRole(ctx context.Context, request openapi.CreateRoleRequestObject) (openapi.CreateRoleResponseObject, error) {
	s.hooks.OnRecordBeforeCreateRequest.Publish(ctx, "roles", request.Body)

	role, err := s.queries.CreateRole(ctx, sqlc.CreateRoleParams{
		ID:          generateID("r"),
		Name:        request.Body.Name,
		Permissions: permission.ToJSONArray(ctx, request.Body.Permissions),
	})
	if err != nil {
		return nil, err
	}

	response := openapi.Role{
		Created:     role.Created,
		Id:          role.ID,
		Name:        role.Name,
		Permissions: permission.FromJSONArray(ctx, role.Permissions),
		Updated:     role.Updated,
	}

	s.hooks.OnRecordAfterCreateRequest.Publish(ctx, "roles", response)

	return openapi.CreateRole200JSONResponse(response), nil
}

func (s *Service) DeleteRole(ctx context.Context, request openapi.DeleteRoleRequestObject) (openapi.DeleteRoleResponseObject, error) {
	s.hooks.OnRecordBeforeDeleteRequest.Publish(ctx, "roles", request.Id)

	err := s.queries.DeleteRole(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterDeleteRequest.Publish(ctx, "roles", request.Id)

	return openapi.DeleteRole204Response{}, nil
}

func (s *Service) GetRole(ctx context.Context, request openapi.GetRoleRequestObject) (openapi.GetRoleResponseObject, error) {
	role, err := s.queries.GetRole(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	response := openapi.Role{
		Created:     role.Created,
		Id:          role.ID,
		Name:        role.Name,
		Permissions: permission.FromJSONArray(ctx, role.Permissions),
		Updated:     role.Updated,
	}

	s.hooks.OnRecordViewRequest.Publish(ctx, "roles", response)

	return openapi.GetRole200JSONResponse(response), nil
}

func (s *Service) UpdateRole(ctx context.Context, request openapi.UpdateRoleRequestObject) (openapi.UpdateRoleResponseObject, error) {
	s.hooks.OnRecordBeforeUpdateRequest.Publish(ctx, "roles", request.Body)

	var permissions sql.NullString

	if request.Body.Permissions != nil {
		permissions = sql.NullString{String: permission.ToJSONArray(ctx, *request.Body.Permissions), Valid: true}
	}

	role, err := s.queries.UpdateRole(ctx, sqlc.UpdateRoleParams{
		Name:        toNullString(request.Body.Name),
		Permissions: permissions,
		ID:          request.Id,
	})
	if err != nil {
		return nil, err
	}

	response := openapi.Role{
		Created:     role.Created,
		Id:          role.ID,
		Name:        role.Name,
		Permissions: permission.FromJSONArray(ctx, role.Permissions),
		Updated:     role.Updated,
	}

	s.hooks.OnRecordAfterUpdateRequest.Publish(ctx, "roles", response)

	return openapi.UpdateRole200JSONResponse(response), nil
}

func (s *Service) AddRoleParent(ctx context.Context, request openapi.AddRoleParentRequestObject) (openapi.AddRoleParentResponseObject, error) {
	s.hooks.OnRecordBeforeCreateRequest.Publish(ctx, "role_parents", request.Body)

	err := s.queries.AssignParentRole(ctx, sqlc.AssignParentRoleParams{
		ChildRoleID:  request.Id,
		ParentRoleID: request.Body.RoleId,
	})
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterCreateRequest.Publish(ctx, "role_parents", request.Body)

	return openapi.AddRoleParent201Response{}, nil
}

func (s *Service) RemoveRoleParent(ctx context.Context, request openapi.RemoveRoleParentRequestObject) (openapi.RemoveRoleParentResponseObject, error) {
	s.hooks.OnRecordBeforeDeleteRequest.Publish(ctx, "role_parents", request.Id)

	err := s.queries.RemoveParentRole(ctx, sqlc.RemoveParentRoleParams{
		ChildRoleID:  request.Id,
		ParentRoleID: request.ParentRoleId,
	})
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterDeleteRequest.Publish(ctx, "role_parents", request.Id)

	return openapi.RemoveRoleParent204Response{}, nil
}

func (s *Service) AddUserRole(ctx context.Context, request openapi.AddUserRoleRequestObject) (openapi.AddUserRoleResponseObject, error) {
	s.hooks.OnRecordBeforeCreateRequest.Publish(ctx, "user_roles", request.Body)

	err := s.queries.AssignRoleToUser(ctx, sqlc.AssignRoleToUserParams{
		UserID: request.Id,
		RoleID: request.Body.RoleId,
	})
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterCreateRequest.Publish(ctx, "user_roles", request.Body)

	return openapi.AddUserRole201Response{}, nil
}

func (s *Service) RemoveUserRole(ctx context.Context, request openapi.RemoveUserRoleRequestObject) (openapi.RemoveUserRoleResponseObject, error) {
	s.hooks.OnRecordBeforeDeleteRequest.Publish(ctx, "user_roles", request.Id)

	err := s.queries.RemoveRoleFromUser(ctx, sqlc.RemoveRoleFromUserParams{
		UserID: request.Id,
		RoleID: request.RoleId,
	})
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterDeleteRequest.Publish(ctx, "user_roles", request.Id)

	return openapi.RemoveUserRole204Response{}, nil
}

func (s *Service) ListUserPermissions(ctx context.Context, request openapi.ListUserPermissionsRequestObject) (openapi.ListUserPermissionsResponseObject, error) {
	permissions, err := s.queries.ListUserPermissions(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordsListRequest.Publish(ctx, "user_permissions", permissions)

	return openapi.ListUserPermissions200JSONResponse(permissions), nil
}

func (s *Service) ListUserRoles(ctx context.Context, request openapi.ListUserRolesRequestObject) (openapi.ListUserRolesResponseObject, error) {
	roles, err := s.queries.ListUserRoles(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	response := make([]openapi.UserRole, 0, len(roles))
	for _, role := range roles {
		response = append(response, openapi.UserRole{
			Created:     role.Created,
			Id:          role.ID,
			Name:        role.Name,
			Permissions: permission.FromJSONArray(ctx, role.Permissions),
			Type:        role.RoleType,
			Updated:     role.Updated,
		})
	}

	s.hooks.OnRecordsListRequest.Publish(ctx, "user_roles", response)

	return openapi.ListUserRoles200JSONResponse(response), nil
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
	s.hooks.OnRecordBeforeCreateRequest.Publish(ctx, "webhooks", request.Body)

	webhook, err := s.queries.CreateWebhook(ctx, sqlc.CreateWebhookParams{
		ID:          generateID("w"),
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

	s.hooks.OnRecordAfterCreateRequest.Publish(ctx, "webhooks", response)

	return openapi.CreateWebhook200JSONResponse(response), nil
}

func (s *Service) DeleteWebhook(ctx context.Context, request openapi.DeleteWebhookRequestObject) (openapi.DeleteWebhookResponseObject, error) {
	s.hooks.OnRecordBeforeDeleteRequest.Publish(ctx, "webhooks", request.Id)

	err := s.queries.DeleteWebhook(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.hooks.OnRecordAfterDeleteRequest.Publish(ctx, "webhooks", request.Id)

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

	s.hooks.OnRecordViewRequest.Publish(ctx, "webhooks", response)

	return openapi.GetWebhook200JSONResponse(response), nil
}

func (s *Service) UpdateWebhook(ctx context.Context, request openapi.UpdateWebhookRequestObject) (openapi.UpdateWebhookResponseObject, error) {
	s.hooks.OnRecordBeforeUpdateRequest.Publish(ctx, "webhooks", request.Body)

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

	s.hooks.OnRecordAfterUpdateRequest.Publish(ctx, "webhooks", response)

	return openapi.UpdateWebhook200JSONResponse(response), nil
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

func generateID(prefix string) string {
	return prefix + "-" + uuid.New().String()
}
