package data

import (
	"context"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

var (
	triggerSchedule = map[string]any{"expression": "12 * * * *"}
	triggerWebhook  = map[string]any{"token": "1234567890", "path": "webhook"}
	triggerHook     = map[string]any{"collections": []any{"tickets"}, "events": []any{"create"}}
)

func generateDemoReactions(ctx context.Context, queries *sqlc.Queries) error {
	_, err := queries.CreateReaction(ctx, sqlc.CreateReactionParams{
		ID:          "r-schedule",
		Name:        "Create New Ticket",
		Trigger:     "schedule",
		Triggerdata: marshal(triggerSchedule),
		Action:      "python",
		Actiondata: marshal(map[string]any{
			"requirements": "requests",
			"script":       createTicketPy,
		}),
		Created: dateTime(gofakeit.PastDate()),
		Updated: dateTime(gofakeit.PastDate()),
	})
	if err != nil {
		return fmt.Errorf("failed to create reaction for schedule trigger: %w", err)
	}

	_, err = queries.CreateReaction(ctx, sqlc.CreateReactionParams{
		ID:          "r-webhook",
		Name:        "Alert Ingest Webhook",
		Trigger:     "webhook",
		Triggerdata: marshal(triggerWebhook),
		Action:      "python",
		Actiondata: marshal(map[string]any{
			"requirements": "requests",
			"script":       alertIngestPy,
		}),
		Created: dateTime(gofakeit.PastDate()),
		Updated: dateTime(gofakeit.PastDate()),
	})
	if err != nil {
		return fmt.Errorf("failed to create reaction for webhook trigger: %w", err)
	}

	_, err = queries.CreateReaction(ctx, sqlc.CreateReactionParams{
		ID:          "r-hook",
		Name:        "Assign new Tickets",
		Trigger:     "hook",
		Triggerdata: marshal(triggerHook),
		Action:      "python",
		Actiondata: marshal(map[string]any{
			"requirements": "requests",
			"script":       assignTicketsPy,
		}),
		Created: dateTime(gofakeit.PastDate()),
		Updated: dateTime(gofakeit.PastDate()),
	})
	if err != nil {
		return fmt.Errorf("failed to create reaction for hook trigger: %w", err)
	}

	return nil
}

const createTicketPy = `import sys
import json
import random
import os

import requests

url = os.environ["CATALYST_APP_URL"]
header = {"Authorization": "Bearer " + os.environ["CATALYST_TOKEN"]}

newtickets = requests.get(url + "/api/tickets", headers=header).json()
for ticket in newtickets.items:
    requests.delete(url + "/api/tickets/" + ticket.id, headers=header)

# Create a new ticket
requests.post(url + "/api/tickets", headers=header, json={
    "name": "New Ticket",
    "type": "alert",
    "open": True,
})`

const alertIngestPy = `import sys
import json
import random
import os

import requests

# Parse the event from the webhook payload
event = json.loads(sys.argv[1])
body = json.loads(event["body"])

url = os.environ["CATALYST_APP_URL"]
header = {"Authorization": "Bearer " + os.environ["CATALYST_TOKEN"]}

# Create a new ticket
requests.post(url + "/api/tickets", headers=header, json={
	"name": body["name"],
	"type": "alert",
	"open": True,
})`

const assignTicketsPy = `import sys
import json
import random
import os

import requests

# Parse the ticket from the input
ticket = json.loads(sys.argv[1])

url = os.environ["CATALYST_APP_URL"]
header = {"Authorization": "Bearer " + os.environ["CATALYST_TOKEN"]}

# Get a random user
users = requests.get(url + "/api/users", headers=header).json()
random_user = random.choice(users.items)

# Assign the ticket to the random user
requests.patch(url + "/api/tickets/" + ticket["record"]["id"], {
	"owner": random_user.id,
})`

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

const Script = `import sys
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
