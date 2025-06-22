package data

import (
	"context"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

var ticketCreated = time.Date(2025, 2, 1, 11, 29, 35, 0, time.UTC)

func generateDemoTickets(ctx context.Context, queries *sqlc.Queries, users []sqlc.User, types []sqlc.ListTypesRow, count int) error {
	created := time.Now()
	number := gofakeit.Number(200*count, 300*count)

	for range count {
		number -= gofakeit.Number(100, 200)
		created = created.Add(time.Duration(-gofakeit.Number(13, 37)) * time.Hour)

		ticketType := random(types)

		newTicket, err := queries.CreateTicket(ctx, sqlc.CreateTicketParams{
			ID:          database.GenerateID(ticketType.Singular),
			Name:        fakeTicketTitle(),
			Type:        ticketType.ID,
			Description: fakeTicketDescription(),
			Open:        gofakeit.Bool(),
			Schema:      marshal(map[string]any{"type": "object", "properties": map[string]any{"tlp": map[string]any{"title": "TLP", "type": "string"}}}),
			State:       marshal(map[string]any{"severity": "Medium"}),
			Owner:       random(users).ID,
			Created:     dateTime(created),
			Updated:     dateTime(created.Add(time.Minute * 5)),
		})
		if err != nil {
			return fmt.Errorf("failed to create ticket: %w", err)
		}

		if err := createDemoComments(ctx, queries, users, newTicket); err != nil {
			return fmt.Errorf("failed to create comments for ticket %s: %w", newTicket.ID, err)
		}

		if err := createDemoTimeline(ctx, queries, created, newTicket); err != nil {
			return fmt.Errorf("failed to create timeline for ticket %s: %w", newTicket.ID, err)
		}

		if err := createDemoTasks(ctx, queries, users, newTicket); err != nil {
			return fmt.Errorf("failed to create tasks for ticket %s: %w", newTicket.ID, err)
		}

		if err := createDemoLinks(ctx, queries, newTicket); err != nil {
			return fmt.Errorf("failed to create links for ticket %s: %w", newTicket.ID, err)
		}
	}

	return nil
}

func CreateUpgradeTestDataTickets() map[string]sqlc.Ticket {
	return map[string]sqlc.Ticket{
		"t_0": {
			ID:          "t_0",
			Created:     dateTime(ticketCreated),
			Updated:     dateTime(ticketCreated.Add(time.Minute * 5)),
			Name:        "phishing-123",
			Type:        "alert",
			Description: "Phishing email reported by several employees.",
			Open:        true,
			Schema:      `{"type":"object","properties":{"tlp":{"title":"TLP","type":"string"}}}`,
			State:       `{"severity":"Medium"}`,
			Owner:       "u_test",
		},
	}
}
