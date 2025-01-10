package fakedata

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/security"

	"github.com/SecurityBrewery/catalyst/migrations"
)

const (
	minimumUserCount   = 1
	minimumTicketCount = 1
)

func Generate(ctx context.Context, app core.App, userCount, ticketCount int) error {
	if userCount < minimumUserCount {
		userCount = minimumUserCount
	}

	if ticketCount < minimumTicketCount {
		ticketCount = minimumTicketCount
	}

	records, err := Records(app, userCount, ticketCount)
	if err != nil {
		return fmt.Errorf("failed to generate fake records: %w", err)
	}

	for _, record := range records {
		if err := app.SaveNoValidateWithContext(ctx, record); err != nil {
			return fmt.Errorf("failed to save fake record: %w", err)
		}
	}

	return nil
}

func Records(app core.App, userCount int, ticketCount int) ([]*core.Record, error) {
	types, err := app.FindAllRecords(migrations.TypeCollectionName)
	if err != nil {
		return nil, err
	}

	users := userRecords(app, userCount)
	tickets := ticketRecords(app, users, types, ticketCount)
	reactions := reactionRecords(app)

	var records []*core.Record
	records = append(records, users...)
	records = append(records, types...)
	records = append(records, tickets...)
	records = append(records, reactions...)

	return records, nil
}

func userRecords(app core.App, count int) []*core.Record {
	collection, err := app.FindCollectionByNameOrId(migrations.UserCollectionID)
	if err != nil {
		panic(err)
	}

	records := make([]*core.Record, 0, count)

	// create the test user
	if _, err := app.FindRecordById(migrations.UserCollectionID, "u_test"); err != nil {
		record := core.NewRecord(collection)
		record.Id = "u_test"
		record.Set("username", "u_test")
		record.SetPassword("1234567890")
		record.Set("name", gofakeit.Name())
		record.Set("email", "user@catalyst-soar.com")
		record.SetVerified(true)

		records = append(records, record)
	}

	for range count - 1 {
		record := core.NewRecord(collection)
		record.Id = "u_" + security.PseudorandomString(10)
		record.Set("username", "u_"+security.RandomStringWithAlphabet(5, "123456789"))
		record.SetPassword("1234567890")
		record.Set("name", gofakeit.Name())
		record.Set("email", gofakeit.Username()+"@catalyst-soar.com")
		record.SetVerified(true)

		records = append(records, record)
	}

	return records
}

func ticketRecords(app core.App, users, types []*core.Record, count int) []*core.Record {
	collection, err := app.FindCollectionByNameOrId(migrations.TicketCollectionName)
	if err != nil {
		panic(err)
	}

	records := make([]*core.Record, 0, count)

	created := time.Now()
	number := gofakeit.Number(200*count, 300*count)

	for range count {
		number -= gofakeit.Number(100, 200)
		created = created.Add(time.Duration(-gofakeit.Number(13, 37)) * time.Hour)

		record := core.NewRecord(collection)
		record.Id = "t_" + security.PseudorandomString(10)

		updated := gofakeit.DateRange(created, time.Now())

		ticketType := random(types)

		record.Set("created", created.Format("2006-01-02T15:04:05Z"))
		record.Set("updated", updated.Format("2006-01-02T15:04:05Z"))

		record.Set("name", fmt.Sprintf("%s-%d", strings.ToUpper(ticketType.GetString("singular")), number))
		record.Set("type", ticketType.Id)
		record.Set("description", fakeTicketDescription())
		record.Set("open", gofakeit.Bool())
		record.Set("schema", `{"type":"object","properties":{"tlp":{"title":"TLP","type":"string"}}}`)
		record.Set("state", `{"severity":"Medium"}`)
		record.Set("owner", random(users).Id)

		records = append(records, record)

		// Add comments
		records = append(records, commentRecords(app, users, created, record)...)
		records = append(records, timelineRecords(app, created, record)...)
		records = append(records, taskRecords(app, users, created, record)...)
		records = append(records, linkRecords(app, created, record)...)
	}

	return records
}

func commentRecords(app core.App, users []*core.Record, created time.Time, record *core.Record) []*core.Record {
	commentCollection, err := app.FindCollectionByNameOrId(migrations.CommentCollectionName)
	if err != nil {
		panic(err)
	}

	records := make([]*core.Record, 0, 5)

	for range gofakeit.IntN(5) {
		commentCreated := gofakeit.DateRange(created, time.Now())
		commentUpdated := gofakeit.DateRange(commentCreated, time.Now())

		commentRecord := core.NewRecord(commentCollection)
		commentRecord.Id = "c_" + security.PseudorandomString(10)
		commentRecord.Set("created", commentCreated.Format("2006-01-02T15:04:05Z"))
		commentRecord.Set("updated", commentUpdated.Format("2006-01-02T15:04:05Z"))
		commentRecord.Set("ticket", record.Id)
		commentRecord.Set("author", random(users).Id)
		commentRecord.Set("message", fakeTicketComment())

		records = append(records, commentRecord)
	}

	return records
}

func timelineRecords(app core.App, created time.Time, record *core.Record) []*core.Record {
	timelineCollection, err := app.FindCollectionByNameOrId(migrations.TimelineCollectionName)
	if err != nil {
		panic(err)
	}

	records := make([]*core.Record, 0, 5)

	for range gofakeit.IntN(5) {
		timelineCreated := gofakeit.DateRange(created, time.Now())
		timelineUpdated := gofakeit.DateRange(timelineCreated, time.Now())

		timelineRecord := core.NewRecord(timelineCollection)
		timelineRecord.Id = "tl_" + security.PseudorandomString(10)
		timelineRecord.Set("created", timelineCreated.Format("2006-01-02T15:04:05Z"))
		timelineRecord.Set("updated", timelineUpdated.Format("2006-01-02T15:04:05Z"))
		timelineRecord.Set("ticket", record.Id)
		timelineRecord.Set("time", gofakeit.DateRange(created, time.Now()).Format("2006-01-02T15:04:05Z"))
		timelineRecord.Set("message", fakeTicketTimelineMessage())

		records = append(records, timelineRecord)
	}

	return records
}

func taskRecords(app core.App, users []*core.Record, created time.Time, record *core.Record) []*core.Record {
	taskCollection, err := app.FindCollectionByNameOrId(migrations.TaskCollectionName)
	if err != nil {
		panic(err)
	}

	records := make([]*core.Record, 0, 5)

	for range gofakeit.IntN(5) {
		taskCreated := gofakeit.DateRange(created, time.Now())
		taskUpdated := gofakeit.DateRange(taskCreated, time.Now())

		taskRecord := core.NewRecord(taskCollection)
		taskRecord.Id = "ts_" + security.PseudorandomString(10)
		taskRecord.Set("created", taskCreated.Format("2006-01-02T15:04:05Z"))
		taskRecord.Set("updated", taskUpdated.Format("2006-01-02T15:04:05Z"))
		taskRecord.Set("ticket", record.Id)
		taskRecord.Set("name", fakeTicketTask())
		taskRecord.Set("open", gofakeit.Bool())
		taskRecord.Set("owner", random(users).Id)

		records = append(records, taskRecord)
	}

	return records
}

func linkRecords(app core.App, created time.Time, record *core.Record) []*core.Record {
	linkCollection, err := app.FindCollectionByNameOrId(migrations.LinkCollectionName)
	if err != nil {
		panic(err)
	}

	records := make([]*core.Record, 0, 5)

	for range gofakeit.IntN(5) {
		linkCreated := gofakeit.DateRange(created, time.Now())
		linkUpdated := gofakeit.DateRange(linkCreated, time.Now())

		linkRecord := core.NewRecord(linkCollection)
		linkRecord.Id = "l_" + security.PseudorandomString(10)
		linkRecord.Set("created", linkCreated.Format("2006-01-02T15:04:05Z"))
		linkRecord.Set("updated", linkUpdated.Format("2006-01-02T15:04:05Z"))
		linkRecord.Set("ticket", record.Id)
		linkRecord.Set("url", gofakeit.URL())
		linkRecord.Set("name", random([]string{"Blog", "Forum", "Wiki", "Documentation"}))

		records = append(records, linkRecord)
	}

	return records
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

const (
	triggerSchedule = `{"expression":"12 * * * *"}`
	triggerWebhook  = `{"token":"1234567890","path":"webhook"}`
	triggerHook     = `{"collections":["tickets"],"events":["create"]}`
)

func reactionRecords(app core.App) []*core.Record {
	var records []*core.Record

	collection, err := app.FindCollectionByNameOrId(migrations.ReactionCollectionName)
	if err != nil {
		panic(err)
	}

	createTicketActionData, err := json.Marshal(map[string]interface{}{
		"requirements": "pocketbase",
		"script":       createTicketPy,
	})
	if err != nil {
		panic(err)
	}

	record := core.NewRecord(collection)
	record.Id = "w_" + security.PseudorandomString(10)
	record.Set("name", "Create New Ticket")
	record.Set("trigger", "schedule")
	record.Set("triggerdata", triggerSchedule)
	record.Set("action", "python")
	record.Set("actiondata", string(createTicketActionData))

	records = append(records, record)

	alertIngestActionData, err := json.Marshal(map[string]interface{}{
		"requirements": "pocketbase",
		"script":       alertIngestPy,
	})
	if err != nil {
		panic(err)
	}

	record = core.NewRecord(collection)
	record.Id = "w_" + security.PseudorandomString(10)
	record.Set("name", "Alert Ingest Webhook")
	record.Set("trigger", "webhook")
	record.Set("triggerdata", triggerWebhook)
	record.Set("action", "python")
	record.Set("actiondata", string(alertIngestActionData))

	records = append(records, record)

	assignTicketsActionData, err := json.Marshal(map[string]interface{}{
		"requirements": "pocketbase",
		"script":       assignTicketsPy,
	})
	if err != nil {
		panic(err)
	}

	record = core.NewRecord(collection)
	record.Id = "w_" + security.PseudorandomString(10)
	record.Set("name", "Assign new Tickets")
	record.Set("trigger", "hook")
	record.Set("triggerdata", triggerHook)
	record.Set("action", "python")
	record.Set("actiondata", string(assignTicketsActionData))

	records = append(records, record)

	return records
}
