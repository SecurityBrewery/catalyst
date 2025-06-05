package database

import (
	"testing"

	"github.com/SecurityBrewery/catalyst/app/auth/password"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/permission"
)

const (
	AdminEmail   = "admin@catalyst-soar.com"
	AnalystEmail = "analyst@catalyst-soar.com"
)

func DefaultTestData(t *testing.T, queries *sqlc.Queries) {
	t.Helper()

	userTestData(t, queries)
	ticketTestData(t, queries)
	reactionTestData(t, queries)
}

func userTestData(t *testing.T, queries *sqlc.Queries) {
	t.Helper()

	_, err := queries.CreateRole(t.Context(), sqlc.CreateRoleParams{
		ID:          "r_analyst",
		Name:        "Analyst",
		Permissions: permission.ToJSONArray(t.Context(), permission.Default()),
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = queries.CreateRole(t.Context(), sqlc.CreateRoleParams{
		ID:          "r_admin",
		Name:        "Admin",
		Permissions: permission.ToJSONArray(t.Context(), permission.AllPermissions()),
	})
	if err != nil {
		t.Fatal(err)
	}

	passwordHash, tokenKey, err := password.Hash("password")
	if err != nil {
		t.Fatal(err)
	}

	_, err = queries.CreateUser(t.Context(), sqlc.CreateUserParams{
		ID:           "u_bob_analyst",
		Username:     "u_bob_analyst",
		Email:        AnalystEmail,
		Verified:     true,
		Name:         "Bob Analyst",
		PasswordHash: passwordHash,
		TokenKey:     tokenKey,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = queries.AssignRoleToUser(t.Context(), sqlc.AssignRoleToUserParams{
		UserID: "u_bob_analyst",
		RoleID: "r_analyst",
	})
	if err != nil {
		t.Fatal(err)
	}

	passwordHash, tokenKey, err = password.Hash("password123")
	if err != nil {
		t.Fatal(err)
	}

	_, err = queries.CreateUser(t.Context(), sqlc.CreateUserParams{
		ID:           "u_admin",
		Username:     "u_admin",
		Email:        AdminEmail,
		Verified:     true,
		Name:         "Admin User",
		PasswordHash: passwordHash,
		TokenKey:     tokenKey,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = queries.AssignRoleToUser(t.Context(), sqlc.AssignRoleToUserParams{
		UserID: "u_admin",
		RoleID: "r_admin",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func ticketTestData(t *testing.T, queries *sqlc.Queries) {
	t.Helper()

	ticket, err := queries.CreateTicket(t.Context(), sqlc.CreateTicketParams{
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

	if _, err := queries.CreateComment(t.Context(), sqlc.CreateCommentParams{
		ID:      "c_test_comment",
		Author:  "u_bob_analyst",
		Message: "Initial comment on the test ticket.",
		Ticket:  ticket.ID,
	}); err != nil {
		t.Fatal(err)
	}
}

func reactionTestData(t *testing.T, queries *sqlc.Queries) {
	t.Helper()

	if _, err := queries.CreateReaction(t.Context(), sqlc.CreateReactionParams{
		ID:          "r-test-webhook",
		Name:        "Reaction",
		Trigger:     "webhook",
		Triggerdata: `{"token":"1234567890","path":"test"}`,
		Action:      "python",
		Actiondata:  `{"requirements":"requests","script":"print('Hello, World!')"}`,
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := queries.CreateReaction(t.Context(), sqlc.CreateReactionParams{
		ID:          "r-test-proxy",
		Name:        "Reaction",
		Trigger:     "webhook",
		Triggerdata: `{"path":"test2"}`,
		Action:      "webhook",
		Actiondata:  `{"headers":{"Content-Type":"application/json"},"url":"http://127.0.0.1:12345/webhook"}`,
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := queries.CreateReaction(t.Context(), sqlc.CreateReactionParams{
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
