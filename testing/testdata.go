package testing

import (
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"

	"github.com/SecurityBrewery/catalyst/migrations"
)

const (
	adminEmail   = "admin@catalyst-soar.com"
	analystEmail = "analyst@catalyst-soar.com"
)

func defaultTestData(t *testing.T, app core.App) {
	t.Helper()

	adminTestData(t, app)
	userTestData(t, app)
	ticketTestData(t, app)
	reactionTestData(t, app)
}

func adminTestData(t *testing.T, app core.App) {
	t.Helper()

	admin := &models.Admin{Email: adminEmail}

	if err := admin.SetPassword("password"); err != nil {
		t.Fatal(err)
	}

	if err := app.Dao().SaveAdmin(admin); err != nil {
		t.Fatal(err)
	}
}

func userTestData(t *testing.T, app core.App) {
	t.Helper()

	collection, err := app.Dao().FindCollectionByNameOrId(migrations.UserCollectionName)
	if err != nil {
		t.Fatal(err)
	}

	record := models.NewRecord(collection)
	record.SetId("u_bob_analyst")
	_ = record.SetUsername("u_bob_analyst")
	_ = record.SetPassword("password")
	record.Set("name", "Bob Analyst")
	record.Set("email", analystEmail)
	_ = record.SetVerified(true)

	if err := app.Dao().SaveRecord(record); err != nil {
		t.Fatal(err)
	}
}

func ticketTestData(t *testing.T, app core.App) {
	t.Helper()

	collection, err := app.Dao().FindCollectionByNameOrId(migrations.TicketCollectionName)
	if err != nil {
		t.Fatal(err)
	}

	record := models.NewRecord(collection)
	record.SetId("t_test")

	record.Set("name", "Test Ticket")
	record.Set("type", "incident")
	record.Set("description", "This is a test ticket.")
	record.Set("open", true)
	record.Set("schema", `{"type":"object","properties":{"tlp":{"title":"TLP","type":"string"}}}`)
	record.Set("state", `{"tlp":"AMBER"}`)
	record.Set("owner", "u_bob_analyst")

	if err := app.Dao().SaveRecord(record); err != nil {
		t.Fatal(err)
	}
}

func reactionTestData(t *testing.T, app core.App) {
	t.Helper()

	collection, err := app.Dao().FindCollectionByNameOrId(migrations.ReactionCollectionName)
	if err != nil {
		t.Fatal(err)
	}

	record := models.NewRecord(collection)
	record.SetId("r_reaction")
	record.Set("name", "Reaction")
	record.Set("trigger", "webhook")
	record.Set("triggerdata", `{"token":"1234567890","path":"test"}`)
	record.Set("action", "python")
	record.Set("actiondata", `{"requirements":"requests","script":"print('Hello, World!')"}`)

	if err := app.Dao().SaveRecord(record); err != nil {
		t.Fatal(err)
	}

	record = models.NewRecord(collection)
	record.SetId("r_reaction_webhook")
	record.Set("name", "Reaction")
	record.Set("trigger", "webhook")
	record.Set("triggerdata", `{"path":"test2"}`)
	record.Set("action", "webhook")
	record.Set("actiondata", `{"headers":{"Content-Type":"application/json"},"url":"http://127.0.0.1:12345/webhook"}`)

	if err := app.Dao().SaveRecord(record); err != nil {
		t.Fatal(err)
	}

	record = models.NewRecord(collection)
	record.SetId("r_reaction_hook")
	record.Set("name", "Hook")
	record.Set("trigger", "hook")
	record.Set("triggerdata", `{"collections":["tickets"],"events":["create"]}`)
	record.Set("action", "python")
	record.Set("actiondata", `{"requirements":"requests","script":"import requests\nrequests.post('http://127.0.0.1:12346/test', json={'test':True})"}`)

	if err := app.Dao().SaveRecord(record); err != nil {
		t.Fatal(err)
	}
}
