package data

import (
	"context"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

func createDemoComments(ctx context.Context, queries *sqlc.Queries, users []sqlc.User, record sqlc.Ticket) error {
	for range gofakeit.IntN(5) {
		_, err := queries.CreateComment(ctx, sqlc.CreateCommentParams{
			ID:      "test-" + gofakeit.UUID(),
			Ticket:  record.ID,
			Author:  random(users).ID,
			Message: fakeTicketComment(),
			Created: dateTime(gofakeit.PastDate()),
			Updated: dateTime(gofakeit.PastDate()),
		})
		if err != nil {
			return fmt.Errorf("failed to create comment for ticket %s: %w", record.ID, err)
		}
	}

	return nil
}

func CreateUpgradeTestDataComments() map[string]sqlc.Comment {
	return map[string]sqlc.Comment{
		"c_0": {
			ID:      "c_0",
			Created: dateTime(ticketCreated.Add(time.Minute * 10)),
			Updated: dateTime(ticketCreated.Add(time.Minute * 15)),
			Ticket:  "t_0",
			Author:  "u_test",
			Message: "This is a test comment.",
		},
	}
}
