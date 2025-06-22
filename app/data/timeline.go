package data

import (
	"context"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

func createDemoTimeline(ctx context.Context, queries *sqlc.Queries, created time.Time, record sqlc.Ticket) error {
	for range gofakeit.IntN(5) {
		_, err := queries.CreateTimeline(ctx, sqlc.CreateTimelineParams{
			ID:      "test-" + gofakeit.UUID(),
			Message: fakeTicketTimelineMessage(),
			Ticket:  record.ID,
			Time:    created.Format("2006-01-02T15:04:05Z"),
			Created: dateTime(gofakeit.PastDate()),
			Updated: dateTime(gofakeit.PastDate()),
		})
		if err != nil {
			return fmt.Errorf("failed to create timeline for ticket %s: %w", record.ID, err)
		}
	}

	return nil
}

func CreateUpgradeTestDataTimeline() map[string]sqlc.Timeline {
	return map[string]sqlc.Timeline{
		"tl_0": {
			ID:      "tl_0",
			Created: dateTime(ticketCreated.Add(time.Minute * 15)),
			Updated: dateTime(ticketCreated.Add(time.Minute * 20)),
			Ticket:  "t_0",
			Time:    dateTime(ticketCreated.Add(time.Minute * 15)),
			Message: "This is a test timeline message.",
		},
	}
}
