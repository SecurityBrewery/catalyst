package fakedata

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/security"

	"github.com/SecurityBrewery/catalyst/migrations"
	"github.com/SecurityBrewery/catalyst/playbook"
)

const (
	minimumUserCount   = 1
	minimumTicketCount = 1
)

func Generate(db dbx.Builder, userCount, ticketCount int) error {
	dao := daos.New(db)

	if userCount < minimumUserCount {
		userCount = minimumUserCount
	}

	if ticketCount < minimumTicketCount {
		ticketCount = minimumTicketCount
	}

	types, err := dao.FindRecordsByExpr(migrations.TypeCollectionName)
	if err != nil {
		return err
	}

	users := userRecords(dao, userCount)
	playbooks := playbookRecords(dao)
	tickets := ticketRecords(dao, users, types, playbooks, ticketCount)
	webhooks := webhookRecords(dao)

	for _, records := range [][]*models.Record{users, tickets, webhooks, playbooks} {
		for _, record := range records {
			if err := dao.SaveRecord(record); err != nil {
				return err
			}
		}
	}

	return nil
}

func userRecords(dao *daos.Dao, count int) []*models.Record {
	collection, err := dao.FindCollectionByNameOrId(migrations.UserCollectionName)
	if err != nil {
		panic(err)
	}

	var records []*models.Record

	// create the test user
	if _, err := dao.FindRecordById(migrations.UserCollectionName, "u_test"); err != nil {
		record := models.NewRecord(collection)
		record.SetId("u_test")
		record.SetUsername("u_test")
		record.SetPassword("1234567890")
		record.Set("name", gofakeit.Name())
		record.Set("email", "user@catalyst-soar.com")
		record.SetVerified(true)

		records = append(records, record)
	}

	for range count - 1 {
		record := models.NewRecord(collection)
		record.SetId("u_" + security.PseudorandomString(5))
		record.SetUsername("u_" + security.RandomStringWithAlphabet(5, "123456789"))
		record.SetPassword("1234567890")
		record.Set("name", gofakeit.Name())
		record.Set("email", gofakeit.Username()+"@catalyst-soar.com")
		record.SetVerified(true)

		records = append(records, record)
	}

	return records
}

func ticketRecords(dao *daos.Dao, users, types, playbooks []*models.Record, count int) []*models.Record {
	collection, err := dao.FindCollectionByNameOrId(migrations.TicketCollectionName)
	if err != nil {
		panic(err)
	}

	var records []*models.Record

	created := time.Now()
	number := gofakeit.Number(200*count, 300*count)

	for range count {
		number -= gofakeit.Number(100, 200)
		created = created.Add(time.Duration(-gofakeit.Number(13, 37)) * time.Hour)

		record := models.NewRecord(collection)
		record.SetId("t_" + security.PseudorandomString(5))

		updated := gofakeit.DateRange(created, time.Now())

		ticketType := random(types)

		record.Set("created", created.Format("2006-01-02T15:04:05Z"))
		record.Set("updated", updated.Format("2006-01-02T15:04:05Z"))

		record.Set("name", fmt.Sprintf("%s-%d", strings.ToUpper(ticketType.GetString("singular")), number))
		record.Set("type", ticketType.GetId())
		record.Set("description", fakeTicketDescription())
		record.Set("open", gofakeit.Bool())
		record.Set("schema", `{"type":"object","properties":{"tlp":{"title":"TLP","type":"string"}}}`)
		record.Set("state", `{"tlp":"AMBER"}`)
		record.Set("owner", random(users).GetId())

		records = append(records, record)

		// Add comments
		for range gofakeit.IntN(5) {
			commentCollection, err := dao.FindCollectionByNameOrId(migrations.CommentCollectionName)
			if err != nil {
				panic(err)
			}

			commentCreated := gofakeit.DateRange(created, time.Now())
			commentUpdated := gofakeit.DateRange(commentCreated, time.Now())

			commentRecord := models.NewRecord(commentCollection)
			commentRecord.SetId("c_" + security.PseudorandomString(5))
			commentRecord.Set("created", commentCreated.Format("2006-01-02T15:04:05Z"))
			commentRecord.Set("updated", commentUpdated.Format("2006-01-02T15:04:05Z"))
			commentRecord.Set("ticket", record.GetId())
			commentRecord.Set("author", random(users).GetId())
			commentRecord.Set("message", fakeTicketComment())

			records = append(records, commentRecord)
		}

		// Add timeline
		for range gofakeit.IntN(5) {
			timelineCollection, err := dao.FindCollectionByNameOrId(migrations.TimelineCollectionName)
			if err != nil {
				panic(err)
			}

			timelineCreated := gofakeit.DateRange(created, time.Now())
			timelineUpdated := gofakeit.DateRange(timelineCreated, time.Now())

			timelineRecord := models.NewRecord(timelineCollection)
			timelineRecord.SetId("tl_" + security.PseudorandomString(5))
			timelineRecord.Set("created", timelineCreated.Format("2006-01-02T15:04:05Z"))
			timelineRecord.Set("updated", timelineUpdated.Format("2006-01-02T15:04:05Z"))
			timelineRecord.Set("ticket", record.GetId())
			timelineRecord.Set("time", gofakeit.DateRange(created, time.Now()).Format("2006-01-02T15:04:05Z"))
			timelineRecord.Set("message", fakeTicketTimelineMessage())

			records = append(records, timelineRecord)
		}

		// Add tasks
		for range gofakeit.IntN(5) {
			taskCollection, err := dao.FindCollectionByNameOrId(migrations.TaskCollectionName)
			if err != nil {
				panic(err)
			}

			taskCreated := gofakeit.DateRange(created, time.Now())
			taskUpdated := gofakeit.DateRange(taskCreated, time.Now())

			taskRecord := models.NewRecord(taskCollection)
			taskRecord.SetId("ts_" + security.PseudorandomString(5))
			taskRecord.Set("created", taskCreated.Format("2006-01-02T15:04:05Z"))
			taskRecord.Set("updated", taskUpdated.Format("2006-01-02T15:04:05Z"))
			taskRecord.Set("ticket", record.GetId())
			taskRecord.Set("name", fakeTicketTask())
			taskRecord.Set("open", gofakeit.Bool())
			taskRecord.Set("owner", random(users).GetId())

			records = append(records, taskRecord)
		}

		// Add links
		for range gofakeit.IntN(5) {
			linkCollection, err := dao.FindCollectionByNameOrId(migrations.LinkCollectionName)
			if err != nil {
				panic(err)
			}

			linkCreated := gofakeit.DateRange(created, time.Now())
			linkUpdated := gofakeit.DateRange(linkCreated, time.Now())

			linkRecord := models.NewRecord(linkCollection)
			linkRecord.SetId("l_" + security.PseudorandomString(5))
			linkRecord.Set("created", linkCreated.Format("2006-01-02T15:04:05Z"))
			linkRecord.Set("updated", linkUpdated.Format("2006-01-02T15:04:05Z"))
			linkRecord.Set("ticket", record.GetId())
			linkRecord.Set("url", gofakeit.URL())
			linkRecord.Set("name", random([]string{"Blog", "Forum", "Wiki", "Documentation"}))

			records = append(records, linkRecord)
		}

		// Add runs
		for range gofakeit.IntN(2) {
			runCollection, err := dao.FindCollectionByNameOrId(migrations.RunCollectionName)
			if err != nil {
				panic(err)
			}

			runCreated := gofakeit.DateRange(created, time.Now())
			runUpdated := gofakeit.DateRange(runCreated, time.Now())

			runPlaybook := random(playbooks)
			runPlaybookStepsJSON := runPlaybook.GetString("steps")

			var runPlaybookSteps []playbook.Step

			if err := json.Unmarshal([]byte(runPlaybookStepsJSON), &runPlaybookSteps); err != nil {
				continue
			}

			var steps []any
			for _, step := range runPlaybookSteps {
				steps = append(steps, map[string]any{
					"name":        step.Name,
					"type":        step.Type,
					"status":      gofakeit.RandomString([]string{"open", "completed"}),
					"description": step.Description,
					"schema":      step.Schema,
					"state":       map[string]any{},
				})
			}

			runRecord := models.NewRecord(runCollection)
			runRecord.SetId("r_" + security.PseudorandomString(5))
			runRecord.Set("created", runCreated.Format("2006-01-02T15:04:05Z"))
			runRecord.Set("updated", runUpdated.Format("2006-01-02T15:04:05Z"))
			runRecord.Set("ticket", record.GetId())
			runRecord.Set("name", runPlaybook.Get("name"))
			runRecord.Set("steps", steps)

			records = append(records, runRecord)
		}
	}

	return records
}

func webhookRecords(dao *daos.Dao) []*models.Record {
	collection, err := dao.FindCollectionByNameOrId(migrations.WebhookCollectionName)
	if err != nil {
		panic(err)
	}

	// var records []*models.Record

	record := models.NewRecord(collection)
	record.SetId("w_" + security.PseudorandomString(5))
	record.Set("name", "Test Webhook")
	record.Set("collection", "tickets")
	record.Set("destination", "http://localhost:8080/webhook")

	return []*models.Record{record}
}

func playbookRecords(dao *daos.Dao) []*models.Record {
	playbookCollection, err := dao.FindCollectionByNameOrId(migrations.PlaybookCollectionName)
	if err != nil {
		panic(err)
	}

	playbook1 := []playbook.Step{
		{
			Name:        "Detection and Isolation",
			Type:        "task",
			Description: "Monitor and identify malware presence, then isolate affected systems from the network.",
			// Schema: playbook.JSONSchema{
			// 	Type: "object",
			// 	Properties: map[string]playbook.JSONProperty{
			// 		"name": {
			// 			Title: "Name",
			// 			Type:  "string",
			// 		},
			// 	},
			// },
		},
		{
			Name:        "Containment and Mitigation",
			Type:        "task",
			Description: "Contain the malware and mitigate the impact on the network.",
		},
		{
			Name:        "Eradication and Recovery",
			Type:        "task",
			Description: "Remove the malware from the network and recover affected systems.",
		},
		{
			Name:        "Post-Incident Analysis",
			Type:        "task",
			Description: "Analyze the incident and identify the root cause of the malware infection.",
		},
	}

	b, err := json.Marshal(playbook1)
	if err != nil {
		panic(err)
	}

	record1 := models.NewRecord(playbookCollection)
	record1.SetId("p_" + security.PseudorandomString(5))
	record1.Set("name", "Malware Infection")
	record1.Set("steps", string(b))

	playbook2 := []playbook.Step{
		{
			Name:        "Data Collection",
			Type:        "task",
			Description: "Collect customer information and store it in a secure location.",
		},
		{
			Name:        "Data Analysis",
			Type:        "task",
			Description: "Analyze the collected data and identify patterns and trends.",
		},
		{
			Name:        "Data Reporting",
			Type:        "task",
			Description: "Generate reports based on the analyzed data and share them with the customer.",
		},
	}

	b, err = json.Marshal(playbook2)
	if err != nil {
		panic(err)
	}

	record2 := models.NewRecord(playbookCollection)
	record2.SetId("p_" + security.PseudorandomString(5))
	record2.Set("name", "Customer information")
	record2.Set("steps", string(b))

	return []*models.Record{record1, record2}
}
