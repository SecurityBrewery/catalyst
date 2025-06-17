package upgradetest

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/SecurityBrewery/catalyst/app"
	"github.com/SecurityBrewery/catalyst/app/auth/password"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

func generateUpgradeTestData(ctx context.Context, queries *sqlc.Queries) error { //nolint:cyclop
	passwordHash, tokenKey, err := password.Hash("1234567890")
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// users
	_, err = queries.CreateUser(ctx, sqlc.CreateUserParams{
		ID:           "u_test",
		Username:     "u_test",
		Name:         "Test User",
		Email:        "user@catalyst-soar.com",
		Verified:     true,
		PasswordHash: passwordHash,
		TokenKey:     tokenKey,
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

func validateUpgradeTestData(t *testing.T, app *app.App) {
	t.Helper()

	// users
	userRecord, err := app.Queries.UserByUserName(t.Context(), "u_test")
	require.NoError(t, err, "failed to find user by username")

	assert.Equal(t, "u_test", userRecord.ID, "user ID does not match expected value")
	assert.Equal(t, "u_test", userRecord.Username, "username does not match expected value")
	assert.Equal(t, "Test User", userRecord.Name, "name does not match expected value")
	assert.Equal(t, "user@catalyst-soar.com", userRecord.Email, "email does not match expected value")
	assert.True(t, userRecord.Verified, "user should be verified")
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(userRecord.Passwordhash), []byte("1234567890")), "password hash does not match expected value")

	tickets, comments, timelines, tasks, links, reactions := defaultData()

	for _, ticket := range tickets {
		ticket, err := app.Queries.Ticket(t.Context(), ticket.ID)
		require.NoError(t, err, "failed to find ticket")

		assert.Equal(t, "phishing-123", ticket.Name)
		assert.Equal(t, "Phishing email reported by several employees.", ticket.Description)
		assert.True(t, ticket.Open)
		assert.Equal(t, "alert", ticket.Type)
		assert.Equal(t, "u_test", ticket.Owner)
		assert.JSONEq(t, `{"type":"object","properties":{"tlp":{"title":"TLP","type":"string"}}}`, ticket.Schema)
		assert.JSONEq(t, `{"severity":"Medium"}`, ticket.State)
	}

	for _, comment := range comments {
		comment, err := app.Queries.GetComment(t.Context(), comment.ID)
		require.NoError(t, err, "failed to find comment")

		assert.Equal(t, "This is a test comment.", comment.Message)
	}

	for _, timeline := range timelines {
		timeline, err := app.Queries.GetTimeline(t.Context(), timeline.ID)
		require.NoError(t, err, "failed to find timeline")

		assert.Equal(t, "This is a test timeline message.", timeline.Message)
	}

	for _, task := range tasks {
		task, err := app.Queries.GetTask(t.Context(), task.ID)
		require.NoError(t, err, "failed to find task")

		assert.Equal(t, "This is a test task.", task.Name)
	}

	for _, link := range links {
		link, err := app.Queries.GetLink(t.Context(), link.ID)
		require.NoError(t, err, "failed to find link")

		assert.Equal(t, "https://www.example.com", link.Url)
		assert.Equal(t, "This is a test link.", link.Name)
	}

	for _, reaction := range reactions {
		reaction, err := app.Queries.GetReaction(t.Context(), reaction.ID)
		require.NoError(t, err, "failed to find reaction")

		assert.Equal(t, "Create New Ticket", reaction.Name)
		assert.Equal(t, "schedule", reaction.Trigger)
		assert.JSONEq(t, "{\"expression\":\"12 * * * *\"}", reaction.Triggerdata)
		assert.Equal(t, "python", reaction.Action)
		assert.JSONEq(t, marshal(map[string]interface{}{"requirements": "pocketbase", "script": script}), reaction.Actiondata)
	}
}

const script = `import sys
import json
import random
import os

from pocketbase import PocketBase

# Connect to the PocketBase server
client = PocketBase(os.environ["CATALYST_APP_URL"])
client.auth_store.save(token=os.environ["CATALYST_TOKEN"])

newtickets = client.collection("tickets").get_list(1, 200, {"filter": 'name = "New Ticket"'})
for ticket in newtickets.items:
	client.collection("tickets").delete(ticket.id)

# Create a new ticket
client.collection("tickets").create({
	"name": "New Ticket",
	"type": "alert",
	"open": True,
})`

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

	createTicketActionData := marshal(map[string]any{
		"requirements": "pocketbase",
		"script":       script,
	})

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

func marshal(m map[string]interface{}) string {
	b, _ := json.Marshal(m) //nolint:errchkjson

	return string(b)
}
