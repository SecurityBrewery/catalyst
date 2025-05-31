package app2

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/SecurityBrewery/catalyst/app2/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app2/openapi"
)

const (
	defaultLimit  = 100
	defaultOffset = 0
)

type Service struct {
	Queries *sqlc.Queries

	OnRecordsListRequest        *Hook
	OnRecordViewRequest         *Hook
	OnRecordBeforeCreateRequest *Hook
	OnRecordAfterCreateRequest  *Hook
	OnRecordBeforeUpdateRequest *Hook
	OnRecordAfterUpdateRequest  *Hook
	OnRecordBeforeDeleteRequest *Hook
	OnRecordAfterDeleteRequest  *Hook
}

func NewService(queries *sqlc.Queries) *Service {
	return &Service{
		Queries:                     queries,
		OnRecordsListRequest:        &Hook{},
		OnRecordViewRequest:         &Hook{},
		OnRecordBeforeCreateRequest: &Hook{},
		OnRecordAfterCreateRequest:  &Hook{},
		OnRecordBeforeUpdateRequest: &Hook{},
		OnRecordAfterUpdateRequest:  &Hook{},
		OnRecordBeforeDeleteRequest: &Hook{},
		OnRecordAfterDeleteRequest:  &Hook{},
	}
}

type Hook struct {
	subscribers []func(ctx context.Context, table string, record any)
}

func (h *Hook) Publish(ctx context.Context, table string, record any) {
	for _, subscriber := range h.subscribers {
		subscriber(ctx, table, record)
	}
}

func (h *Hook) Subscribe(fn func(ctx context.Context, table string, record any)) {
	h.subscribers = append(h.subscribers, fn)
}

func (s *Service) ListComments(ctx context.Context, request openapi.ListCommentsRequestObject) (openapi.ListCommentsResponseObject, error) {
	comments, err := s.Queries.ListComments(ctx, sqlc.ListCommentsParams{
		Ticket: toString(request.Params.Ticket, ""),
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
	})
	if err != nil {
		return nil, err
	}

	var response []openapi.ExtendedComment
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

	s.OnRecordsListRequest.Publish(ctx, "comments", response)

	return openapi.ListComments200JSONResponse{
		Body: response,
		Headers: openapi.ListComments200ResponseHeaders{
			XTotalCount: totalCount,
		},
	}, nil
}

func (s *Service) CreateComment(ctx context.Context, request openapi.CreateCommentRequestObject) (openapi.CreateCommentResponseObject, error) {
	s.OnRecordBeforeCreateRequest.Publish(ctx, "comments", request.Body)

	comment, err := s.Queries.CreateComment(ctx, sqlc.CreateCommentParams{
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

	s.OnRecordAfterCreateRequest.Publish(ctx, "comments", response)

	return openapi.CreateComment200JSONResponse(response), nil
}

func (s *Service) DeleteComment(ctx context.Context, request openapi.DeleteCommentRequestObject) (openapi.DeleteCommentResponseObject, error) {
	s.OnRecordBeforeDeleteRequest.Publish(ctx, "comments", request.Id)

	err := s.Queries.DeleteComment(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.OnRecordAfterDeleteRequest.Publish(ctx, "comments", request.Id)

	return openapi.DeleteComment204Response{}, nil
}

func (s *Service) GetComment(ctx context.Context, request openapi.GetCommentRequestObject) (openapi.GetCommentResponseObject, error) {
	comment, err := s.Queries.GetComment(ctx, request.Id)
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

	s.OnRecordViewRequest.Publish(ctx, "comments", response)

	return openapi.GetComment200JSONResponse(response), nil
}

func (s *Service) UpdateComment(ctx context.Context, request openapi.UpdateCommentRequestObject) (openapi.UpdateCommentResponseObject, error) {
	s.OnRecordBeforeUpdateRequest.Publish(ctx, "comments", request.Body)

	comment, err := s.Queries.UpdateComment(ctx, sqlc.UpdateCommentParams{
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

	s.OnRecordAfterUpdateRequest.Publish(ctx, "comments", response)

	return openapi.UpdateComment200JSONResponse(response), nil
}

func (s *Service) GetDashboardCounts(ctx context.Context, request openapi.GetDashboardCountsRequestObject) (openapi.GetDashboardCountsResponseObject, error) {
	counts, err := s.Queries.GetDashboardCounts(ctx)
	if err != nil {
		return nil, err
	}

	var response []openapi.DashboardCounts
	for _, count := range counts {
		response = append(response, openapi.DashboardCounts{
			Id:    count.ID,
			Count: int(count.Count),
		})
	}

	s.OnRecordsListRequest.Publish(ctx, "dashboard_counts", response)

	return openapi.GetDashboardCounts200JSONResponse(response), nil
}

func (s *Service) ListFeatures(ctx context.Context, request openapi.ListFeaturesRequestObject) (openapi.ListFeaturesResponseObject, error) {
	features, err := s.Queries.ListFeatures(ctx, sqlc.ListFeaturesParams{
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
	})
	if err != nil {
		return nil, err
	}

	var response []openapi.Feature
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

	s.OnRecordsListRequest.Publish(ctx, "features", response)

	return openapi.ListFeatures200JSONResponse{
		Body: response,
		Headers: openapi.ListFeatures200ResponseHeaders{
			XTotalCount: totalCount,
		},
	}, nil
}

func (s *Service) CreateFeature(ctx context.Context, request openapi.CreateFeatureRequestObject) (openapi.CreateFeatureResponseObject, error) {
	s.OnRecordBeforeCreateRequest.Publish(ctx, "features", request.Body)

	feature, err := s.Queries.CreateFeature(ctx, request.Body.Name)
	if err != nil {
		return nil, err
	}

	response := openapi.Feature{
		Id:      feature.ID,
		Name:    feature.Name,
		Created: feature.Created,
		Updated: feature.Updated,
	}

	s.OnRecordAfterCreateRequest.Publish(ctx, "features", response)

	return openapi.CreateFeature200JSONResponse(response), nil
}

func (s *Service) DeleteFeature(ctx context.Context, request openapi.DeleteFeatureRequestObject) (openapi.DeleteFeatureResponseObject, error) {
	s.OnRecordBeforeDeleteRequest.Publish(ctx, "features", request.Id)

	err := s.Queries.DeleteFeature(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.OnRecordAfterDeleteRequest.Publish(ctx, "features", request.Id)

	return openapi.DeleteFeature204Response{}, nil
}

func (s *Service) GetFeature(ctx context.Context, request openapi.GetFeatureRequestObject) (openapi.GetFeatureResponseObject, error) {
	feature, err := s.Queries.GetFeature(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	response := openapi.Feature{
		Id:      feature.ID,
		Name:    feature.Name,
		Created: feature.Created,
		Updated: feature.Updated,
	}

	s.OnRecordViewRequest.Publish(ctx, "features", response)

	return openapi.GetFeature200JSONResponse(response), nil
}

func (s *Service) ListFiles(ctx context.Context, request openapi.ListFilesRequestObject) (openapi.ListFilesResponseObject, error) {
	files, err := s.Queries.ListFiles(ctx, sqlc.ListFilesParams{
		Ticket: toString(request.Params.Ticket, ""),
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
	})
	if err != nil {
		return nil, err
	}

	var response []openapi.File
	for _, file := range files {
		response = append(response, openapi.File{
			Blob:    file.Blob,
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
	s.OnRecordBeforeCreateRequest.Publish(ctx, "files", request.Body)

	file, err := s.Queries.CreateFile(ctx, sqlc.CreateFileParams{
		ID:     generateID("f"),
		Name:   request.Body.Name,
		Blob:   request.Body.Blob,
		Size:   int64(request.Body.Size),
		Ticket: request.Body.Ticket,
	})
	if err != nil {
		return nil, err
	}

	s.OnRecordAfterCreateRequest.Publish(ctx, "files", file)

	return openapi.CreateFile200JSONResponse(openapi.File{
		Blob:    file.Blob,
		Created: file.Created,
		Id:      file.ID,
		Name:    file.Name,
		Size:    file.Size,
		Ticket:  file.Ticket,
		Updated: file.Updated,
	}), nil
}

func (s *Service) DeleteFile(ctx context.Context, request openapi.DeleteFileRequestObject) (openapi.DeleteFileResponseObject, error) {
	s.OnRecordBeforeDeleteRequest.Publish(ctx, "files", request.Id)

	err := s.Queries.DeleteFile(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.OnRecordAfterDeleteRequest.Publish(ctx, "files", request.Id)

	return openapi.DeleteFile204Response{}, nil
}

func (s *Service) GetFile(ctx context.Context, request openapi.GetFileRequestObject) (openapi.GetFileResponseObject, error) {
	file, err := s.Queries.GetFile(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	response := openapi.File{
		Blob:    file.Blob,
		Created: file.Created,
		Id:      file.ID,
		Name:    file.Name,
		Size:    file.Size,
		Ticket:  file.Ticket,
		Updated: file.Updated,
	}

	s.OnRecordViewRequest.Publish(ctx, "files", response)

	return openapi.GetFile200JSONResponse(response), nil
}

func (s *Service) UpdateFile(ctx context.Context, request openapi.UpdateFileRequestObject) (openapi.UpdateFileResponseObject, error) {
	s.OnRecordBeforeUpdateRequest.Publish(ctx, "files", request.Body)

	file, err := s.Queries.UpdateFile(ctx, sqlc.UpdateFileParams{
		ID:   request.Id,
		Name: toNullString(request.Body.Name),
		Blob: toNullString(request.Body.Blob),
		Size: toNullInt64(request.Body.Size),
	})
	if err != nil {
		return nil, err
	}

	response := openapi.File{
		Blob:    file.Blob,
		Created: file.Created,
		Id:      file.ID,
		Name:    file.Name,
		Size:    file.Size,
		Ticket:  file.Ticket,
		Updated: file.Updated,
	}

	s.OnRecordAfterUpdateRequest.Publish(ctx, "files", response)

	return openapi.UpdateFile200JSONResponse(response), nil
}

func (s *Service) ListLinks(ctx context.Context, request openapi.ListLinksRequestObject) (openapi.ListLinksResponseObject, error) {
	links, err := s.Queries.ListLinks(ctx, sqlc.ListLinksParams{
		Ticket: toString(request.Params.Ticket, ""),
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
	})
	if err != nil {
		return nil, err
	}

	var response []openapi.Link
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

func (s *Service) CreateLink(ctx context.Context, request openapi.CreateLinkRequestObject) (openapi.CreateLinkResponseObject, error) {
	s.OnRecordBeforeCreateRequest.Publish(ctx, "links", request.Body)

	link, err := s.Queries.CreateLink(ctx, sqlc.CreateLinkParams{
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

	s.OnRecordAfterCreateRequest.Publish(ctx, "links", response)

	return openapi.CreateLink200JSONResponse(response), nil
}

func (s *Service) DeleteLink(ctx context.Context, request openapi.DeleteLinkRequestObject) (openapi.DeleteLinkResponseObject, error) {
	s.OnRecordBeforeDeleteRequest.Publish(ctx, "links", request.Id)

	err := s.Queries.DeleteLink(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.OnRecordAfterDeleteRequest.Publish(ctx, "links", request.Id)

	return openapi.DeleteLink204Response{}, nil
}

func (s *Service) GetLink(ctx context.Context, request openapi.GetLinkRequestObject) (openapi.GetLinkResponseObject, error) {
	link, err := s.Queries.GetLink(ctx, request.Id)
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

	s.OnRecordViewRequest.Publish(ctx, "links", response)

	return openapi.GetLink200JSONResponse(response), nil
}

func (s *Service) UpdateLink(ctx context.Context, request openapi.UpdateLinkRequestObject) (openapi.UpdateLinkResponseObject, error) {
	s.OnRecordBeforeUpdateRequest.Publish(ctx, "links", request.Body)

	link, err := s.Queries.UpdateLink(ctx, sqlc.UpdateLinkParams{
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

	s.OnRecordAfterUpdateRequest.Publish(ctx, "links", response)

	return openapi.UpdateLink200JSONResponse(response), nil
}

func (s *Service) ListReactions(ctx context.Context, request openapi.ListReactionsRequestObject) (openapi.ListReactionsResponseObject, error) {
	reactions, err := s.Queries.ListReactions(ctx, sqlc.ListReactionsParams{
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
	})
	if err != nil {
		return nil, err
	}

	var response []openapi.Reaction
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

	s.OnRecordsListRequest.Publish(ctx, "reactions", response)

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
	s.OnRecordBeforeCreateRequest.Publish(ctx, "reactions", request.Body)

	reaction, err := s.Queries.CreateReaction(ctx, sqlc.CreateReactionParams{
		ID:      generateID("r"),
		Name:    request.Body.Name,
		Action:  request.Body.Action,
		Trigger: request.Body.Trigger,
	})
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

	s.OnRecordAfterCreateRequest.Publish(ctx, "reactions", response)

	return openapi.CreateReaction200JSONResponse(response), nil
}

func (s *Service) DeleteReaction(ctx context.Context, request openapi.DeleteReactionRequestObject) (openapi.DeleteReactionResponseObject, error) {
	s.OnRecordBeforeDeleteRequest.Publish(ctx, "reactions", request.Id)

	err := s.Queries.DeleteReaction(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.OnRecordAfterDeleteRequest.Publish(ctx, "reactions", request.Id)

	return openapi.DeleteReaction204Response{}, nil
}

func (s *Service) GetReaction(ctx context.Context, request openapi.GetReactionRequestObject) (openapi.GetReactionResponseObject, error) {
	reaction, err := s.Queries.GetReaction(ctx, request.Id)
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

	s.OnRecordViewRequest.Publish(ctx, "reactions", response)

	return openapi.GetReaction200JSONResponse(response), nil
}

func (s *Service) UpdateReaction(ctx context.Context, request openapi.UpdateReactionRequestObject) (openapi.UpdateReactionResponseObject, error) {
	s.OnRecordBeforeUpdateRequest.Publish(ctx, "reactions", request.Body)

	reaction, err := s.Queries.UpdateReaction(ctx, sqlc.UpdateReactionParams{
		ID:      request.Id,
		Name:    toNullString(request.Body.Name),
		Action:  toNullString(request.Body.Action),
		Trigger: toNullString(request.Body.Trigger),
	})
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

	s.OnRecordAfterUpdateRequest.Publish(ctx, "reactions", response)

	return openapi.UpdateReaction200JSONResponse(response), nil
}

func (s *Service) GetSidebar(ctx context.Context, request openapi.GetSidebarRequestObject) (openapi.GetSidebarResponseObject, error) {
	sidebar, err := s.Queries.GetSidebar(ctx)
	if err != nil {
		return nil, err
	}

	var response []openapi.Sidebar
	for _, s := range sidebar {
		response = append(response, openapi.Sidebar{
			Id:       s.ID,
			Singular: s.Singular,
			Plural:   s.Plural,
			Icon:     s.Icon,
			Count:    int(s.Count),
		})
	}

	s.OnRecordsListRequest.Publish(ctx, "sidebar", response)

	return openapi.GetSidebar200JSONResponse(response), nil
}

func (s *Service) ListTasks(ctx context.Context, request openapi.ListTasksRequestObject) (openapi.ListTasksResponseObject, error) {
	tasks, err := s.Queries.ListTasks(ctx, sqlc.ListTasksParams{
		Ticket: toString(request.Params.Ticket, ""),
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
	})
	if err != nil {
		return nil, err
	}

	var response []openapi.ExtendedTask
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

	s.OnRecordsListRequest.Publish(ctx, "tasks", response)

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
	s.OnRecordBeforeCreateRequest.Publish(ctx, "tasks", request.Body)

	task, err := s.Queries.CreateTask(ctx, sqlc.CreateTaskParams{
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

	s.OnRecordAfterCreateRequest.Publish(ctx, "tasks", response)

	return openapi.CreateTask200JSONResponse(response), nil
}

func (s *Service) DeleteTask(ctx context.Context, request openapi.DeleteTaskRequestObject) (openapi.DeleteTaskResponseObject, error) {
	s.OnRecordBeforeDeleteRequest.Publish(ctx, "tasks", request.Id)

	err := s.Queries.DeleteTask(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.OnRecordAfterDeleteRequest.Publish(ctx, "tasks", request.Id)

	return openapi.DeleteTask204Response{}, nil
}

func (s *Service) GetTask(ctx context.Context, request openapi.GetTaskRequestObject) (openapi.GetTaskResponseObject, error) {
	task, err := s.Queries.GetTask(ctx, request.Id)
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

	s.OnRecordViewRequest.Publish(ctx, "tasks", response)

	return openapi.GetTask200JSONResponse(response), nil
}

func (s *Service) UpdateTask(ctx context.Context, request openapi.UpdateTaskRequestObject) (openapi.UpdateTaskResponseObject, error) {
	s.OnRecordBeforeUpdateRequest.Publish(ctx, "tasks", request.Body)

	task, err := s.Queries.UpdateTask(ctx, sqlc.UpdateTaskParams{
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

	s.OnRecordAfterUpdateRequest.Publish(ctx, "tasks", response)

	return openapi.UpdateTask200JSONResponse(response), nil
}

func (s *Service) SearchTickets(ctx context.Context, request openapi.SearchTicketsRequestObject) (openapi.SearchTicketsResponseObject, error) {
	tickets, err := s.Queries.SearchTickets(ctx, sqlc.SearchTicketsParams{
		Query:  toNullString(request.Params.Query),
		Type:   toNullString(request.Params.Type),
		Open:   toNullBool(request.Params.Open),
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
	})
	if err != nil {
		return nil, err
	}

	response := []openapi.TicketSearch{}

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

	s.OnRecordsListRequest.Publish(ctx, "tickets", response)

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
	tickets, err := s.Queries.ListTickets(ctx, sqlc.ListTicketsParams{
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
	})
	if err != nil {
		return nil, err
	}

	var response []openapi.ExtendedTicket
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

	s.OnRecordsListRequest.Publish(ctx, "tickets", response)

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
	s.OnRecordBeforeCreateRequest.Publish(ctx, "tickets", request.Body)

	ticket, err := s.Queries.CreateTicket(ctx, sqlc.CreateTicketParams{
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

	s.OnRecordAfterCreateRequest.Publish(ctx, "tickets", response)

	return openapi.CreateTicket200JSONResponse(response), nil
}

func (s *Service) DeleteTicket(ctx context.Context, request openapi.DeleteTicketRequestObject) (openapi.DeleteTicketResponseObject, error) {
	s.OnRecordBeforeDeleteRequest.Publish(ctx, "tickets", request.Id)

	err := s.Queries.DeleteTicket(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.OnRecordAfterDeleteRequest.Publish(ctx, "tickets", request.Id)

	return openapi.DeleteTicket204Response{}, nil
}

func (s *Service) GetTicket(ctx context.Context, request openapi.GetTicketRequestObject) (openapi.GetTicketResponseObject, error) {
	ticket, err := s.Queries.Ticket(ctx, request.Id)
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

	s.OnRecordViewRequest.Publish(ctx, "tickets", response)

	return openapi.GetTicket200JSONResponse(response), nil
}

func (s *Service) UpdateTicket(ctx context.Context, request openapi.UpdateTicketRequestObject) (openapi.UpdateTicketResponseObject, error) {
	s.OnRecordBeforeUpdateRequest.Publish(ctx, "tickets", request.Body)

	ticket, err := s.Queries.UpdateTicket(ctx, sqlc.UpdateTicketParams{
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

	s.OnRecordAfterUpdateRequest.Publish(ctx, "tickets", response)

	return openapi.UpdateTicket200JSONResponse(response), nil
}

func (s *Service) ListTimeline(ctx context.Context, request openapi.ListTimelineRequestObject) (openapi.ListTimelineResponseObject, error) {
	timeline, err := s.Queries.ListTimeline(ctx, sqlc.ListTimelineParams{
		Ticket: toString(request.Params.Ticket, ""),
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
	})
	if err != nil {
		return nil, err
	}

	var response []openapi.TimelineEntry
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

	s.OnRecordsListRequest.Publish(ctx, "timeline", response)

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
	s.OnRecordBeforeCreateRequest.Publish(ctx, "timeline", request.Body)

	timeline, err := s.Queries.CreateTimeline(ctx, sqlc.CreateTimelineParams{
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

	s.OnRecordAfterCreateRequest.Publish(ctx, "timeline", response)

	return openapi.CreateTimeline200JSONResponse(response), nil
}

func (s *Service) DeleteTimeline(ctx context.Context, request openapi.DeleteTimelineRequestObject) (openapi.DeleteTimelineResponseObject, error) {
	s.OnRecordBeforeDeleteRequest.Publish(ctx, "timeline", request.Id)

	err := s.Queries.DeleteTimeline(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.OnRecordAfterDeleteRequest.Publish(ctx, "timeline", request.Id)

	return openapi.DeleteTimeline204Response{}, nil
}

func (s *Service) GetTimeline(ctx context.Context, request openapi.GetTimelineRequestObject) (openapi.GetTimelineResponseObject, error) {
	timeline, err := s.Queries.GetTimeline(ctx, request.Id)
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

	s.OnRecordViewRequest.Publish(ctx, "timeline", response)

	return openapi.GetTimeline200JSONResponse(response), nil
}

func (s *Service) UpdateTimeline(ctx context.Context, request openapi.UpdateTimelineRequestObject) (openapi.UpdateTimelineResponseObject, error) {
	s.OnRecordBeforeUpdateRequest.Publish(ctx, "timeline", request.Body)

	timeline, err := s.Queries.UpdateTimeline(ctx, sqlc.UpdateTimelineParams{
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

	s.OnRecordAfterUpdateRequest.Publish(ctx, "timeline", response)

	return openapi.UpdateTimeline200JSONResponse(response), nil
}

func (s *Service) ListTypes(ctx context.Context, request openapi.ListTypesRequestObject) (openapi.ListTypesResponseObject, error) {
	types, err := s.Queries.ListTypes(ctx)
	if err != nil {
		return nil, err
	}

	var response []openapi.Type
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

	s.OnRecordsListRequest.Publish(ctx, "types", response)

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
	s.OnRecordBeforeCreateRequest.Publish(ctx, "types", request.Body)

	t, err := s.Queries.CreateType(ctx, sqlc.CreateTypeParams{
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

	s.OnRecordAfterCreateRequest.Publish(ctx, "types", response)

	return openapi.CreateType200JSONResponse(response), nil
}

func (s *Service) DeleteType(ctx context.Context, request openapi.DeleteTypeRequestObject) (openapi.DeleteTypeResponseObject, error) {
	s.OnRecordBeforeDeleteRequest.Publish(ctx, "types", request.Id)

	err := s.Queries.DeleteType(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.OnRecordAfterDeleteRequest.Publish(ctx, "types", request.Id)

	return openapi.DeleteType204Response{}, nil
}

func (s *Service) GetType(ctx context.Context, request openapi.GetTypeRequestObject) (openapi.GetTypeResponseObject, error) {
	t, err := s.Queries.GetType(ctx, request.Id)
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

	s.OnRecordViewRequest.Publish(ctx, "types", response)

	return openapi.GetType200JSONResponse(response), nil
}

func (s *Service) UpdateType(ctx context.Context, request openapi.UpdateTypeRequestObject) (openapi.UpdateTypeResponseObject, error) {
	s.OnRecordBeforeUpdateRequest.Publish(ctx, "types", request.Body)

	t, err := s.Queries.UpdateType(ctx, sqlc.UpdateTypeParams{
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

	s.OnRecordAfterUpdateRequest.Publish(ctx, "types", response)

	return openapi.UpdateType200JSONResponse(response), nil
}

func (s *Service) ListUsers(ctx context.Context, request openapi.ListUsersRequestObject) (openapi.ListUsersResponseObject, error) {
	users, err := s.Queries.ListUsers(ctx, sqlc.ListUsersParams{
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
	})
	if err != nil {
		return nil, err
	}

	var response []openapi.User
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
			PasswordHash:           user.Passwordhash,
			TokenKey:               user.Tokenkey,
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
	s.OnRecordBeforeCreateRequest.Publish(ctx, "users", request.Body)

	if request.Body.Password != request.Body.PasswordConfirm {
		return nil, errors.New("passwords do not match")
	}

	passwordHash, err := hashPassword(request.Body.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.Queries.CreateUser(ctx, sqlc.CreateUserParams{
		ID:              generateID("u"),
		Name:            request.Body.Name,
		Email:           request.Body.Email,
		EmailVisibility: request.Body.EmailVisibility,
		Username:        request.Body.Username,
		PasswordHash:    passwordHash,
		TokenKey:        "", // TODO
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
		PasswordHash:           user.Passwordhash,
		TokenKey:               user.Tokenkey,
		Updated:                user.Updated,
		Username:               user.Username,
		Verified:               user.Verified,
	}

	s.OnRecordAfterCreateRequest.Publish(ctx, "users", response)

	return openapi.CreateUser200JSONResponse(response), nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (s *Service) DeleteUser(ctx context.Context, request openapi.DeleteUserRequestObject) (openapi.DeleteUserResponseObject, error) {
	s.OnRecordBeforeDeleteRequest.Publish(ctx, "users", request.Id)

	err := s.Queries.DeleteUser(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.OnRecordAfterDeleteRequest.Publish(ctx, "users", request.Id)

	return openapi.DeleteUser204Response{}, nil
}

func (s *Service) GetUser(ctx context.Context, request openapi.GetUserRequestObject) (openapi.GetUserResponseObject, error) {
	user, err := s.Queries.GetUser(ctx, request.Id)
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
		PasswordHash:           user.Passwordhash,
		TokenKey:               user.Tokenkey,
		Updated:                user.Updated,
		Username:               user.Username,
		Verified:               user.Verified,
	}

	s.OnRecordViewRequest.Publish(ctx, "users", response)

	return openapi.GetUser200JSONResponse(response), nil
}

func (s *Service) UpdateUser(ctx context.Context, request openapi.UpdateUserRequestObject) (openapi.UpdateUserResponseObject, error) {
	s.OnRecordBeforeUpdateRequest.Publish(ctx, "users", request.Body)

	var passwordHash sql.NullString

	switch {
	case request.Body.Password == nil && request.Body.PasswordConfirm == nil:
	case request.Body.Password != nil && request.Body.PasswordConfirm != nil:
		if request.Body.Password != request.Body.PasswordConfirm {
			return nil, errors.New("passwords do not match")
		}

		passwordHashS, err := hashPassword(*request.Body.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}

		passwordHash = sql.NullString{String: passwordHashS, Valid: true}
	default:
		return nil, errors.New("password and password confirm must be provided together")
	}

	user, err := s.Queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		Name:            toNullString(request.Body.Name),
		Email:           toNullString(request.Body.Email),
		EmailVisibility: toNullBool(request.Body.EmailVisibility),
		Username:        toNullString(request.Body.Username),
		PasswordHash:    passwordHash,
		TokenKey:        sql.NullString{},
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
		PasswordHash:           user.Passwordhash,
		TokenKey:               user.Tokenkey,
		Updated:                user.Updated,
		Username:               user.Username,
		Verified:               user.Verified,
	}

	s.OnRecordAfterUpdateRequest.Publish(ctx, "users", response)

	return openapi.UpdateUser200JSONResponse(response), nil
}

func (s *Service) ListWebhooks(ctx context.Context, request openapi.ListWebhooksRequestObject) (openapi.ListWebhooksResponseObject, error) {
	webhooks, err := s.Queries.ListWebhooks(ctx, sqlc.ListWebhooksParams{
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
	})
	if err != nil {
		return nil, err
	}

	var response []openapi.Webhook
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
	s.OnRecordBeforeCreateRequest.Publish(ctx, "webhooks", request.Body)

	webhook, err := s.Queries.CreateWebhook(ctx, sqlc.CreateWebhookParams{
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

	s.OnRecordAfterCreateRequest.Publish(ctx, "webhooks", response)

	return openapi.CreateWebhook200JSONResponse(response), nil
}

func (s *Service) DeleteWebhook(ctx context.Context, request openapi.DeleteWebhookRequestObject) (openapi.DeleteWebhookResponseObject, error) {
	s.OnRecordBeforeDeleteRequest.Publish(ctx, "webhooks", request.Id)

	err := s.Queries.DeleteWebhook(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	s.OnRecordAfterDeleteRequest.Publish(ctx, "webhooks", request.Id)

	return openapi.DeleteWebhook204Response{}, nil
}

func (s *Service) GetWebhook(ctx context.Context, request openapi.GetWebhookRequestObject) (openapi.GetWebhookResponseObject, error) {
	webhook, err := s.Queries.GetWebhook(ctx, request.Id)
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

	s.OnRecordViewRequest.Publish(ctx, "webhooks", response)

	return openapi.GetWebhook200JSONResponse(response), nil
}

func (s *Service) UpdateWebhook(ctx context.Context, request openapi.UpdateWebhookRequestObject) (openapi.UpdateWebhookResponseObject, error) {
	s.OnRecordBeforeUpdateRequest.Publish(ctx, "webhooks", request.Body)

	webhook, err := s.Queries.UpdateWebhook(ctx, sqlc.UpdateWebhookParams{
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

	s.OnRecordAfterUpdateRequest.Publish(ctx, "webhooks", response)

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

func toNullInt64(value *int64) sql.NullInt64 {
	if value == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(*value), Valid: true}
}

func toNullBool(value *bool) sql.NullBool {
	if value == nil {
		return sql.NullBool{}
	}
	return sql.NullBool{Bool: *value, Valid: true}
}

func marshal(state map[string]interface{}) string {
	b, _ := json.Marshal(state)
	return string(b)
}

func marshalPointer(state *map[string]interface{}) string {
	if state == nil {
		return "{}"
	}

	b, _ := json.Marshal(*state)
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
