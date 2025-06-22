package data

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

//go:embed scripts/upgradetest.py
var Script string

func GenerateUpgradeTestData(ctx context.Context, queries *sqlc.Queries) error { //nolint:cyclop
	if _, err := createTestUser(ctx, queries); err != nil {
		return err
	}

	for _, ticket := range CreateUpgradeTestDataTickets() {
		_, err := queries.InsertTicket(ctx, sqlc.InsertTicketParams{
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
		_, err := queries.InsertComment(ctx, sqlc.InsertCommentParams{
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
		_, err := queries.InsertTimeline(ctx, sqlc.InsertTimelineParams{
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
		_, err := queries.InsertTask(ctx, sqlc.InsertTaskParams{
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
		_, err := queries.InsertLink(ctx, sqlc.InsertLinkParams{
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
		_, err := queries.InsertReaction(ctx, sqlc.InsertReactionParams{
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

func CreateUpgradeTestDataTickets() map[string]sqlc.Ticket {
	return map[string]sqlc.Ticket{
		"t_0": {
			ID:          "t_0",
			Created:     dateTime(ticketCreated),
			Updated:     dateTime(ticketCreated.Add(time.Minute * 5)),
			Name:        "phishing-123",
			Type:        "alert",
			Description: "Phishing email reported by several employees.",
			Open:        true,
			Schema:      `{"type":"object","properties":{"tlp":{"title":"TLP","type":"string"}}}`,
			State:       `{"severity":"Medium"}`,
			Owner:       "u_test",
		},
	}
}

func CreateUpgradeTestDataComments() map[string]sqlc.Comment {
	return map[string]sqlc.Comment{
		"c_0": {
			ID:      "c_0",
			Created: dateTime(ticketCreated.Add(time.Minute * 10)),
			Updated: dateTime(ticketCreated.Add(time.Minute * 15)),
			Ticket:  "t_0",
			Author:  "u_test",
			Message: "This is a test comment.",
		},
	}
}

func CreateUpgradeTestDataTimeline() map[string]sqlc.Timeline {
	return map[string]sqlc.Timeline{
		"tl_0": {
			ID:      "tl_0",
			Created: dateTime(ticketCreated.Add(time.Minute * 15)),
			Updated: dateTime(ticketCreated.Add(time.Minute * 20)),
			Ticket:  "t_0",
			Time:    dateTime(ticketCreated.Add(time.Minute * 15)),
			Message: "This is a test timeline message.",
		},
	}
}

func CreateUpgradeTestDataTasks() map[string]sqlc.Task {
	return map[string]sqlc.Task{
		"ts_0": {
			ID:      "ts_0",
			Created: dateTime(ticketCreated.Add(time.Minute * 20)),
			Updated: dateTime(ticketCreated.Add(time.Minute * 25)),
			Ticket:  "t_0",
			Name:    "This is a test task.",
			Open:    true,
			Owner:   "u_test",
		},
	}
}

func CreateUpgradeTestDataLinks() map[string]sqlc.Link {
	return map[string]sqlc.Link{
		"l_0": {
			ID:      "l_0",
			Created: dateTime(ticketCreated.Add(time.Minute * 25)),
			Updated: dateTime(ticketCreated.Add(time.Minute * 30)),
			Ticket:  "t_0",
			Url:     "https://www.example.com",
			Name:    "This is a test link.",
		},
	}
}

func CreateUpgradeTestDataReaction() map[string]sqlc.Reaction {
	var (
		reactionCreated = time.Date(2025, 2, 1, 11, 30, 0, 0, time.UTC)
		reactionUpdated = reactionCreated.Add(time.Minute * 5)
	)

	createTicketActionData := marshal(map[string]any{
		"requirements": "pocketbase",
		"script":       Script,
	})

	return map[string]sqlc.Reaction{
		"w_0": {
			ID:          "w_0",
			Created:     dateTime(reactionCreated),
			Updated:     dateTime(reactionUpdated),
			Name:        "Create New Ticket",
			Trigger:     "schedule",
			Triggerdata: `{"expression":"12 * * * *"}`,
			Action:      "python",
			Actiondata:  createTicketActionData,
		},
	}
}

func dateTime(updated time.Time) string {
	return updated.UTC().Format(time.RFC3339)
}

func marshal(m map[string]any) string {
	b, _ := json.Marshal(m) //nolint:errchkjson

	return string(b)
}
