package fakedata

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/SecurityBrewery/catalyst/app2/database/sqlc"
)

const (
	minimumUserCount   = 1
	minimumTicketCount = 1
)

func Generate(queries *sqlc.Queries, userCount, ticketCount int) error {
	if userCount < minimumUserCount {
		userCount = minimumUserCount
	}

	if ticketCount < minimumTicketCount {
		ticketCount = minimumTicketCount
	}

	return Records(context.Background(), queries, userCount, ticketCount)
}

func Records(ctx context.Context, queries *sqlc.Queries, userCount int, ticketCount int) error {
	types, err := queries.ListTypes(ctx)
	if err != nil {
		return fmt.Errorf("failed to list types: %w", err)
	}

	users, err := userRecords(ctx, queries, userCount)
	if err != nil {
		return fmt.Errorf("failed to create user records: %w", err)
	}

	err = ticketRecords(ctx, queries, users, types, ticketCount)
	if err != nil {
		return fmt.Errorf("failed to create ticket records: %w", err)
	}

	err = reactionRecords(ctx, queries)
	if err != nil {
		return fmt.Errorf("failed to create reaction records: %w", err)
	}

	return nil
}

func userRecords(ctx context.Context, queries *sqlc.Queries, count int) ([]sqlc.User, error) {
	records := make([]sqlc.User, 0, count)

	// create the test user
	user, err := queries.GetUser(ctx, "u_test")
	if err != nil {
		newUser, err := queries.CreateUser(ctx, sqlc.CreateUserParams{
			Username: "u_test",
			Name:     "Test User",
			Email:    "user@catalyst-soar.com",
		})
		if err != nil {
			return nil, err
		}

		records = append(records, newUser)
	} else {
		records = append(records, user)
	}

	for range count - 1 {
		newUser, err := queries.CreateUser(ctx, sqlc.CreateUserParams{
			Username: "u_" + gofakeit.RandomString([]string{"123456789"}),
			Name:     gofakeit.Name(),
			Email:    gofakeit.Username() + "@catalyst-soar.com",
		})
		if err != nil {
			return nil, err
		}

		records = append(records, newUser)
	}

	return records, nil
}

func ticketRecords(ctx context.Context, queries *sqlc.Queries, users []sqlc.User, types []sqlc.Type, count int) error {
	created := time.Now()
	number := gofakeit.Number(200*count, 300*count)

	for range count {
		number -= gofakeit.Number(100, 200)
		created = created.Add(time.Duration(-gofakeit.Number(13, 37)) * time.Hour)

		ticketType := random(types)

		newTicket, err := queries.CreateTicket(ctx, sqlc.CreateTicketParams{
			Name:        fmt.Sprintf("%s-%d", strings.ToUpper(ticketType.Singular), number),
			Type:        ticketType.ID,
			Description: fakeTicketDescription(),
			Open:        gofakeit.Bool(),
			Schema:      marshal(map[string]interface{}{"type": "object", "properties": map[string]interface{}{"tlp": map[string]interface{}{"title": "TLP", "type": "string"}}}),
			State:       marshal(map[string]interface{}{"severity": "Medium"}),
			Owner:       random(users).ID,
		})
		if err != nil {
			return fmt.Errorf("failed to create ticket: %w", err)
		}

		if err := commentRecords(ctx, queries, users, newTicket); err != nil {
			return fmt.Errorf("failed to create comments for ticket %s: %w", newTicket.ID, err)
		}

		if err := timelineRecords(ctx, queries, created, newTicket); err != nil {
			return fmt.Errorf("failed to create timeline for ticket %s: %w", newTicket.ID, err)
		}

		if err := taskRecords(ctx, queries, users, created, newTicket); err != nil {
			return fmt.Errorf("failed to create tasks for ticket %s: %w", newTicket.ID, err)
		}

		if err := linkRecords(ctx, queries, created, newTicket); err != nil {
			return fmt.Errorf("failed to create links for ticket %s: %w", newTicket.ID, err)
		}
	}

	return nil
}

func commentRecords(ctx context.Context, queries *sqlc.Queries, users []sqlc.User, record sqlc.Ticket) error {
	for range gofakeit.IntN(5) {
		_, err := queries.CreateComment(ctx, sqlc.CreateCommentParams{
			Ticket:  record.ID,
			Author:  random(users).ID,
			Message: fakeTicketComment(),
		})
		if err != nil {
			return fmt.Errorf("failed to create comment for ticket %s: %w", record.ID, err)
		}
	}

	return nil
}

func timelineRecords(ctx context.Context, queries *sqlc.Queries, created time.Time, record sqlc.Ticket) error {
	for range gofakeit.IntN(5) {
		_, err := queries.CreateTimeline(ctx, sqlc.CreateTimelineParams{
			Message: fakeTicketTimelineMessage(),
			Ticket:  record.ID,
			Time:    created.Format("2006-01-02T15:04:05Z"),
		})
		if err != nil {
			return fmt.Errorf("failed to create timeline for ticket %s: %w", record.ID, err)
		}
	}

	return nil
}

func taskRecords(ctx context.Context, queries *sqlc.Queries, users []sqlc.User, created time.Time, record sqlc.Ticket) error {
	for range gofakeit.IntN(5) {
		_, err := queries.CreateTask(ctx, sqlc.CreateTaskParams{
			Name:   fakeTicketTask(),
			Open:   gofakeit.Bool(),
			Owner:  random(users).ID,
			Ticket: record.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to create task for ticket %s: %w", record.ID, err)
		}
	}

	return nil
}

func linkRecords(ctx context.Context, queries *sqlc.Queries, created time.Time, record sqlc.Ticket) error {
	for range gofakeit.IntN(5) {
		_, err := queries.CreateLink(ctx, sqlc.CreateLinkParams{
			Ticket: record.ID,
			Url:    gofakeit.URL(),
			Name:   random([]string{"Blog", "Forum", "Wiki", "Documentation"}),
		})
		if err != nil {
			return fmt.Errorf("failed to create link for ticket %s: %w", record.ID, err)
		}
	}

	return nil
}

const createTicketPy = `import sys
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

const alertIngestPy = `import sys
import json
import random
import os

from pocketbase import PocketBase

# Parse the event from the webhook payload
event = json.loads(sys.argv[1])
body = json.loads(event["body"])

# Connect to the PocketBase server
client = PocketBase(os.environ["CATALYST_APP_URL"])
client.auth_store.save(token=os.environ["CATALYST_TOKEN"])

# Create a new ticket
client.collection("tickets").create({
	"name": body["name"],
	"type": "alert",
	"open": True,
})`

const assignTicketsPy = `import sys
import json
import random
import os

from pocketbase import PocketBase

# Parse the ticket from the input
ticket = json.loads(sys.argv[1])

# Connect to the PocketBase server
client = PocketBase(os.environ["CATALYST_APP_URL"])
client.auth_store.save(token=os.environ["CATALYST_TOKEN"])

# Get a random user
users = client.collection("users").get_list(1, 200)
random_user = random.choice(users.items)

# Assign the ticket to the random user
client.collection("tickets").update(ticket["record"]["id"], {
	"owner": random_user.id,
})`

var (
	triggerSchedule = map[string]any{"expression": "12 * * * *"}
	triggerWebhook  = map[string]any{"token": "1234567890", "path": "webhook"}
	triggerHook     = map[string]any{"collections": []any{"tickets"}, "events": []any{"create"}}
)

func reactionRecords(ctx context.Context, queries *sqlc.Queries) error {
	_, err := queries.CreateReaction(ctx, sqlc.CreateReactionParams{
		Name:        "Create New Ticket",
		Trigger:     "schedule",
		Triggerdata: marshal(triggerSchedule),
		Action:      "python",
		Actiondata: marshal(map[string]interface{}{
			"requirements": "pocketbase",
			"script":       createTicketPy,
		}),
	})
	if err != nil {
		return fmt.Errorf("failed to create reaction for schedule trigger: %w", err)
	}

	_, err = queries.CreateReaction(ctx, sqlc.CreateReactionParams{
		Name:        "Alert Ingest Webhook",
		Trigger:     "webhook",
		Triggerdata: marshal(triggerWebhook),
		Action:      "python",
		Actiondata: marshal(map[string]interface{}{
			"requirements": "pocketbase",
			"script":       alertIngestPy,
		}),
	})
	if err != nil {
		return fmt.Errorf("failed to create reaction for webhook trigger: %w", err)
	}

	_, err = queries.CreateReaction(ctx, sqlc.CreateReactionParams{
		Name:        "Assign new Tickets",
		Trigger:     "hook",
		Triggerdata: marshal(triggerHook),
		Action:      "python",
		Actiondata: marshal(map[string]interface{}{
			"requirements": "pocketbase",
			"script":       assignTicketsPy,
		}),
	})
	if err != nil {
		return fmt.Errorf("failed to create reaction for hook trigger: %w", err)
	}

	return nil
}

func marshal(m map[string]interface{}) string {
	b, _ := json.Marshal(m)
	return string(b)
}
