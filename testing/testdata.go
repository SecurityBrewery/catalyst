package testing

import (
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/SecurityBrewery/catalyst/app2"
	"github.com/SecurityBrewery/catalyst/app2/database/sqlc"
)

const (
	adminEmail   = "admin@catalyst-soar.com"
	analystEmail = "analyst@catalyst-soar.com"
)

func defaultTestData(t *testing.T, app *app2.App2) {
	t.Helper()

	userTestData(t, app)
	ticketTestData(t, app)
	reactionTestData(t, app)
}

func userTestData(t *testing.T, app *app2.App2) {
	t.Helper()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	app.Queries.CreateUser(t.Context(), sqlc.CreateUserParams{
		Username:     "u_bob_analyst",
		Email:        analystEmail,
		Verified:     true,
		Name:         "Bob Analyst",
		PasswordHash: string(passwordHash),
	})
}

func ticketTestData(t *testing.T, app *app2.App2) {
	t.Helper()

	app.Queries.CreateTicket(t.Context(), sqlc.CreateTicketParams{
		Name:        "Test Ticket",
		Type:        "incident",
		Description: "This is a test ticket.",
		Open:        true,
		Owner:       "u_bob_analyst",
		Schema:      `{"type":"object","properties":{"tlp":{"title":"TLP","type":"string"}}}`,
		State:       `{"tlp":"AMBER"}`,
	})
}

func reactionTestData(t *testing.T, app *app2.App2) {
	t.Helper()

	app.Queries.CreateReaction(t.Context(), sqlc.CreateReactionParams{
		Name:        "Reaction",
		Trigger:     "webhook",
		Triggerdata: `{"token":"1234567890","path":"test"}`,
		Action:      "python",
		Actiondata:  `{"requirements":"requests","script":"print('Hello, World!')"}`,
	})

	app.Queries.CreateReaction(t.Context(), sqlc.CreateReactionParams{
		Name:        "Reaction",
		Trigger:     "webhook",
		Triggerdata: `{"path":"test2"}`,
		Action:      "webhook",
		Actiondata:  `{"headers":{"Content-Type":"application/json"},"url":"http://127.0.0.1:12345/webhook"}`,
	})

	app.Queries.CreateReaction(t.Context(), sqlc.CreateReactionParams{
		Name:        "Hook",
		Trigger:     "hook",
		Triggerdata: `{"collections":["tickets"],"events":["create"]}`,
		Action:      "python",
		Actiondata:  `{"requirements":"requests","script":"import requests\nrequests.post('http://127.0.0.1:12346/test', json={'test':True})"}`,
	})
}
