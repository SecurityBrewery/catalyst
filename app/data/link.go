package data

import (
	"context"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

func createDemoLinks(ctx context.Context, queries *sqlc.Queries, record sqlc.Ticket) error {
	for range gofakeit.IntN(5) {
		_, err := queries.CreateLink(ctx, sqlc.CreateLinkParams{
			ID:      "test-" + gofakeit.UUID(),
			Ticket:  record.ID,
			Url:     gofakeit.URL(),
			Name:    random([]string{"Blog", "Forum", "Wiki", "Documentation"}),
			Created: dateTime(gofakeit.PastDate()),
			Updated: dateTime(gofakeit.PastDate()),
		})
		if err != nil {
			return fmt.Errorf("failed to create link for ticket %s: %w", record.ID, err)
		}
	}

	return nil
}

func CreateUpgradeTestDataLinks() map[string]sqlc.Link {
	return map[string]sqlc.Link{
		"l_0": {
			ID:      "l_0",
			Created: dateTime(ticketCreated.Add(time.Minute * 25)),
			Updated: dateTime(ticketCreated.Add(time.Minute * 30)),
			Ticket:  "t_0",
			Url:     "https://www.example.com",
			Name:    "This is a test link.",
		},
	}
}
