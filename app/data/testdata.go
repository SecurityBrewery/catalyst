package data

import (
	"encoding/json"
	"log/slog"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/pointer"
)

const (
	AdminEmail   = "admin@catalyst-soar.com"
	AnalystEmail = "analyst@catalyst-soar.com"
)

func DefaultTestData(t *testing.T, dir string, queries *sqlc.Queries) {
	t.Helper()

	parseTime := func(s string) time.Time {
		t, _ := time.Parse(time.RFC3339Nano, s)

		return t
	}

	ctx := t.Context()

	// Insert users
	_, err := queries.InsertUser(ctx, sqlc.InsertUserParams{
		Created:      parseTime("2025-06-21T22:21:26.271Z"),
		Updated:      parseTime("2025-06-21T22:21:26.271Z"),
		Email:        pointer.Pointer("analyst@catalyst-soar.com"),
		Username:     "u_bob_analyst",
		Name:         pointer.Pointer("Bob Analyst"),
		PasswordHash: "$2a$10$ZEHNh9ZKJ81N717wovDnMuLwZOLa6.g22IRzRr4goG6zGN.57UzJG",
		TokenKey:     "z3Jj8bbzcq_cSZs07XKoGlB0UtvmQiphHgwNkE4akoY=",
		Active:       true,
		ID:           "u_bob_analyst",
	})
	require.NoError(t, err, "failed to insert analyst user")

	_, err = queries.InsertUser(ctx, sqlc.InsertUserParams{
		Created:      parseTime("2025-06-21T22:21:26.271Z"),
		Updated:      parseTime("2025-06-21T22:21:26.271Z"),
		Email:        pointer.Pointer("admin@catalyst-soar.com"),
		Username:     "u_admin",
		Name:         pointer.Pointer("Admin User"),
		PasswordHash: "$2a$10$Z3/0HHWau6oi1t1aRPiI0uiVOWI.IosTAYEL0DJ2XJaalP9kesgBa",
		TokenKey:     "5BWDKLIAn3SQkpQlBUGrS_XEbFf91DsDpuh_Xmt4Nwg=",
		Active:       true,
		ID:           "u_admin",
	})
	require.NoError(t, err, "failed to insert admin user")

	// Insert webhooks
	_, err = queries.InsertWebhook(ctx, sqlc.InsertWebhookParams{
		ID:          "w_test_webhook",
		Name:        "Test Webhook",
		Collection:  "tickets",
		Destination: "https://example.com",
		Created:     parseTime("2025-06-21T22:21:26.271Z"),
		Updated:     parseTime("2025-06-21T22:21:26.271Z"),
	})
	require.NoError(t, err, "failed to insert webhook")

	// Insert types
	_, err = queries.InsertType(ctx, sqlc.InsertTypeParams{
		ID:       "test-type",
		Singular: "Test",
		Plural:   "Tests",
		Schema:   []byte(`{}`),
		Created:  parseTime("2025-06-21T22:21:26.271Z"),
		Updated:  parseTime("2025-06-21T22:21:26.271Z"),
	})
	require.NoError(t, err, "failed to insert type")

	// Insert tickets
	_, err = queries.InsertTicket(ctx, sqlc.InsertTicketParams{
		Created:     parseTime("2025-06-21T22:21:26.271Z"),
		Description: "This is a test ticket.",
		ID:          "test-ticket",
		Name:        "Test Ticket",
		Open:        true,
		Owner:       pointer.Pointer("u_bob_analyst"),
		Schema:      json.RawMessage([]byte(`{"type":"object","properties":{"tlp":{"title":"TLP","type":"string"}}}`)),
		State:       json.RawMessage([]byte(`{"tlp":"AMBER"}`)),
		Type:        "incident",
		Updated:     parseTime("2025-06-21T22:21:26.271Z"),
	})
	require.NoError(t, err, "failed to insert ticket")

	// Insert tasks
	_, err = queries.InsertTask(ctx, sqlc.InsertTaskParams{
		Created: parseTime("2025-06-21T22:21:26.271Z"),
		ID:      "k_test_task",
		Name:    "Test Task",
		Open:    true,
		Owner:   pointer.Pointer("u_bob_analyst"),
		Ticket:  "test-ticket",
		Updated: parseTime("2025-06-21T22:21:26.271Z"),
	})
	require.NoError(t, err, "failed to insert task")

	// Insert comments
	_, err = queries.InsertComment(ctx, sqlc.InsertCommentParams{
		Author:  "u_bob_analyst",
		Created: parseTime("2025-06-21T22:21:26.271Z"),
		ID:      "c_test_comment",
		Message: "Initial comment on the test ticket.",
		Ticket:  "test-ticket",
		Updated: parseTime("2025-06-21T22:21:26.271Z"),
	})
	require.NoError(t, err, "failed to insert comment")

	// Insert timeline
	_, err = queries.InsertTimeline(ctx, sqlc.InsertTimelineParams{
		Created: parseTime("2025-06-21T22:21:26.271Z"),
		ID:      "h_test_timeline",
		Message: "Initial timeline entry.",
		Ticket:  "test-ticket",
		Time:    parseTime("2023-01-01T00:00:00Z"),
		Updated: parseTime("2025-06-21T22:21:26.271Z"),
	})
	require.NoError(t, err, "failed to insert timeline entry")

	// Insert links
	_, err = queries.InsertLink(ctx, sqlc.InsertLinkParams{
		Created: parseTime("2025-06-21T22:21:26.271Z"),
		ID:      "l_test_link",
		Name:    "Catalyst",
		Ticket:  "test-ticket",
		Updated: parseTime("2025-06-21T22:21:26.271Z"),
		Url:     "https://example.com",
	})
	require.NoError(t, err, "failed to insert link")

	// Insert files
	_, err = queries.InsertFile(ctx, sqlc.InsertFileParams{
		Created: parseTime("2025-06-21T22:21:26.271Z"),
		ID:      "b_test_file",
		Name:    "hello.txt",
		Size:    5,
		Ticket:  "test-ticket",
		Updated: parseTime("2025-06-21T22:21:26.271Z"),
		Blob:    "hello_a20DUE9c77rj.txt",
	})
	require.NoError(t, err, "failed to insert file")

	// Insert features
	_, err = queries.CreateFeature(ctx, "dev")
	require.NoError(t, err, "failed to insert feature 'dev'")

	// Insert reactions
	_, err = queries.InsertReaction(ctx, sqlc.InsertReactionParams{
		ID:          "r-test-webhook",
		Name:        "Reaction",
		Action:      "python",
		Actiondata:  []byte(`{"requirements":"requests","script":"print('Hello, World!')"}`),
		Trigger:     "webhook",
		Triggerdata: []byte(`{"token":"1234567890","path":"test"}`),
		Created:     parseTime("2025-06-21T22:21:26.271Z"),
		Updated:     parseTime("2025-06-21T22:21:26.271Z"),
	})
	require.NoError(t, err, "failed to insert reaction")

	_, err = queries.InsertReaction(ctx, sqlc.InsertReactionParams{
		ID:          "r-test-proxy",
		Action:      "webhook",
		Name:        "Reaction",
		Actiondata:  []byte(`{"headers":{"Content-Type":"application/json"},"url":"http://127.0.0.1:12345/webhook"}`),
		Trigger:     "webhook",
		Triggerdata: []byte(`{"path":"test2"}`),
		Created:     parseTime("2025-06-21T22:21:26.271Z"),
		Updated:     parseTime("2025-06-21T22:21:26.271Z"),
	})
	require.NoError(t, err, "failed to insert reaction")

	_, err = queries.InsertReaction(ctx, sqlc.InsertReactionParams{
		ID:          "r-test-hook",
		Name:        "Hook",
		Action:      "python",
		Actiondata:  []byte(`{"requirements":"requests","script":"import requests\nrequests.post('http://127.0.0.1:12346/test', json={'test':True})"}`),
		Trigger:     "hook",
		Triggerdata: json.RawMessage([]byte(`{"collections":["tickets"],"events":["create"]}`)),
		Created:     parseTime("2025-06-21T22:21:26.271Z"),
		Updated:     parseTime("2025-06-21T22:21:26.271Z"),
	})
	require.NoError(t, err, "failed to insert reaction")

	// Insert user_groups
	err = queries.AssignGroupToUser(ctx, sqlc.AssignGroupToUserParams{
		UserID:  "u_bob_analyst",
		GroupID: "analyst",
	})
	require.NoError(t, err, "failed to assign analyst group to user")

	err = queries.AssignGroupToUser(ctx, sqlc.AssignGroupToUserParams{
		UserID:  "u_admin",
		GroupID: "admin",
	})
	require.NoError(t, err, "failed to assign admin group to user")

	files, err := queries.ListFiles(t.Context(), sqlc.ListFilesParams{
		Limit: 1000, // TODO
	})
	require.NoError(t, err, "failed to list files")

	for _, file := range files {
		_ = os.MkdirAll(path.Join(dir, "uploads", file.ID), 0o755)

		infoFilePath := path.Join(dir, "uploads", file.ID+".info")
		slog.InfoContext(t.Context(), "Creating file info", "path", infoFilePath)

		err = os.WriteFile(infoFilePath, []byte(`{"MetaData":{"filetype":"text/plain"}}`), 0o600)
		require.NoError(t, err, "failed to write file info")

		err = os.WriteFile(path.Join(dir, "uploads", file.ID, file.Blob), []byte("hello"), 0o600)
		require.NoError(t, err, "failed to write file blob")
	}
}
