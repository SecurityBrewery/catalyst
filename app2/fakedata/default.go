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
				Triggerdata: `{"schedule":"0 0 * * *"}`,
				Action:      "python",
				Actiondata:  createTicketActionData,
			},
		}
}

func dateTime(updated time.Time) string {
	return updated.Format(time.RFC3339)
}

func GenerateDefaultData(ctx context.Context, queries *sqlc.Queries) error {
	// users
	_, err := queries.CreateUser(ctx, sqlc.CreateUserParams{
		Username: "u_test",
		Name:     "Test User",
		Email:    "user@catalyst-soar.com",
	})
	if err != nil {
		return err
	}

	tickets, comments, timelines, tasks, links, reactions := defaultData()

	for _, ticket := range tickets {
		_, err := queries.CreateTicket(ctx, sqlc.CreateTicketParams{
			Name:        ticket.Name,
			Type:        ticket.Type,
			Description: ticket.Description,
			Open:        ticket.Open,
		})
		if err != nil {
			return err
		}
	}

	for _, comment := range comments {
		_, err := queries.CreateComment(ctx, sqlc.CreateCommentParams{
			Ticket:  comment.Ticket,
			Author:  comment.Author,
			Message: comment.Message,
		})
		if err != nil {
			return err
		}
	}

	for _, timeline := range timelines {
		_, err := queries.CreateTimeline(ctx, sqlc.CreateTimelineParams{
			Ticket:  timeline.Ticket,
			Time:    timeline.Time,
			Message: timeline.Message,
		})
		if err != nil {
			return err
		}
	}

	for _, task := range tasks {
		_, err := queries.CreateTask(ctx, sqlc.CreateTaskParams{
			Ticket: task.Ticket,
			Name:   task.Name,
			Open:   task.Open,
			Owner:  task.Owner,
		})
		if err != nil {
			return err
		}
	}

	for _, link := range links {
		_, err := queries.CreateLink(ctx, sqlc.CreateLinkParams{
			Ticket: link.Ticket,
			Url:    link.Url,
			Name:   link.Name,
		})
		if err != nil {
			return err
		}
	}

	for _, reaction := range reactions {
		_, err := queries.CreateReaction(ctx, sqlc.CreateReactionParams{
			Name:        reaction.Name,
			Trigger:     reaction.Trigger,
			Triggerdata: reaction.Triggerdata,
			Action:      reaction.Action,
			Actiondata:  reaction.Actiondata,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func ValidateDefaultData(ctx context.Context, t *testing.T, app app2.App2) error { //nolint:cyclop,gocognit
	// users
	userRecord, err := app.Queries.GetUser(ctx, "u_test")
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
			return err
		}

		assert.Equal(t, ticket.Name, "phishing-123")
		assert.Equal(t, ticket.Description, "Phishing email reported by several employees.")
		assert.Equal(t, ticket.Open, true)
		assert.Equal(t, ticket.Type, "alert")
		assert.Equal(t, ticket.Owner, "u_test")
		assert.Equal(t, ticket.Schema, `{"type":"object","properties":{"tlp":{"title":"TLP","type":"string"}}}`)
		assert.Equal(t, ticket.State, `{"severity":"Medium"}`)
	}

	for _, comment := range comments {
		comment, err := app.Queries.GetComment(ctx, comment.ID)
		if err != nil {
			return err
		}

		assert.Equal(t, comment.Message, "This is a test comment.")
	}

	for _, timeline := range timelines {
		timeline, err := app.Queries.GetTimeline(ctx, timeline.ID)
		if err != nil {
			return err
		}

		assert.Equal(t, timeline.Message, "This is a test timeline message.")
	}

	for _, task := range tasks {
		task, err := app.Queries.GetTask(ctx, task.ID)
		if err != nil {
			return err
		}

		assert.Equal(t, task.Name, "This is a test task.")
	}

	for _, link := range links {
		link, err := app.Queries.GetLink(ctx, link.ID)
		if err != nil {
			return err
		}

		assert.Equal(t, link.Url, "https://www.example.com")
		assert.Equal(t, link.Name, "This is a test link.")
	}

	for _, reaction := range reactions {
		reaction, err := app.Queries.GetReaction(ctx, reaction.ID)
		if err != nil {
			return err
		}

		assert.Equal(t, reaction.Name, "Create New Ticket")
		assert.Equal(t, reaction.Trigger, "schedule")
		assert.Equal(t, reaction.Triggerdata, `{"schedule":"0 0 * * *"}`)
		assert.Equal(t, reaction.Action, "python")
		// assert.Equal(t, reaction.Actiondata, createTicketActionData)
	}

	return nil
}
