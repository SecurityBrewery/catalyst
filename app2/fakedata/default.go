package fakedata

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/SecurityBrewery/catalyst/app2"
	"github.com/SecurityBrewery/catalyst/app2/database/sqlc"
)

func defaultData() ([]sqlc.Ticket, []sqlc.Comment, []sqlc.Timeline, []sqlc.Task, []sqlc.Link, []sqlc.Reaction) {
	var (
		ticketCreated   = time.Date(2025, 2, 1, 11, 29, 35, 0, time.UTC)
		ticketUpdated   = ticketCreated.Add(time.Minute * 5)
		commentCreated  = ticketCreated.Add(time.Minute * 10)
		commentUpdated  = commentCreated.Add(time.Minute * 5)
		timelineCreated = ticketCreated.Add(time.Minute * 15)
		timelineUpdated = timelineCreated.Add(time.Minute * 5)
		taskCreated     = ticketCreated.Add(time.Minute * 20)
		taskUpdated     = taskCreated.Add(time.Minute * 5)
		linkCreated     = ticketCreated.Add(time.Minute * 25)
		linkUpdated     = linkCreated.Add(time.Minute * 5)
		reactionCreated = time.Date(2025, 2, 1, 11, 30, 0, 0, time.UTC)
		reactionUpdated = reactionCreated.Add(time.Minute * 5)
	)

	createTicketActionData := `{"requirements":"pocketbase","script":"import sys\nimport json\nimport random\nimport os\n\nfrom pocketbase import PocketBase\n\n# Connect to the PocketBase server\nclient = PocketBase(os.environ[\"CATALYST_APP_URL\"])\nclient.auth_store.save(token=os.environ[\"CATALYST_TOKEN\"])\n\nnewtickets = client.collection(\"tickets\").get_list(1, 200, {\"filter\": 'name = \"New Ticket\"'})\nfor ticket in newtickets.items:\n\tclient.collection(\"tickets\").delete(ticket.id)\n\n# Create a new ticket\nclient.collection(\"tickets\").create({\n\t\"name\": \"New Ticket\",\n\t\"type\": \"alert\",\n\t\"open\": True,\n})"}`

	return []sqlc.Ticket{
			{
				ID:          "t_0",
				Created:     dateTime(ticketCreated),
				Updated:     dateTime(ticketUpdated),
				Name:        "phishing-123",
				Type:        "alert",
				Description: "Phishing email reported by several employees.",
				Open:        true,
				Schema:      `{"type":"object","properties":{"tlp":{"title":"TLP","type":"string"}}}`,
				State:       `{"severity":"Medium"}`,
				Owner:       "u_test",
			},
		}, []sqlc.Comment{
			{
				ID:      "c_0",
				Created: dateTime(commentCreated),
				Updated: dateTime(commentUpdated),
				Ticket:  "t_0",
				Author:  "u_test",
				Message: "This is a test comment.",
			},
		}, []sqlc.Timeline{
			{
				ID:      "tl_0",
				Created: dateTime(timelineCreated),
				Updated: dateTime(timelineUpdated),
				Ticket:  "t_0",
				Time:    dateTime(timelineCreated),
				Message: "This is a test timeline message.",
			},
		}, []sqlc.Task{
			{
				ID:      "ts_0",
				Created: dateTime(taskCreated),
				Updated: dateTime(taskUpdated),
				Ticket:  "t_0",
				Name:    "This is a test task.",
				Open:    true,
				Owner:   "u_test",
			},
		}, []sqlc.Link{
			{
				ID:      "l_0",
				Created: dateTime(linkCreated),
				Updated: dateTime(linkUpdated),
				Ticket:  "t_0",
				Url:     "https://www.example.com",
				Name:    "This is a test link.",
			},
		}, []sqlc.Reaction{
			{
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
	return updated.Format(time.RFC3339)
}

func GenerateDefaultData(ctx context.Context, queries *sqlc.Queries) error { //nolint:cyclop
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("1234567890"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to generate password hash: %w", err)
	}

	// users
	_, err = queries.CreateUser(ctx, sqlc.CreateUserParams{
		ID:           "u_test",
		Username:     "u_test",
		Name:         "Test User",
		Email:        "user@catalyst-soar.com",
		Verified:     true,
		PasswordHash: string(passwordHash),
	})
	if err != nil {
		return fmt.Errorf("failed to create test user: %w", err)
	}

	tickets, comments, timelines, tasks, links, reactions := defaultData()

	for _, ticket := range tickets {
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
		})
		if err != nil {
			return fmt.Errorf("failed to create ticket: %w", err)
		}
	}

	for _, comment := range comments {
		_, err := queries.CreateComment(ctx, sqlc.CreateCommentParams{
			ID:      comment.ID,
			Ticket:  comment.Ticket,
			Author:  comment.Author,
			Message: comment.Message,
		})
		if err != nil {
			return fmt.Errorf("failed to create comment: %w", err)
		}
	}

	for _, timeline := range timelines {
		_, err := queries.CreateTimeline(ctx, sqlc.CreateTimelineParams{
			ID:      timeline.ID,
			Ticket:  timeline.Ticket,
			Time:    timeline.Time,
			Message: timeline.Message,
		})
		if err != nil {
			return fmt.Errorf("failed to create timeline: %w", err)
		}
	}

	for _, task := range tasks {
		_, err := queries.CreateTask(ctx, sqlc.CreateTaskParams{
			ID:     task.ID,
			Ticket: task.Ticket,
			Name:   task.Name,
			Open:   task.Open,
			Owner:  task.Owner,
		})
		if err != nil {
			return fmt.Errorf("failed to create task: %w", err)
		}
	}

	for _, link := range links {
		_, err := queries.CreateLink(ctx, sqlc.CreateLinkParams{
			ID:     link.ID,
			Ticket: link.Ticket,
			Url:    link.Url,
			Name:   link.Name,
		})
		if err != nil {
			return fmt.Errorf("failed to create link: %w", err)
		}
	}

	for _, reaction := range reactions {
		_, err := queries.CreateReaction(ctx, sqlc.CreateReactionParams{
			ID:          reaction.ID,
			Name:        reaction.Name,
			Trigger:     reaction.Trigger,
			Triggerdata: reaction.Triggerdata,
			Action:      reaction.Action,
			Actiondata:  reaction.Actiondata,
		})
		if err != nil {
			return fmt.Errorf("failed to create reaction: %w", err)
		}
	}

	return nil
}

func ValidateDefaultData(ctx context.Context, t *testing.T, app *app2.App2) error { //nolint:cyclop
	t.Helper()

	// users
	userRecord, err := app.Queries.UserByUserName(ctx, "u_test")
	if err != nil {
		return fmt.Errorf("failed to find user record: %w", err)
	}

	if userRecord.Username != "u_test" {
		return fmt.Errorf(`username does not match: got %q, want "u_test"`, userRecord.Username)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userRecord.Passwordhash), []byte("1234567890")); err != nil {
		return fmt.Errorf("password hash does not match: %w", err)
	}

	if userRecord.Name != "Test User" {
		return fmt.Errorf(`name does not match: got %q, want "Test User"`, userRecord.Name)
	}

	if userRecord.Email != "user@catalyst-soar.com" {
		return fmt.Errorf(`email does not match: got %q, want "user@catalyst-soar.com"`, userRecord.Email)
	}

	if !userRecord.Verified {
		return errors.New("user is not verified")
	}

	tickets, comments, timelines, tasks, links, reactions := defaultData()

	for _, ticket := range tickets {
		ticket, err := app.Queries.Ticket(ctx, ticket.ID)
		if err != nil {
			return fmt.Errorf("failed to find ticket: %w", err)
		}

		assert.Equal(t, "phishing-123", ticket.Name)
		assert.Equal(t, "Phishing email reported by several employees.", ticket.Description)
		assert.True(t, ticket.Open)
		assert.Equal(t, "alert", ticket.Type)
		assert.Equal(t, "u_test", ticket.Owner)
		assert.JSONEq(t, `{"type":"object","properties":{"tlp":{"title":"TLP","type":"string"}}}`, ticket.Schema)
		assert.JSONEq(t, `{"severity":"Medium"}`, ticket.State)
	}

	for _, comment := range comments {
		comment, err := app.Queries.GetComment(ctx, comment.ID)
		if err != nil {
			return fmt.Errorf("failed to find comment: %w", err)
		}

		assert.Equal(t, "This is a test comment.", comment.Message)
	}

	for _, timeline := range timelines {
		timeline, err := app.Queries.GetTimeline(ctx, timeline.ID)
		if err != nil {
			return fmt.Errorf("failed to find timeline: %w", err)
		}

		assert.Equal(t, "This is a test timeline message.", timeline.Message)
	}

	for _, task := range tasks {
		task, err := app.Queries.GetTask(ctx, task.ID)
		if err != nil {
			return fmt.Errorf("failed to find task: %w", err)
		}

		assert.Equal(t, "This is a test task.", task.Name)
	}

	for _, link := range links {
		link, err := app.Queries.GetLink(ctx, link.ID)
		if err != nil {
			return fmt.Errorf("failed to find link: %w", err)
		}

		assert.Equal(t, "https://www.example.com", link.Url)
		assert.Equal(t, "This is a test link.", link.Name)
	}

	for _, reaction := range reactions {
		reaction, err := app.Queries.GetReaction(ctx, reaction.ID)
		if err != nil {
			return fmt.Errorf("failed to find reaction: %w", err)
		}

		assert.Equal(t, "Create New Ticket", reaction.Name)
		assert.Equal(t, "schedule", reaction.Trigger)
		assert.JSONEq(t, "{\"expression\":\"12 * * * *\"}", reaction.Triggerdata)
		assert.Equal(t, "python", reaction.Action)
	}

	return nil
}
