package testing

import (
	"testing"

	"github.com/pocketbase/pocketbase/core"

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

	superusers, err := app.FindCollectionByNameOrId(core.CollectionNameSuperusers)
	if err != nil {
		t.Fatalf("failed to find superusers collection: %v", err)
	}

	admin := core.NewRecord(superusers)
	admin.SetEmail(adminEmail)
	admin.SetPassword("password")

	if err := app.Save(admin); err != nil {
		t.Fatalf("failed to save admin record: %v", err)
	}
}

func userTestData(t *testing.T, app core.App) {
	t.Helper()

	collection, err := app.FindCollectionByNameOrId(migrations.UserCollectionID)
	if err != nil {
		t.Fatalf("failed to find user collection: %v", err)
	}

	record := core.NewRecord(collection)
	record.Id = "u_bob_analyst"
	record.Set("username", "u_bob_analyst")
	record.SetPassword("password")
	record.Set("name", "Bob Analyst")
	record.Set("email", analystEmail)
	record.SetVerified(true)

	if err := app.SaveNoValidate(record); err != nil {
		t.Fatalf("failed to save user record: %v", err)
	}
}

func ticketTestData(t *testing.T, app core.App) {
	t.Helper()

	collection, err := app.FindCollectionByNameOrId(migrations.TicketCollectionName)
	if err != nil {
		t.Fatalf("failed to find ticket collection: %v", err)
	}

	record := core.NewRecord(collection)
	record.Id = "t_test"
	record.Set("name", "Test Ticket")
	record.Set("type", "incident")
	record.Set("description", "This is a test ticket.")
	record.Set("open", true)
	record.Set("schema", `{"type":"object","properties":{"tlp":{"title":"TLP","type":"string"}}}`)
	record.Set("state", `{"tlp":"AMBER"}`)
	record.Set("owner", "u_bob_analyst")

	if err := app.SaveNoValidate(record); err != nil {
		t.Fatalf("failed to save ticket record: %v", err)
	}
}

func reactionTestData(t *testing.T, app core.App) {
	t.Helper()

	collection, err := app.FindCollectionByNameOrId(migrations.ReactionCollectionName)
	if err != nil {
		t.Fatalf("failed to find reaction collection: %v", err)
	}

	record := core.NewRecord(collection)
	record.Id = "r_reaction"
	record.Set("name", "Reaction")
	record.Set("trigger", "webhook")
	record.Set("triggerdata", `{"token":"1234567890","path":"test"}`)
	record.Set("action", "python")
	record.Set("actiondata", `{"requirements":"requests","script":"print('Hello, World!')"}`)

	if err := app.SaveNoValidate(record); err != nil {
		t.Fatalf("failed to save reaction record: %v", err)
	}

	record = core.NewRecord(collection)
	record.Id = "r_reaction_webhook"
	record.Set("name", "Reaction")
	record.Set("trigger", "webhook")
	record.Set("triggerdata", `{"path":"test2"}`)
	record.Set("action", "webhook")
	record.Set("actiondata", `{"headers":{"Content-Type":"application/json"},"url":"http://127.0.0.1:12345/webhook"}`)

	if err := app.SaveNoValidate(record); err != nil {
		t.Fatalf("failed to save reaction record: %v", err)
	}

	record = core.NewRecord(collection)
	record.Id = "r_reaction_hook"
	record.Set("name", "Hook")
	record.Set("trigger", "hook")
	record.Set("triggerdata", `{"collections":["tickets"],"events":["create"]}`)
	record.Set("action", "python")
	record.Set("actiondata", `{"requirements":"requests","script":"import requests\nrequests.post('http://127.0.0.1:12346/test', json={'test':True})"}`)

	if err := app.SaveNoValidate(record); err != nil {
		t.Fatalf("failed to save reaction record: %v", err)
	}
}
