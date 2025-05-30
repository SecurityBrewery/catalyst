package app2

import (
	"context"
	"database/sql"

	"github.com/SecurityBrewery/catalyst/app2/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app2/openapi"
)

const (
	defaultLimit  = 10
	defaultOffset = 0
)

type Service struct {
	Queries *sqlc.Queries
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

	var response []openapi.Comment
	for _, comment := range comments {
		response = append(response, openapi.Comment{
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
	comment, err := s.Queries.CreateComment(ctx, sqlc.CreateCommentParams{
		Author:  request.Body.Author,
		Message: request.Body.Message,
		Ticket:  request.Body.Ticket,
	})
	if err != nil {
		return nil, err
	}

	return openapi.CreateComment200JSONResponse(openapi.Comment{
		Author:  comment.Author,
		Created: comment.Created,
		Id:      comment.ID,
		Message: comment.Message,
		Ticket:  comment.Ticket,
		Updated: comment.Updated,
	}), nil
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
	comment, err := s.Queries.UpdateComment(ctx, sqlc.UpdateCommentParams{
		Message: toNullString(request.Body.Message),
		ID:      request.Id,
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateComment200JSONResponse{
		Id:      comment.ID,
		Author:  comment.Author,
		Created: comment.Created,
		Updated: comment.Updated,
		Message: comment.Message,
	}, nil
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

	var response []openapi.Feature
	for _, feature := range features {
		response = append(response, openapi.Feature{
			Id:      feature.ID,
			Name:    feature.Name,
			Created: feature.Created,
			Updated: feature.Updated,
		})
	}

	return openapi.ListFeatures200JSONResponse(response), nil
}

func (s *Service) CreateFeature(ctx context.Context, request openapi.CreateFeatureRequestObject) (openapi.CreateFeatureResponseObject, error) {
	feature, err := s.Queries.CreateFeature(ctx, request.Body.Name)
	if err != nil {
		return nil, err
	}

	return openapi.CreateFeature200JSONResponse(openapi.Feature{
		Id:      feature.ID,
		Name:    feature.Name,
		Created: feature.Created,
		Updated: feature.Updated,
	}), nil
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
			Id:      file.ID,
			Name:    file.Name,
			Created: file.Created,
			Updated: file.Updated,
		})
	}

	return openapi.ListFiles200JSONResponse(response), nil
}

func (s *Service) CreateFile(ctx context.Context, request openapi.CreateFileRequestObject) (openapi.CreateFileResponseObject, error) {
	file, err := s.Queries.CreateFile(ctx, sqlc.CreateFileParams{
		Name:   request.Body.Name,
		Blob:   request.Body.Blob,
		Size:   int64(request.Body.Size),
		Ticket: request.Body.Ticket,
	})
	if err != nil {
		return nil, err
	}

	return openapi.CreateFile200JSONResponse(openapi.File{
		Id:      file.ID,
		Name:    file.Name,
		Created: file.Created,
		Updated: file.Updated,
	}), nil
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
	file, err := s.Queries.UpdateFile(ctx, sqlc.UpdateFileParams{
		ID:   request.Id,
		Name: toNullString(request.Body.Name),
		Blob: toNullString(request.Body.Blob),
		Size: toNullInt64(request.Body.Size),
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateFile200JSONResponse{
		Id:      file.ID,
		Name:    file.Name,
		Created: file.Created,
		Updated: file.Updated,
	}, nil
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

	return openapi.ListLinks200JSONResponse(response), nil
}

func (s *Service) CreateLink(ctx context.Context, request openapi.CreateLinkRequestObject) (openapi.CreateLinkResponseObject, error) {
	link, err := s.Queries.CreateLink(ctx, sqlc.CreateLinkParams{
		Name:   request.Body.Name,
		Url:    request.Body.Url,
		Ticket: request.Body.Ticket,
	})
	if err != nil {
		return nil, err
	}

	return openapi.CreateLink200JSONResponse(openapi.Link{
		Id:      link.ID,
		Name:    link.Name,
		Created: link.Created,
		Updated: link.Updated,
		Url:     link.Url,
	}), nil
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
	link, err := s.Queries.UpdateLink(ctx, sqlc.UpdateLinkParams{
		ID:   request.Id,
		Name: toNullString(request.Body.Name),
		Url:  toNullString(request.Body.Url),
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateLink200JSONResponse{
		Id:      link.ID,
		Name:    link.Name,
		Created: link.Created,
		Updated: link.Updated,
		Url:     link.Url,
	}, nil
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
	reaction, err := s.Queries.CreateReaction(ctx, sqlc.CreateReactionParams{
		Name:    request.Body.Name,
		Action:  request.Body.Action,
		Trigger: request.Body.Trigger,
	})
	if err != nil {
		return nil, err
	}

	return openapi.CreateReaction200JSONResponse(openapi.Reaction{
		Id:      reaction.ID,
		Name:    reaction.Name,
		Created: reaction.Created,
		Updated: reaction.Updated,
		Action:  reaction.Action,
	}), nil
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
	reaction, err := s.Queries.UpdateReaction(ctx, sqlc.UpdateReactionParams{
		ID:      request.Id,
		Name:    toNullString(request.Body.Name),
		Action:  toNullString(request.Body.Action),
		Trigger: toNullString(request.Body.Trigger),
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateReaction200JSONResponse{
		Id:      reaction.ID,
		Name:    reaction.Name,
		Created: reaction.Created,
		Updated: reaction.Updated,
		Action:  reaction.Action,
	}, nil
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
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
	})
	if err != nil {
		return nil, err
	}

	var response []openapi.Task
	for _, task := range tasks {
		response = append(response, openapi.Task{
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
	task, err := s.Queries.CreateTask(ctx, sqlc.CreateTaskParams{
		Name:   request.Body.Name,
		Open:   request.Body.Open,
		Owner:  request.Body.Owner,
		Ticket: request.Body.Ticket,
	})
	if err != nil {
		return nil, err
	}

	return openapi.CreateTask200JSONResponse(openapi.Task{
		Id:      task.ID,
		Name:    task.Name,
		Created: task.Created,
		Updated: task.Updated,
	}), nil
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
	task, err := s.Queries.UpdateTask(ctx, sqlc.UpdateTaskParams{
		ID:    request.Id,
		Name:  toNullString(request.Body.Name),
		Open:  toNullBool(request.Body.Open),
		Owner: toNullString(request.Body.Owner),
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateTask200JSONResponse{
		Id:      task.ID,
		Name:    task.Name,
		Created: task.Created,
		Updated: task.Updated,
		Open:    task.Open,
	}, nil
}

func (s *Service) SearchTickets(ctx context.Context, request openapi.SearchTicketsRequestObject) (openapi.SearchTicketsResponseObject, error) {
	tickets, err := s.Queries.SearchTickets(ctx, sqlc.SearchTicketsParams{
		Query:  sql.NullString{String: request.Params.Query, Valid: true},
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
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
		Offset: toInt64(request.Params.Offset, defaultOffset),
		Limit:  toInt64(request.Params.Limit, defaultLimit),
	})
	if err != nil {
		return nil, err
	}

	var response []openapi.Ticket
	for _, ticket := range tickets {
		response = append(response, openapi.Ticket{
			Id:          ticket.ID,
			Name:        ticket.Name,
			Description: ticket.Description,
			Owner:       ticket.Owner,
		})
	}

	return openapi.ListTickets200JSONResponse(response), nil
}

func (s *Service) CreateTicket(ctx context.Context, request openapi.CreateTicketRequestObject) (openapi.CreateTicketResponseObject, error) {
	ticket, err := s.Queries.CreateTicket(ctx, sqlc.CreateTicketParams{
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

	return openapi.CreateTicket200JSONResponse(openapi.Ticket{
		Id:          ticket.ID,
		Name:        ticket.Name,
		Description: ticket.Description,
		Owner:       ticket.Owner,
		Created:     ticket.Created,
		Updated:     ticket.Updated,
		Open:        ticket.Open,
		Resolution:  ticket.Resolution,
		Type:        ticket.Type,
		// State:       ticket.State, // TODO
	}), nil
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
		// State:       ticket.State, // TODO
	}, nil
}

func (s *Service) UpdateTicket(ctx context.Context, request openapi.UpdateTicketRequestObject) (openapi.UpdateTicketResponseObject, error) {
	ticket, err := s.Queries.UpdateTicket(ctx, sqlc.UpdateTicketParams{
		ID:          request.Id,
		Name:        toNullString(request.Body.Name),
		Description: toNullString(request.Body.Description),
		Owner:       toNullString(request.Body.Owner),
		Open:        toNullBool(request.Body.Open),
		Resolution:  toNullString(request.Body.Resolution),
		Type:        toNullString(request.Body.Type),
		// State:       toNullString(request.Body.State), // TODO
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateTicket200JSONResponse{
		Id:          ticket.ID,
		Name:        ticket.Name,
		Description: ticket.Description,
		Owner:       ticket.Owner,
	}, nil
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

	return openapi.ListTimeline200JSONResponse(response), nil
}

func (s *Service) CreateTimeline(ctx context.Context, request openapi.CreateTimelineRequestObject) (openapi.CreateTimelineResponseObject, error) {
	timeline, err := s.Queries.CreateTimeline(ctx, sqlc.CreateTimelineParams{
		Message: request.Body.Message,
		Time:    request.Body.Time,
		Ticket:  request.Body.Ticket,
	})
	if err != nil {
		return nil, err
	}

	return openapi.CreateTimeline200JSONResponse(openapi.TimelineEntry{
		Id:      timeline.ID,
		Message: timeline.Message,
		Created: timeline.Created,
		Updated: timeline.Updated,
		Time:    timeline.Time,
	}), nil
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
	timeline, err := s.Queries.UpdateTimeline(ctx, sqlc.UpdateTimelineParams{
		ID:      request.Id,
		Message: toNullString(request.Body.Message),
		Time:    toNullString(request.Body.Time),
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateTimeline200JSONResponse{
		Id:      timeline.ID,
		Message: timeline.Message,
		Created: timeline.Created,
		Updated: timeline.Updated,
		Time:    timeline.Time,
	}, nil
}

func (s *Service) ListTypes(ctx context.Context, request openapi.ListTypesRequestObject) (openapi.ListTypesResponseObject, error) {
	types, err := s.Queries.ListTypes(ctx)
	if err != nil {
		return nil, err
	}

	var response []openapi.Type
	for _, t := range types {
		response = append(response, openapi.Type{
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
	t, err := s.Queries.CreateType(ctx, sqlc.CreateTypeParams{
		Icon:     request.Body.Icon,
		Plural:   request.Body.Plural,
		Singular: request.Body.Singular,
		Schema:   request.Body.Schema,
	})
	if err != nil {
		return nil, err
	}

	return openapi.CreateType200JSONResponse(openapi.Type{
		Id:       t.ID,
		Created:  t.Created,
		Updated:  t.Updated,
		Icon:     t.Icon,
		Plural:   t.Plural,
		Singular: t.Singular,
	}), nil
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
	t, err := s.Queries.UpdateType(ctx, sqlc.UpdateTypeParams{
		ID:       request.Id,
		Icon:     toNullString(request.Body.Icon),
		Plural:   toNullString(request.Body.Plural),
		Singular: toNullString(request.Body.Singular),
		Schema:   request.Body.Schema,
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateType200JSONResponse{
		Id:       t.ID,
		Created:  t.Created,
		Updated:  t.Updated,
		Icon:     t.Icon,
		Plural:   t.Plural,
		Singular: t.Singular,
	}, nil
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
			Id: user.ID,
		})
	}

	return openapi.ListUsers200JSONResponse(response), nil
}

func (s *Service) CreateUser(ctx context.Context, request openapi.CreateUserRequestObject) (openapi.CreateUserResponseObject, error) {
	user, err := s.Queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:    request.Body.Email,
		Username: request.Body.Username,
		// TokenKey:     request.Body.TokenKey, // TODO
		Name: request.Body.Name,
		// PasswordHash: request.Body.PasswordHash, // TODO
	})
	if err != nil {
		return nil, err
	}

	return openapi.CreateUser200JSONResponse(openapi.User{
		Id:       user.ID,
		Created:  user.Created,
		Updated:  user.Updated,
		Name:     user.Name,
		Email:    user.Email,
		Username: user.Username,
		// TokenKey:     user.TokenKey, // TODO
		// PasswordHash: user.PasswordHash, // TODO
	}), nil
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
	user, err := s.Queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:       request.Id,
		Name:     toNullString(request.Body.Name),
		Email:    toNullString(request.Body.Email),
		Username: toNullString(request.Body.Username),
		// TokenKey:     toNullString(request.Body.TokenKey), // TODO
		// PasswordHash: toNullString(request.Body.PasswordHash), // TODO
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateUser200JSONResponse{
		Id:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Username: user.Username,
		Verified: user.Verified,
	}, nil
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

	return openapi.ListWebhooks200JSONResponse(response), nil
}

func (s *Service) CreateWebhook(ctx context.Context, request openapi.CreateWebhookRequestObject) (openapi.CreateWebhookResponseObject, error) {
	webhook, err := s.Queries.CreateWebhook(ctx, sqlc.CreateWebhookParams{
		Name:        request.Body.Name,
		Destination: request.Body.Destination,
		Collection:  request.Body.Collection,
	})
	if err != nil {
		return nil, err
	}

	return openapi.CreateWebhook200JSONResponse(openapi.Webhook{
		Id:          webhook.ID,
		Name:        webhook.Name,
		Created:     webhook.Created,
		Updated:     webhook.Updated,
		Destination: webhook.Destination,
		Collection:  webhook.Collection,
	}), nil
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
	webhook, err := s.Queries.UpdateWebhook(ctx, sqlc.UpdateWebhookParams{
		ID:          request.Id,
		Name:        toNullString(request.Body.Name),
		Destination: toNullString(request.Body.Destination),
		Collection:  toNullString(request.Body.Collection),
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateWebhook200JSONResponse{
		Id:          webhook.ID,
		Name:        webhook.Name,
		Created:     webhook.Created,
		Updated:     webhook.Updated,
		Destination: webhook.Destination,
		Collection:  webhook.Collection,
	}, nil
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

func toNullInt64(value *int) sql.NullInt64 {
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
