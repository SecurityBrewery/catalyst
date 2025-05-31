package testing

import (
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/SecurityBrewery/catalyst/app"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

const (
	adminEmail   = "admin@catalyst-soar.com"
	analystEmail = "analyst@catalyst-soar.com"
)

func defaultTestData(t *testing.T, app *app.App) {
	t.Helper()

	userTestData(t, app)
	ticketTestData(t, app)
	reactionTestData(t, app)
}

func userTestData(t *testing.T, app *app.App) {
	t.Helper()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	_, err = app.Queries.CreateUser(t.Context(), sqlc.CreateUserParams{
		ID:           "u_bob_analyst",
		Username:     "u_bob_analyst",
		Email:        analystEmail,
		Verified:     true,
		Name:         "Bob Analyst",
		PasswordHash: string(passwordHash),
	})
	if err != nil {
		t.Fatal(err)
	}
}

func ticketTestData(t *testing.T, app *app.App) {
	t.Helper()

	_, err := app.Queries.CreateTicket(t.Context(), sqlc.CreateTicketParams{
		ID:          "test-ticket",
		Name:        "Test Ticket",
		Type:        "incident",
		Description: "This is a test ticket.",
		Open:        true,
		Owner:       "u_bob_analyst",
		Schema:      `{"type":"object","properties":{"tlp":{"title":"TLP","type":"string"}}}`,
		State:       `{"tlp":"AMBER"}`,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func reactionTestData(t *testing.T, app *app.App) {
	t.Helper()

	if _, err := app.Queries.CreateReaction(t.Context(), sqlc.CreateReactionParams{
		ID:          "r-test-webhook",
		Name:        "Reaction",
		Trigger:     "webhook",
		Triggerdata: `{"token":"1234567890","path":"test"}`,
		Action:      "python",
		Actiondata:  `{"requirements":"requests","script":"print('Hello, World!')"}`,
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := app.Queries.CreateReaction(t.Context(), sqlc.CreateReactionParams{
		ID:          "r-test-proxy",
		Name:        "Reaction",
		Trigger:     "webhook",
		Triggerdata: `{"path":"test2"}`,
		Action:      "webhook",
		Actiondata:  `{"headers":{"Content-Type":"application/json"},"url":"http://127.0.0.1:12345/webhook"}`,
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := app.Queries.CreateReaction(t.Context(), sqlc.CreateReactionParams{
		ID:          "r-test-hook",
		Name:        "Hook",
		Trigger:     "hook",
		Triggerdata: `{"collections":["tickets"],"events":["create"]}`,
		Action:      "python",
		Actiondata:  `{"requirements":"requests","script":"import requests\nrequests.post('http://127.0.0.1:12346/test', json={'test':True})"}`,
	}); err != nil {
		t.Fatal(err)
	}
}
