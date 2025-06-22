package data

import (
	"context"
	"fmt"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

const (
	minimumUserCount   = 1
	minimumTicketCount = 1
)

func GenerateDemoData(ctx context.Context, queries *sqlc.Queries, userCount, ticketCount int) error {
	if userCount < minimumUserCount {
		userCount = minimumUserCount
	}

	if ticketCount < minimumTicketCount {
		ticketCount = minimumTicketCount
	}

	types, err := queries.ListTypes(ctx, sqlc.ListTypesParams{
		Limit:  100,
		Offset: 0,
	})
	if err != nil {
		return fmt.Errorf("failed to list types: %w", err)
	}

	users, err := generateDemoUsers(ctx, queries, userCount)
	if err != nil {
		return fmt.Errorf("failed to create user records: %w", err)
	}

	if len(types) == 0 {
		return fmt.Errorf("no types found")
	}

	if len(users) == 0 {
		return fmt.Errorf("no users found")
	}

	if err := generateDemoTickets(ctx, queries, users, types, ticketCount); err != nil {
		return fmt.Errorf("failed to create ticket records: %w", err)
	}

	if err := generateDemoReactions(ctx, queries); err != nil {
		return fmt.Errorf("failed to create reaction records: %w", err)
	}

	if err := generateDemoGroups(ctx, queries, users); err != nil {
		return fmt.Errorf("failed to create group records: %w", err)
	}

	return nil
}
