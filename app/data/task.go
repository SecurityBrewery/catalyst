package data

import (
	"context"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

func createDemoTasks(ctx context.Context, queries *sqlc.Queries, users []sqlc.User, record sqlc.Ticket) error {
	for range gofakeit.IntN(5) {
		_, err := queries.CreateTask(ctx, sqlc.CreateTaskParams{
			ID:      "test-" + gofakeit.UUID(),
			Name:    fakeTicketTask(),
			Open:    gofakeit.Bool(),
			Owner:   random(users).ID,
			Ticket:  record.ID,
			Created: dateTime(gofakeit.PastDate()),
			Updated: dateTime(gofakeit.PastDate()),
		})
		if err != nil {
			return fmt.Errorf("failed to create task for ticket %s: %w", record.ID, err)
		}
	}

	return nil
}

func CreateUpgradeTestDataTasks() map[string]sqlc.Task {
	return map[string]sqlc.Task{
		"ts_0": {
			ID:      "ts_0",
			Created: dateTime(ticketCreated.Add(time.Minute * 20)),
			Updated: dateTime(ticketCreated.Add(time.Minute * 25)),
			Ticket:  "t_0",
			Name:    "This is a test task.",
			Open:    true,
			Owner:   "u_test",
		},
	}
}
