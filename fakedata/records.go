package fakedata

import (
	"fmt"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/security"

	"github.com/SecurityBrewery/catalyst/migrations"
)

const (
	minimumUserCount   = 1
	minimumTicketCount = 1
)

func Generate(app core.App, userCount, ticketCount int) error {
	if userCount < minimumUserCount {
		userCount = minimumUserCount
	}

	if ticketCount < minimumTicketCount {
		ticketCount = minimumTicketCount
	}

	types, err := app.Dao().FindRecordsByExpr(migrations.TypeCollectionName)
	if err != nil {
		return err
	}

	users := userRecords(app.Dao(), userCount)
	tickets := ticketRecords(app.Dao(), users, types, ticketCount)
	webhooks := webhookRecords(app.Dao())
	reactions := reactionRecords(app.Dao())

	for _, records := range [][]*models.Record{users, tickets, webhooks, reactions} {
		for _, record := range records {
			if err := app.Dao().SaveRecord(record); err != nil {
				app.Logger().Error(err.Error())
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
		_ = record.SetUsername("u_test")
		_ = record.SetPassword("1234567890")
		record.Set("name", gofakeit.Name())
		record.Set("email", "user@catalyst-soar.com")
		_ = record.SetVerified(true)

		records = append(records, record)
	}

	for range count - 1 {
		record := models.NewRecord(collection)
		record.SetId("u_" + security.PseudorandomString(10))
		_ = record.SetUsername("u_" + security.RandomStringWithAlphabet(5, "123456789"))
		_ = record.SetPassword("1234567890")
		record.Set("name", gofakeit.Name())
		record.Set("email", gofakeit.Username()+"@catalyst-soar.com")
		_ = record.SetVerified(true)

		records = append(records, record)
	}

	return records
}

func ticketRecords(dao *daos.Dao, users, types []*models.Record, count int) []*models.Record {
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
		record.SetId("t_" + security.PseudorandomString(10))

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
			commentRecord.SetId("c_" + security.PseudorandomString(10))
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
			timelineRecord.SetId("tl_" + security.PseudorandomString(10))
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
			taskRecord.SetId("ts_" + security.PseudorandomString(10))
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
			linkRecord.SetId("l_" + security.PseudorandomString(10))
			linkRecord.Set("created", linkCreated.Format("2006-01-02T15:04:05Z"))
			linkRecord.Set("updated", linkUpdated.Format("2006-01-02T15:04:05Z"))
			linkRecord.Set("ticket", record.GetId())
			linkRecord.Set("url", gofakeit.URL())
			linkRecord.Set("name", random([]string{"Blog", "Forum", "Wiki", "Documentation"}))

			records = append(records, linkRecord)
		}
	}

	return records
}

func webhookRecords(dao *daos.Dao) []*models.Record {
	collection, err := dao.FindCollectionByNameOrId(migrations.WebhookCollectionName)
	if err != nil {
		panic(err)
	}

	record := models.NewRecord(collection)
	record.SetId("w_" + security.PseudorandomString(10))
	record.Set("name", "Test Webhook")
	record.Set("collection", "tickets")
	record.Set("destination", "http://localhost:8080/webhook")

	return []*models.Record{record}
}

const (
	triggerWebhook  = `{"token":"1234567890","path":"webhook"}`
	reactionPython  = `{"requirements":"requests","script":"import sys\n\nprint(sys.argv[1])"}`
	triggerHook     = `{"collections":["tickets","comments"],"events":["create","update","delete"]}`
	reactionWebhook = `{"headers":["Content-Type: application/json"],"url":"http://localhost:8080/webhook"}`
)

func reactionRecords(dao *daos.Dao) []*models.Record {
	var records []*models.Record

	collection, err := dao.FindCollectionByNameOrId(migrations.ReactionCollectionName)
	if err != nil {
		panic(err)
	}

	record := models.NewRecord(collection)
	record.SetId("w_" + security.PseudorandomString(10))
	record.Set("name", "Test Reaction")
	record.Set("trigger", "webhook")
	record.Set("triggerdata", triggerWebhook)
	record.Set("reaction", "python")
	record.Set("reactiondata", reactionPython)

	records = append(records, record)

	record = models.NewRecord(collection)
	record.SetId("w_" + security.PseudorandomString(10))
	record.Set("name", "Test Reaction 2")
	record.Set("trigger", "hook")
	record.Set("triggerdata", triggerHook)
	record.Set("reaction", "webhook")
	record.Set("reactiondata", reactionWebhook)

	records = append(records, record)

	return records
}
