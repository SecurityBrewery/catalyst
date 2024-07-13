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
	record.Set("triggerdata", `{"path":"test"}`)
	record.Set("reaction", "python")
	record.Set("reactiondata", `{"script":"print('Hello, World!')"}`)

	if err := app.Dao().SaveRecord(record); err != nil {
		t.Fatal(err)
	}
}
