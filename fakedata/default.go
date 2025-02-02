package fakedata

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"

	"github.com/SecurityBrewery/catalyst/migrations"
)

func defaultData() map[string]map[string]map[string]any {
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

	return map[string]map[string]map[string]any{
		migrations.TicketCollectionName: {
			"t_0": {
				"created":     dateTime(ticketCreated),
				"updated":     dateTime(ticketUpdated),
				"name":        "phishing-123",
				"type":        "alert",
				"description": "Phishing email reported by several employees.",
				"open":        true,
				"schema":      types.JsonRaw(`{"type":"object","properties":{"tlp":{"title":"TLP","type":"string"}}}`),
				"state":       types.JsonRaw(`{"severity":"Medium"}`),
				"owner":       "u_test",
			},
		},
		migrations.CommentCollectionName: {
			"c_0": {
				"created": dateTime(commentCreated),
				"updated": dateTime(commentUpdated),
				"ticket":  "t_0",
				"author":  "u_test",
				"message": "This is a test comment.",
			},
		},
		migrations.TimelineCollectionName: {
			"tl_0": {
				"created": dateTime(timelineCreated),
				"updated": dateTime(timelineUpdated),
				"ticket":  "t_0",
				"time":    dateTime(timelineCreated),
				"message": "This is a test timeline message.",
			},
		},
		migrations.TaskCollectionName: {
			"ts_0": {
				"created": dateTime(taskCreated),
				"updated": dateTime(taskUpdated),
				"ticket":  "t_0",
				"name":    "This is a test task.",
				"open":    true,
				"owner":   "u_test",
			},
		},
		migrations.LinkCollectionName: {
			"l_0": {
				"created": dateTime(linkCreated),
				"updated": dateTime(linkUpdated),
				"ticket":  "t_0",
				"url":     "https://www.example.com",
				"name":    "This is a test link.",
			},
		},
		migrations.ReactionCollectionName: {
			"w_0": {
				"created":     dateTime(reactionCreated),
				"updated":     dateTime(reactionUpdated),
				"name":        "Create New Ticket",
				"trigger":     "schedule",
				"triggerdata": types.JsonRaw(triggerSchedule),
				"action":      "python",
				"actiondata":  types.JsonRaw(createTicketActionData),
			},
		},
	}
}

func GenerateDefaultData(app core.App) error {
	var records []*models.Record

	// users
	userRecord, err := testUser(app.Dao())
	if err != nil {
		return err
	}

	records = append(records, userRecord)

	// records
	for collectionName, collectionRecords := range defaultData() {
		collection, err := app.Dao().FindCollectionByNameOrId(collectionName)
		if err != nil {
			return err
		}

		for id, fields := range collectionRecords {
			record := models.NewRecord(collection)
			record.SetId(id)

			for key, value := range fields {
				record.Set(key, value)
			}

			records = append(records, record)
		}
	}

	for _, record := range records {
		if err := app.Dao().SaveRecord(record); err != nil {
			return err
		}
	}

	return nil
}

func ValidateDefaultData(app core.App) error { //nolint:cyclop,gocognit
	// users
	userRecord, err := app.Dao().FindRecordById(migrations.UserCollectionName, "u_test")
	if err != nil {
		return fmt.Errorf("failed to find user record: %w", err)
	}

	if userRecord == nil {
		return errors.New("user not found")
	}

	if userRecord.Username() != "u_test" {
		return fmt.Errorf(`username does not match: got %q, want "u_test"`, userRecord.Username())
	}

	if !userRecord.ValidatePassword("1234567890") {
		return errors.New("password does not match")
	}

	if userRecord.Get("name") != "Test User" {
		return fmt.Errorf(`name does not match: got %q, want "Test User"`, userRecord.Get("name"))
	}

	if userRecord.Get("email") != "user@catalyst-soar.com" {
		return fmt.Errorf(`email does not match: got %q, want "user@catalyst-soar.com"`, userRecord.Get("email"))
	}

	if !userRecord.Verified() {
		return errors.New("user is not verified")
	}

	// records
	for collectionName, collectionRecords := range defaultData() {
		for id, fields := range collectionRecords {
			record, err := app.Dao().FindRecordById(collectionName, id)
			if err != nil {
				return fmt.Errorf("failed to find record %s: %w", id, err)
			}

			if record == nil {
				return errors.New("record not found")
			}

			for key, value := range fields {
				got := record.Get(key)

				if wantJSON, ok := value.(types.JsonRaw); ok {
					if err := compareJSON(got, wantJSON); err != nil {
						return fmt.Errorf("record field %q does not match: %w", key, err)
					}

					continue
				}

				if got != value {
					return fmt.Errorf("record field %s does not match: got %v (%T), want %v (%T)", key, got, got, value, value)
				}
			}
		}
	}

	return nil
}

func compareJSON(got any, wantJSON types.JsonRaw) error {
	gotJSON, ok := got.(types.JsonRaw)
	if !ok {
		return fmt.Errorf("got %T, want %T", got, wantJSON)
	}

	if !jsonEqual(gotJSON.String(), wantJSON.String()) {
		return fmt.Errorf("got %v, want %v", gotJSON, wantJSON)
	}

	return nil
}

func jsonEqual(a, b string) bool {
	var objA, objB interface{}

	if err := json.Unmarshal([]byte(a), &objA); err != nil {
		return false
	}

	if err := json.Unmarshal([]byte(b), &objB); err != nil {
		return false
	}

	return reflect.DeepEqual(objA, objB)
}

func dateTime(t time.Time) types.DateTime {
	dt := types.DateTime{}
	_ = dt.Scan(t)

	return dt
}
