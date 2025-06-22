package data

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

func GenerateUpgradeTestData(ctx context.Context, queries *sqlc.Queries) error { //nolint:cyclop
	if _, err := createTestUser(ctx, queries); err != nil {
		return err
	}

	for _, ticket := range CreateUpgradeTestDataTickets() {
		_, err := queries.CreateTicket(ctx, sqlc.CreateTicketParams{
			ID:          ticket.ID,
			Name:        ticket.Name,
			Type:        ticket.Type,
			Description: ticket.Description,
			Open:        ticket.Open,
			Schema:      ticket.Schema,
			State:       ticket.State,
			Owner:       ticket.Owner,
			Resolution:  ticket.Resolution,
			Created:     ticket.Created,
			Updated:     ticket.Updated,
		})
		if err != nil {
			return fmt.Errorf("failed to create ticket: %w", err)
		}
	}

	for _, comment := range CreateUpgradeTestDataComments() {
		_, err := queries.CreateComment(ctx, sqlc.CreateCommentParams{
			ID:      comment.ID,
			Ticket:  comment.Ticket,
			Author:  comment.Author,
			Message: comment.Message,
			Created: comment.Created,
			Updated: comment.Updated,
		})
		if err != nil {
			return fmt.Errorf("failed to create comment: %w", err)
		}
	}

	for _, timeline := range CreateUpgradeTestDataTimeline() {
		_, err := queries.CreateTimeline(ctx, sqlc.CreateTimelineParams{
			ID:      timeline.ID,
			Ticket:  timeline.Ticket,
			Time:    timeline.Time,
			Message: timeline.Message,
			Created: timeline.Created,
			Updated: timeline.Updated,
		})
		if err != nil {
			return fmt.Errorf("failed to create timeline: %w", err)
		}
	}

	for _, task := range CreateUpgradeTestDataTasks() {
		_, err := queries.CreateTask(ctx, sqlc.CreateTaskParams{
			ID:      task.ID,
			Ticket:  task.Ticket,
			Name:    task.Name,
			Open:    task.Open,
			Owner:   task.Owner,
			Created: task.Created,
			Updated: task.Updated,
		})
		if err != nil {
			return fmt.Errorf("failed to create task: %w", err)
		}
	}

	for _, link := range CreateUpgradeTestDataLinks() {
		_, err := queries.CreateLink(ctx, sqlc.CreateLinkParams{
			ID:      link.ID,
			Ticket:  link.Ticket,
			Url:     link.Url,
			Name:    link.Name,
			Created: link.Created,
			Updated: link.Updated,
		})
		if err != nil {
			return fmt.Errorf("failed to create link: %w", err)
		}
	}

	for _, reaction := range CreateUpgradeTestDataReaction() {
		_, err := queries.CreateReaction(ctx, sqlc.CreateReactionParams{
			ID:          reaction.ID,
			Name:        reaction.Name,
			Trigger:     reaction.Trigger,
			Triggerdata: reaction.Triggerdata,
			Action:      reaction.Action,
			Actiondata:  reaction.Actiondata,
			Created:     reaction.Created,
			Updated:     reaction.Updated,
		})
		if err != nil {
			return fmt.Errorf("failed to create reaction: %w", err)
		}
	}

	return nil
}

func dateTime(updated time.Time) string {
	return updated.Format(time.RFC3339)
}

func marshal(m map[string]any) string {
	b, _ := json.Marshal(m) //nolint:errchkjson

	return string(b)
}
