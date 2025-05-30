package app2

import (
	"context"
	"database/sql"

	"github.com/SecurityBrewery/catalyst/app2/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app2/openapi"
)

type Service struct {
	Queries *sqlc.Queries
}

func (s *Service) ListComments(ctx context.Context, request openapi.ListCommentsRequestObject) (openapi.ListCommentsResponseObject, error) {
	comments, err := s.Queries.ListComments(ctx, sqlc.ListCommentsParams{
		Ticket: *request.Params.Ticket,
		Offset: int64(*request.Params.Offset),
		Limit:  int64(*request.Params.Limit),
	})
	if err != nil {
		return nil, err
	}

	var response []openapi.Comments
	for _, comment := range comments {
		response = append(response, openapi.Comments{
			Author:  comment.Author,
			Created: comment.Created,
			Id:      comment.ID,
			Message: comment.Message,
			Ticket:  comment.Ticket,
			Updated: comment.Updated,
		})
	}

	return openapi.ListComments200JSONResponse(response), nil
}

func (s *Service) CreateComment(ctx context.Context, request openapi.CreateCommentRequestObject) (openapi.CreateCommentResponseObject, error) {
	_, err := s.Queries.CreateComment(ctx, sqlc.CreateCommentParams{
		Author:  request.Body.Author,
		Message: request.Body.Message,
		Ticket:  request.Body.Ticket,
	})
	if err != nil {
		return nil, err
	}

	return openapi.CreateComment201Response{}, nil
}

func (s *Service) DeleteComment(ctx context.Context, request openapi.DeleteCommentRequestObject) (openapi.DeleteCommentResponseObject, error) {
	err := s.Queries.DeleteComment(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return openapi.DeleteComment204Response{}, nil
}

func (s *Service) GetComment(ctx context.Context, request openapi.GetCommentRequestObject) (openapi.GetCommentResponseObject, error) {
	comment, err := s.Queries.GetComment(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return openapi.GetComment200JSONResponse{
		Author:  comment.Author,
		Created: comment.Created,
		Id:      comment.ID,
		Message: comment.Message,
		Ticket:  comment.Ticket,
		Updated: comment.Updated,
	}, nil
}

func (s *Service) UpdateComment(ctx context.Context, request openapi.UpdateCommentRequestObject) (openapi.UpdateCommentResponseObject, error) {
	err := s.Queries.UpdateComment(ctx, sqlc.UpdateCommentParams{
		Message: request.Body.Message,
		ID:      request.Id,
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateComment200Response{}, nil
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

	return openapi.GetDashboardCounts200JSONResponse(response), nil
}

func (s *Service) ListFeatures(ctx context.Context, request openapi.ListFeaturesRequestObject) (openapi.ListFeaturesResponseObject, error) {
	features, err := s.Queries.ListFeatures(ctx)
	if err != nil {
		return nil, err
	}

	var response []openapi.Features
	for _, feature := range features {
		response = append(response, openapi.Features{
			Id:      feature.ID,
			Name:    feature.Name,
			Created: feature.Created,
			Updated: feature.Updated,
		})
	}

	return openapi.ListFeatures200JSONResponse(response), nil
}

func (s *Service) CreateFeature(ctx context.Context, request openapi.CreateFeatureRequestObject) (openapi.CreateFeatureResponseObject, error) {
	_, err := s.Queries.CreateFeature(ctx, request.Body.Name)
	if err != nil {
		return nil, err
	}

	return openapi.CreateFeature201Response{}, nil
}

func (s *Service) DeleteFeature(ctx context.Context, request openapi.DeleteFeatureRequestObject) (openapi.DeleteFeatureResponseObject, error) {
	err := s.Queries.DeleteFeature(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return openapi.DeleteFeature204Response{}, nil
}

func (s *Service) GetFeature(ctx context.Context, request openapi.GetFeatureRequestObject) (openapi.GetFeatureResponseObject, error) {
	feature, err := s.Queries.GetFeature(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return openapi.GetFeature200JSONResponse{
		Id:      feature.ID,
		Name:    feature.Name,
		Created: feature.Created,
		Updated: feature.Updated,
	}, nil
}

func (s *Service) UpdateFeature(ctx context.Context, request openapi.UpdateFeatureRequestObject) (openapi.UpdateFeatureResponseObject, error) {
	err := s.Queries.UpdateFeature(ctx, sqlc.UpdateFeatureParams{
		ID:   request.Id,
		Name: request.Body.Name,
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateFeature200Response{}, nil
}

func (s *Service) ListFiles(ctx context.Context, request openapi.ListFilesRequestObject) (openapi.ListFilesResponseObject, error) {
	files, err := s.Queries.ListFiles(ctx, sqlc.ListFilesParams{
		Ticket: *request.Params.Ticket,
		Offset: int64(*request.Params.Offset),
		Limit:  int64(*request.Params.Limit),
	})
	if err != nil {
		return nil, err
	}

	var response []openapi.Files
	for _, file := range files {
		response = append(response, openapi.Files{
			Id:      file.ID,
			Name:    file.Name,
			Created: file.Created,
			Updated: file.Updated,
		})
	}

	return openapi.ListFiles200JSONResponse(response), nil
}

func (s *Service) CreateFile(ctx context.Context, request openapi.CreateFileRequestObject) (openapi.CreateFileResponseObject, error) {
	_, err := s.Queries.CreateFile(ctx, sqlc.CreateFileParams{
		Name:   request.Body.Name,
		Blob:   request.Body.Blob,
		Size:   float64(request.Body.Size),
		Ticket: request.Body.Ticket,
	})
	if err != nil {
		return nil, err
	}

	return openapi.CreateFile201Response{}, nil
}

func (s *Service) DeleteFile(ctx context.Context, request openapi.DeleteFileRequestObject) (openapi.DeleteFileResponseObject, error) {
	err := s.Queries.DeleteFile(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return openapi.DeleteFile204Response{}, nil
}

func (s *Service) GetFile(ctx context.Context, request openapi.GetFileRequestObject) (openapi.GetFileResponseObject, error) {
	file, err := s.Queries.GetFile(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return openapi.GetFile200JSONResponse{
		Id:      file.ID,
		Name:    file.Name,
		Created: file.Created,
		Updated: file.Updated,
	}, nil
}

func (s *Service) UpdateFile(ctx context.Context, request openapi.UpdateFileRequestObject) (openapi.UpdateFileResponseObject, error) {
	err := s.Queries.UpdateFile(ctx, sqlc.UpdateFileParams{
		ID:   request.Id,
		Name: request.Body.Name,
		Blob: request.Body.Blob,
		Size: float64(request.Body.Size),
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateFile200Response{}, nil
}

func (s *Service) ListLinks(ctx context.Context, request openapi.ListLinksRequestObject) (openapi.ListLinksResponseObject, error) {
	links, err := s.Queries.ListLinks(ctx, sqlc.ListLinksParams{
		Ticket: *request.Params.Ticket,
		Offset: int64(*request.Params.Offset),
		Limit:  int64(*request.Params.Limit),
	})
	if err != nil {
		return nil, err
	}

	var response []openapi.Links
	for _, link := range links {
		response = append(response, openapi.Links{
			Id:      link.ID,
			Name:    link.Name,
			Created: link.Created,
			Updated: link.Updated,
			Url:     link.Url,
			Ticket:  link.Ticket,
		})
	}

	return openapi.ListLinks200JSONResponse(response), nil
}

func (s *Service) CreateLink(ctx context.Context, request openapi.CreateLinkRequestObject) (openapi.CreateLinkResponseObject, error) {
	_, err := s.Queries.CreateLink(ctx, sqlc.CreateLinkParams{
		Name:   request.Body.Name,
		Url:    request.Body.Url,
		Ticket: request.Body.Ticket,
	})
	if err != nil {
		return nil, err
	}

	return openapi.CreateLink201Response{}, nil
}

func (s *Service) DeleteLink(ctx context.Context, request openapi.DeleteLinkRequestObject) (openapi.DeleteLinkResponseObject, error) {
	err := s.Queries.DeleteLink(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return openapi.DeleteLink204Response{}, nil
}

func (s *Service) GetLink(ctx context.Context, request openapi.GetLinkRequestObject) (openapi.GetLinkResponseObject, error) {
	link, err := s.Queries.GetLink(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return openapi.GetLink200JSONResponse{
		Id:      link.ID,
		Name:    link.Name,
		Created: link.Created,
		Updated: link.Updated,
		Url:     link.Url,
		Ticket:  link.Ticket,
	}, nil
}

func (s *Service) UpdateLink(ctx context.Context, request openapi.UpdateLinkRequestObject) (openapi.UpdateLinkResponseObject, error) {
	err := s.Queries.UpdateLink(ctx, sqlc.UpdateLinkParams{
		ID:   request.Id,
		Name: request.Body.Name,
		Url:  request.Body.Url,
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateLink200Response{}, nil
}

func (s *Service) ListReactions(ctx context.Context, request openapi.ListReactionsRequestObject) (openapi.ListReactionsResponseObject, error) {
	reactions, err := s.Queries.ListReactions(ctx, sqlc.ListReactionsParams{
		Offset: int64(*request.Params.Offset),
		Limit:  int64(*request.Params.Limit),
	})
	if err != nil {
		return nil, err
	}

	var response []openapi.Reactions
	for _, reaction := range reactions {
		response = append(response, openapi.Reactions{
			Id:      reaction.ID,
			Name:    reaction.Name,
			Created: reaction.Created,
			Updated: reaction.Updated,
			Action:  reaction.Action,
			Trigger: reaction.Trigger,
		})
	}

	return openapi.ListReactions200JSONResponse(response), nil
}

func (s *Service) CreateReaction(ctx context.Context, request openapi.CreateReactionRequestObject) (openapi.CreateReactionResponseObject, error) {
	_, err := s.Queries.CreateReaction(ctx, sqlc.CreateReactionParams{
		Name:    request.Body.Name,
		Action:  request.Body.Action,
		Trigger: request.Body.Trigger,
	})
	if err != nil {
		return nil, err
	}

	return openapi.CreateReaction201Response{}, nil
}

func (s *Service) DeleteReaction(ctx context.Context, request openapi.DeleteReactionRequestObject) (openapi.DeleteReactionResponseObject, error) {
	err := s.Queries.DeleteReaction(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return openapi.DeleteReaction204Response{}, nil
}

func (s *Service) GetReaction(ctx context.Context, request openapi.GetReactionRequestObject) (openapi.GetReactionResponseObject, error) {
	reaction, err := s.Queries.GetReaction(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return openapi.GetReaction200JSONResponse{
		Id:      reaction.ID,
		Name:    reaction.Name,
		Created: reaction.Created,
		Updated: reaction.Updated,
		Action:  reaction.Action,
		Trigger: reaction.Trigger,
	}, nil
}

func (s *Service) UpdateReaction(ctx context.Context, request openapi.UpdateReactionRequestObject) (openapi.UpdateReactionResponseObject, error) {
	err := s.Queries.UpdateReaction(ctx, sqlc.UpdateReactionParams{
		ID:      request.Id,
		Name:    request.Body.Name,
		Action:  request.Body.Action,
		Trigger: request.Body.Trigger,
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateReaction200Response{}, nil
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

	return openapi.GetSidebar200JSONResponse(response), nil
}

func (s *Service) ListTasks(ctx context.Context, request openapi.ListTasksRequestObject) (openapi.ListTasksResponseObject, error) {
	tasks, err := s.Queries.ListTasks(ctx, sqlc.ListTasksParams{
		Ticket: *request.Params.Ticket,
		Offset: int64(*request.Params.Offset),
		Limit:  int64(*request.Params.Limit),
	})
	if err != nil {
		return nil, err
	}

	var response []openapi.Tasks
	for _, task := range tasks {
		response = append(response, openapi.Tasks{
			Id:      task.ID,
			Name:    task.Name,
			Created: task.Created,
			Updated: task.Updated,
			Open:    task.Open,
			Owner:   task.Owner,
			Ticket:  task.Ticket,
		})
	}

	return openapi.ListTasks200JSONResponse(response), nil
}

func (s *Service) CreateTask(ctx context.Context, request openapi.CreateTaskRequestObject) (openapi.CreateTaskResponseObject, error) {
	_, err := s.Queries.CreateTask(ctx, sqlc.CreateTaskParams{
		Name:   request.Body.Name,
		Open:   request.Body.Open,
		Owner:  request.Body.Owner,
		Ticket: request.Body.Ticket,
	})
	if err != nil {
		return nil, err
	}

	return openapi.CreateTask201Response{}, nil
}

func (s *Service) DeleteTask(ctx context.Context, request openapi.DeleteTaskRequestObject) (openapi.DeleteTaskResponseObject, error) {
	err := s.Queries.DeleteTask(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return openapi.DeleteTask204Response{}, nil
}

func (s *Service) GetTask(ctx context.Context, request openapi.GetTaskRequestObject) (openapi.GetTaskResponseObject, error) {
	task, err := s.Queries.GetTask(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return openapi.GetTask200JSONResponse{
		Id:      task.ID,
		Name:    task.Name,
		Created: task.Created,
		Updated: task.Updated,
		Open:    task.Open,
		Owner:   task.Owner,
		Ticket:  task.Ticket,
	}, nil
}

func (s *Service) UpdateTask(ctx context.Context, request openapi.UpdateTaskRequestObject) (openapi.UpdateTaskResponseObject, error) {
	err := s.Queries.UpdateTask(ctx, sqlc.UpdateTaskParams{
		ID:    request.Id,
		Name:  request.Body.Name,
		Open:  request.Body.Open,
		Owner: request.Body.Owner,
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateTask200Response{}, nil
}

func (s *Service) SearchTickets(ctx context.Context, request openapi.SearchTicketsRequestObject) (openapi.SearchTicketsResponseObject, error) {
	tickets, err := s.Queries.SearchTickets(ctx, sqlc.SearchTicketsParams{
		Query:  sql.NullString{String: request.Params.Query, Valid: true},
		Offset: int64(*request.Params.Offset),
		Limit:  int64(*request.Params.Limit),
	})
	if err != nil {
		return nil, err
	}

	var response []openapi.TicketSearch
	for _, ticket := range tickets {
		response = append(response, openapi.TicketSearch{
			Id:          ticket.ID,
			Name:        ticket.Name,
			Created:     ticket.Created,
			Description: ticket.Description,
			Open:        ticket.Open,
			Type:        ticket.Type,
			// State:       ticket.State,
		})
	}
	return openapi.SearchTickets200JSONResponse(response), nil
}

func (s *Service) ListTickets(ctx context.Context, request openapi.ListTicketsRequestObject) (openapi.ListTicketsResponseObject, error) {
	tickets, err := s.Queries.ListTickets(ctx, sqlc.ListTicketsParams{
		Offset: int64(*request.Params.Offset),
		Limit:  int64(*request.Params.Limit),
	})
	if err != nil {
		return nil, err
	}

	var response []openapi.Tickets
	for _, ticket := range tickets {
		response = append(response, openapi.Tickets{
			Id:          ticket.ID,
			Name:        ticket.Name,
			Description: ticket.Description,
			Owner:       ticket.Owner,
		})
	}

	return openapi.ListTickets200JSONResponse(response), nil
}

func (s *Service) CreateTicket(ctx context.Context, request openapi.CreateTicketRequestObject) (openapi.CreateTicketResponseObject, error) {
	_, err := s.Queries.CreateTicket(ctx, sqlc.CreateTicketParams{
		Name:        request.Body.Name,
		Description: request.Body.Description,
		Owner:       request.Body.Owner,
		Open:        request.Body.Open,
		Resolution:  request.Body.Resolution,
		Type:        request.Body.Type,
		State:       request.Body.State,
	})
	if err != nil {
		return nil, err
	}

	return openapi.CreateTicket201Response{}, nil
}

func (s *Service) DeleteTicket(ctx context.Context, request openapi.DeleteTicketRequestObject) (openapi.DeleteTicketResponseObject, error) {
	err := s.Queries.DeleteTicket(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return openapi.DeleteTicket204Response{}, nil
}

func (s *Service) GetTicket(ctx context.Context, request openapi.GetTicketRequestObject) (openapi.GetTicketResponseObject, error) {
	ticket, err := s.Queries.Ticket(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return openapi.GetTicket200JSONResponse{
		Id:          ticket.ID,
		Name:        ticket.Name,
		Description: ticket.Description,
		Owner:       ticket.Owner,
		Created:     ticket.Created,
		Updated:     ticket.Updated,
		Open:        ticket.Open,
		Resolution:  ticket.Resolution,
		Type:        ticket.Type,
		// State:       ticket.State,
	}, nil
}

func (s *Service) UpdateTicket(ctx context.Context, request openapi.UpdateTicketRequestObject) (openapi.UpdateTicketResponseObject, error) {
	err := s.Queries.UpdateTicket(ctx, sqlc.UpdateTicketParams{
		ID:          request.Id,
		Name:        request.Body.Name,
		Description: request.Body.Description,
		Owner:       request.Body.Owner,
		Open:        request.Body.Open,
		Resolution:  request.Body.Resolution,
		Type:        request.Body.Type,
		State:       request.Body.State,
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateTicket200Response{}, nil
}

func (s *Service) ListTimeline(ctx context.Context, request openapi.ListTimelineRequestObject) (openapi.ListTimelineResponseObject, error) {
	timeline, err := s.Queries.ListTimeline(ctx, sqlc.ListTimelineParams{
		Ticket: *request.Params.Ticket,
		Offset: int64(*request.Params.Offset),
		Limit:  int64(*request.Params.Limit),
	})
	if err != nil {
		return nil, err
	}

	var response []openapi.Timeline
	for _, timeline := range timeline {
		response = append(response, openapi.Timeline{
			Id:      timeline.ID,
			Message: timeline.Message,
			Created: timeline.Created,
			Updated: timeline.Updated,
			Time:    timeline.Time,
			Ticket:  timeline.Ticket,
		})
	}

	return openapi.ListTimeline200JSONResponse(response), nil
}

func (s *Service) CreateTimeline(ctx context.Context, request openapi.CreateTimelineRequestObject) (openapi.CreateTimelineResponseObject, error) {
	_, err := s.Queries.CreateTimeline(ctx, sqlc.CreateTimelineParams{
		Message: request.Body.Message,
		Time:    request.Body.Time,
		Ticket:  request.Body.Ticket,
	})
	if err != nil {
		return nil, err
	}

	return openapi.CreateTimeline201Response{}, nil
}

func (s *Service) DeleteTimeline(ctx context.Context, request openapi.DeleteTimelineRequestObject) (openapi.DeleteTimelineResponseObject, error) {
	err := s.Queries.DeleteTimeline(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return openapi.DeleteTimeline204Response{}, nil
}

func (s *Service) GetTimeline(ctx context.Context, request openapi.GetTimelineRequestObject) (openapi.GetTimelineResponseObject, error) {
	timeline, err := s.Queries.GetTimeline(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return openapi.GetTimeline200JSONResponse{
		Id:      timeline.ID,
		Message: timeline.Message,
		Created: timeline.Created,
		Updated: timeline.Updated,
		Time:    timeline.Time,
		Ticket:  timeline.Ticket,
	}, nil
}

func (s *Service) UpdateTimeline(ctx context.Context, request openapi.UpdateTimelineRequestObject) (openapi.UpdateTimelineResponseObject, error) {
	err := s.Queries.UpdateTimeline(ctx, sqlc.UpdateTimelineParams{
		ID:      request.Id,
		Message: request.Body.Message,
		Time:    request.Body.Time,
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateTimeline200Response{}, nil
}

func (s *Service) ListTypes(ctx context.Context, request openapi.ListTypesRequestObject) (openapi.ListTypesResponseObject, error) {
	types, err := s.Queries.ListTypes(ctx, sqlc.ListTypesParams{
		Offset: int64(*request.Params.Offset),
		Limit:  int64(*request.Params.Limit),
	})
	if err != nil {
		return nil, err
	}

	var response []openapi.Types
	for _, t := range types {
		response = append(response, openapi.Types{
			Id:       t.ID,
			Created:  t.Created,
			Updated:  t.Updated,
			Icon:     t.Icon,
			Plural:   t.Plural,
			Singular: t.Singular,
		})
	}

	return openapi.ListTypes200JSONResponse(response), nil
}

func (s *Service) CreateType(ctx context.Context, request openapi.CreateTypeRequestObject) (openapi.CreateTypeResponseObject, error) {
	_, err := s.Queries.CreateType(ctx, sqlc.CreateTypeParams{
		Icon:     request.Body.Icon,
		Plural:   request.Body.Plural,
		Singular: request.Body.Singular,
		Schema:   request.Body.Schema,
	})
	if err != nil {
		return nil, err
	}

	return openapi.CreateType201Response{}, nil
}

func (s *Service) DeleteType(ctx context.Context, request openapi.DeleteTypeRequestObject) (openapi.DeleteTypeResponseObject, error) {
	err := s.Queries.DeleteType(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return openapi.DeleteType204Response{}, nil
}

func (s *Service) GetType(ctx context.Context, request openapi.GetTypeRequestObject) (openapi.GetTypeResponseObject, error) {
	t, err := s.Queries.GetType(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return openapi.GetType200JSONResponse{
		Id:       t.ID,
		Created:  t.Created,
		Updated:  t.Updated,
		Icon:     t.Icon,
		Plural:   t.Plural,
		Singular: t.Singular,
	}, nil
}

func (s *Service) UpdateType(ctx context.Context, request openapi.UpdateTypeRequestObject) (openapi.UpdateTypeResponseObject, error) {
	err := s.Queries.UpdateType(ctx, sqlc.UpdateTypeParams{
		ID:       request.Id,
		Icon:     request.Body.Icon,
		Plural:   request.Body.Plural,
		Singular: request.Body.Singular,
		Schema:   request.Body.Schema,
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateType200Response{}, nil
}

func (s *Service) ListUsers(ctx context.Context, request openapi.ListUsersRequestObject) (openapi.ListUsersResponseObject, error) {
	users, err := s.Queries.ListUsers(ctx, sqlc.ListUsersParams{
		Offset: int64(*request.Params.Offset),
		Limit:  int64(*request.Params.Limit),
	})
	if err != nil {
		return nil, err
	}

	var response []openapi.Users
	for _, user := range users {
		response = append(response, openapi.Users{
			Id: user.ID,
		})
	}

	return openapi.ListUsers200JSONResponse(response), nil
}

func (s *Service) CreateUser(ctx context.Context, request openapi.CreateUserRequestObject) (openapi.CreateUserResponseObject, error) {
	_, err := s.Queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        request.Body.Email,
		Username:     request.Body.Username,
		TokenKey:     request.Body.TokenKey,
		Name:         request.Body.Name,
		PasswordHash: request.Body.PasswordHash,
	})
	if err != nil {
		return nil, err
	}

	return openapi.CreateUser201Response{}, nil
}

func (s *Service) DeleteUser(ctx context.Context, request openapi.DeleteUserRequestObject) (openapi.DeleteUserResponseObject, error) {
	err := s.Queries.DeleteUser(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return openapi.DeleteUser204Response{}, nil
}

func (s *Service) GetUser(ctx context.Context, request openapi.GetUserRequestObject) (openapi.GetUserResponseObject, error) {
	user, err := s.Queries.GetUser(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return openapi.GetUser200JSONResponse{
		Id:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Username: user.Username,
		Verified: user.Verified,
	}, nil
}

func (s *Service) UpdateUser(ctx context.Context, request openapi.UpdateUserRequestObject) (openapi.UpdateUserResponseObject, error) {
	err := s.Queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:           request.Id,
		Name:         request.Body.Name,
		Email:        request.Body.Email,
		Username:     request.Body.Username,
		TokenKey:     request.Body.TokenKey,
		PasswordHash: request.Body.PasswordHash,
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateUser200Response{}, nil
}

func (s *Service) ListWebhooks(ctx context.Context, request openapi.ListWebhooksRequestObject) (openapi.ListWebhooksResponseObject, error) {
	webhooks, err := s.Queries.ListWebhooks(ctx, sqlc.ListWebhooksParams{
		Offset: int64(*request.Params.Offset),
		Limit:  int64(*request.Params.Limit),
	})
	if err != nil {
		return nil, err
	}

	var response []openapi.Webhooks
	for _, webhook := range webhooks {
		response = append(response, openapi.Webhooks{
			Id:   webhook.ID,
			Name: webhook.Name,
		})
	}

	return openapi.ListWebhooks200JSONResponse(response), nil
}

func (s *Service) CreateWebhook(ctx context.Context, request openapi.CreateWebhookRequestObject) (openapi.CreateWebhookResponseObject, error) {
	_, err := s.Queries.CreateWebhook(ctx, sqlc.CreateWebhookParams{
		Name:        request.Body.Name,
		Destination: request.Body.Destination,
		Collection:  request.Body.Collection,
	})
	if err != nil {
		return nil, err
	}

	return openapi.CreateWebhook201Response{}, nil
}

func (s *Service) DeleteWebhook(ctx context.Context, request openapi.DeleteWebhookRequestObject) (openapi.DeleteWebhookResponseObject, error) {
	err := s.Queries.DeleteWebhook(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return openapi.DeleteWebhook204Response{}, nil
}

func (s *Service) GetWebhook(ctx context.Context, request openapi.GetWebhookRequestObject) (openapi.GetWebhookResponseObject, error) {
	webhook, err := s.Queries.GetWebhook(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return openapi.GetWebhook200JSONResponse{
		Id:          webhook.ID,
		Name:        webhook.Name,
		Created:     webhook.Created,
		Updated:     webhook.Updated,
		Destination: webhook.Destination,
		Collection:  webhook.Collection,
	}, nil
}

func (s *Service) UpdateWebhook(ctx context.Context, request openapi.UpdateWebhookRequestObject) (openapi.UpdateWebhookResponseObject, error) {
	err := s.Queries.UpdateWebhook(ctx, sqlc.UpdateWebhookParams{
		ID:          request.Id,
		Name:        request.Body.Name,
		Destination: request.Body.Destination,
		Collection:  request.Body.Collection,
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateWebhook200Response{}, nil
}
