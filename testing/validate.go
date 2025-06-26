package testing

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"golang.org/x/crypto/bcrypt"

	"github.com/SecurityBrewery/catalyst/app/data"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

func ValidateUpgradeTestData(t *testing.T, queries *sqlc.Queries) {
	t.Helper()

	// users
	userRecord, err := queries.UserByUserName(t.Context(), "u_test")
	require.NoError(t, err, "failed to find user by username")

	assert.Equal(t, "u_test", userRecord.ID, "user ID does not match expected value")
	assert.Equal(t, "u_test", userRecord.Username, "username does not match expected value")
	assert.Equal(t, "Test User", *userRecord.Name, "name does not match expected value")
	assert.Equal(t, "user@catalyst-soar.com", *userRecord.Email, "email does not match expected value")
	assert.True(t, userRecord.Active, "user should be verified")
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(userRecord.Passwordhash), []byte("1234567890")), "password hash does not match expected value")

	for id := range data.CreateUpgradeTestDataTickets() {
		ticket, err := queries.Ticket(t.Context(), id)
		require.NoError(t, err, "failed to find ticket")

		assert.Equal(t, "phishing-123", ticket.Name)
		assert.Equal(t, "Phishing email reported by several employees.", ticket.Description)
		assert.True(t, ticket.Open)
		assert.Equal(t, "alert", ticket.Type)
		assert.Equal(t, "u_test", *ticket.Owner)
		assert.JSONEq(t, `{"type":"object","properties":{"tlp":{"title":"TLP","type":"string"}}}`, string(ticket.Schema))
		assert.JSONEq(t, `{"severity":"Medium"}`, string(ticket.State))
	}

	for id := range data.CreateUpgradeTestDataComments() {
		comment, err := queries.GetComment(t.Context(), id)
		require.NoError(t, err, "failed to find comment")

		assert.Equal(t, "This is a test comment.", comment.Message)
	}

	for id := range data.CreateUpgradeTestDataTimeline() {
		timeline, err := queries.GetTimeline(t.Context(), id)
		require.NoError(t, err, "failed to find timeline")

		assert.Equal(t, "This is a test timeline message.", timeline.Message)
	}

	for id := range data.CreateUpgradeTestDataTasks() {
		task, err := queries.GetTask(t.Context(), id)
		require.NoError(t, err, "failed to find task")

		assert.Equal(t, "This is a test task.", task.Name)
	}

	for id := range data.CreateUpgradeTestDataLinks() {
		link, err := queries.GetLink(t.Context(), id)
		require.NoError(t, err, "failed to find link")

		assert.Equal(t, "https://www.example.com", link.Url)
		assert.Equal(t, "This is a test link.", link.Name)
	}

	for id := range data.CreateUpgradeTestDataReaction() {
		reaction, err := queries.GetReaction(t.Context(), id)
		require.NoError(t, err, "failed to find reaction")

		assert.Equal(t, "Create New Ticket", reaction.Name)
		assert.Equal(t, "schedule", reaction.Trigger)
		assert.JSONEq(t, "{\"expression\":\"12 * * * *\"}", string(reaction.Triggerdata))
		assert.Equal(t, "python", reaction.Action)
		assert.Equal(t, "pocketbase", gjson.GetBytes(reaction.Actiondata, "requirements").String())
		equalWithoutSpace(t, data.Script, gjson.GetBytes(reaction.Actiondata, "script").String())
	}
}

func equalWithoutSpace(t *testing.T, expected, actual string) {
	t.Helper()

	assert.Equal(t, removeAllWhitespace(expected), removeAllWhitespace(actual))
}

func removeAllWhitespace(s string) string {
	return strings.Map(func(r rune) rune {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			return -1 // remove whitespace characters
		}

		return r // keep other characters
	}, s)
}
