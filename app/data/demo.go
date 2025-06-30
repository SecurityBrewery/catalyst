package data

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/SecurityBrewery/catalyst/app/auth"
	"github.com/SecurityBrewery/catalyst/app/auth/password"
	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/pointer"
)

const (
	minimumUserCount   = 1
	minimumTicketCount = 1
)

var (
	//go:embed scripts/createticket.py
	createTicketPy string
	//go:embed scripts/alertingest.py
	alertIngestPy string
	//go:embed scripts/assigntickets.py
	assignTicketsPy string
)

func GenerateDemoData(ctx context.Context, queries *sqlc.Queries, userCount, ticketCount int) error {
	if userCount < minimumUserCount {
		userCount = minimumUserCount
	}

	if ticketCount < minimumTicketCount {
		ticketCount = minimumTicketCount
	}

	types, err := database.PaginateItems(ctx, func(ctx context.Context, offset, limit int64) ([]sqlc.ListTypesRow, error) {
		return queries.ListTypes(ctx, sqlc.ListTypesParams{Limit: limit, Offset: offset})
	})
	if err != nil {
		return fmt.Errorf("failed to list types: %w", err)
	}

	users, err := generateDemoUsers(ctx, queries, userCount, ticketCount)
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

	if err := generateDemoReactions(ctx, queries, ticketCount); err != nil {
		return fmt.Errorf("failed to create reaction records: %w", err)
	}

	if err := generateDemoGroups(ctx, queries, users, ticketCount); err != nil {
		return fmt.Errorf("failed to create group records: %w", err)
	}

	return nil
}

func generateDemoUsers(ctx context.Context, queries *sqlc.Queries, count, ticketCount int) ([]sqlc.User, error) {
	users := make([]sqlc.User, 0, count)

	// create the test user
	user, err := queries.GetUser(ctx, "u_test")
	if err != nil {
		newUser, err := createTestUser(ctx, queries)
		if err != nil {
			return nil, err
		}

		users = append(users, newUser)
	} else {
		users = append(users, user)
	}

	for range count - 1 {
		newUser, err := createDemoUser(ctx, queries, ticketCount)
		if err != nil {
			return nil, err
		}

		users = append(users, newUser)
	}

	return users, nil
}

func createDemoUser(ctx context.Context, queries *sqlc.Queries, ticketCount int) (sqlc.User, error) {
	username := gofakeit.Username()

	passwordHash, tokenKey, err := password.Hash(gofakeit.Password(true, true, true, true, false, 16))
	if err != nil {
		return sqlc.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	created, updated := dates(ticketCount)

	return queries.InsertUser(ctx, sqlc.InsertUserParams{
		ID:           database.GenerateID("u"),
		Name:         pointer.Pointer(gofakeit.Name()),
		Email:        pointer.Pointer(username + "@catalyst-soar.com"),
		Username:     username,
		PasswordHash: passwordHash,
		TokenKey:     tokenKey,
		Active:       gofakeit.Bool(),
		Created:      created,
		Updated:      updated,
	})
}

var ticketCreated = time.Date(2025, 2, 1, 11, 29, 35, 0, time.UTC)

func generateDemoTickets(ctx context.Context, queries *sqlc.Queries, users []sqlc.User, types []sqlc.ListTypesRow, count int) error { //nolint:cyclop
	for range count {
		newTicket, err := createDemoTicket(ctx, queries, random(types), random(users).ID, fakeTicketTitle(), fakeTicketDescription(), count)
		if err != nil {
			return fmt.Errorf("failed to create ticket: %w", err)
		}

		for range gofakeit.IntN(5) {
			_, err := createDemoComment(ctx, queries, newTicket.ID, random(users).ID, fakeTicketComment(), count)
			if err != nil {
				return fmt.Errorf("failed to create comment for ticket %s: %w", newTicket.ID, err)
			}
		}

		for range gofakeit.IntN(5) {
			_, err := createDemoTimeline(ctx, queries, newTicket.ID, fakeTicketTimelineMessage(), count)
			if err != nil {
				return fmt.Errorf("failed to create timeline for ticket %s: %w", newTicket.ID, err)
			}
		}

		for range gofakeit.IntN(5) {
			_, err := createDemoTask(ctx, queries, newTicket.ID, random(users).ID, fakeTicketTask(), count)
			if err != nil {
				return fmt.Errorf("failed to create task for ticket %s: %w", newTicket.ID, err)
			}
		}

		for range gofakeit.IntN(5) {
			_, err := createDemoLink(ctx, queries, newTicket.ID, random([]string{"Blog", "Forum", "Wiki", "Documentation"}), gofakeit.URL(), count)
			if err != nil {
				return fmt.Errorf("failed to create link for ticket %s: %w", newTicket.ID, err)
			}
		}
	}

	return nil
}

func createDemoTicket(ctx context.Context, queries *sqlc.Queries, ticketType sqlc.ListTypesRow, userID, name, description string, ticketCount int) (sqlc.Ticket, error) {
	created, updated := dates(ticketCount)

	ticket, err := queries.InsertTicket(
		ctx,
		sqlc.InsertTicketParams{
			ID:          database.GenerateID(ticketType.Singular),
			Name:        name,
			Description: description,
			Open:        gofakeit.Bool(),
			Owner:       &userID,
			Schema:      marshal(map[string]any{"type": "object", "properties": map[string]any{"tlp": map[string]any{"title": "TLP", "type": "string"}}}),
			State:       marshal(map[string]any{"severity": "Medium"}),
			Type:        ticketType.ID,
			Created:     created,
			Updated:     updated,
		},
	)
	if err != nil {
		return sqlc.Ticket{}, fmt.Errorf("failed to create ticket for user %s: %w", userID, err)
	}

	return ticket, nil
}

func createDemoComment(ctx context.Context, queries *sqlc.Queries, ticketID, userID, message string, ticketCount int) (*sqlc.Comment, error) {
	created, updated := dates(ticketCount)

	comment, err := queries.InsertComment(ctx, sqlc.InsertCommentParams{
		ID:      database.GenerateID("c"),
		Ticket:  ticketID,
		Author:  userID,
		Message: message,
		Created: created,
		Updated: updated,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create comment for ticket %s: %w", ticketID, err)
	}

	return &comment, nil
}

func createDemoTimeline(ctx context.Context, queries *sqlc.Queries, ticketID, message string, ticketCount int) (*sqlc.Timeline, error) {
	created, updated := dates(ticketCount)

	timeline, err := queries.InsertTimeline(ctx, sqlc.InsertTimelineParams{
		ID:      database.GenerateID("tl"),
		Ticket:  ticketID,
		Message: message,
		Time:    ticketCreated,
		Created: created,
		Updated: updated,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create timeline for ticket %s: %w", ticketID, err)
	}

	return &timeline, nil
}

func createDemoTask(ctx context.Context, queries *sqlc.Queries, ticketID, userID, name string, ticketCount int) (*sqlc.Task, error) {
	created, updated := dates(ticketCount)

	task, err := queries.InsertTask(ctx, sqlc.InsertTaskParams{
		ID:      database.GenerateID("t"),
		Ticket:  ticketID,
		Owner:   &userID,
		Name:    name,
		Open:    gofakeit.Bool(),
		Created: created,
		Updated: updated,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create task for ticket %s: %w", ticketID, err)
	}

	return &task, nil
}

func createDemoLink(ctx context.Context, queries *sqlc.Queries, ticketID, name, url string, ticketCount int) (*sqlc.Link, error) {
	created, updated := dates(ticketCount)

	link, err := queries.InsertLink(ctx, sqlc.InsertLinkParams{
		ID:      database.GenerateID("l"),
		Ticket:  ticketID,
		Name:    name,
		Url:     url,
		Created: created,
		Updated: updated,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create link for ticket %s: %w", ticketID, err)
	}

	return &link, nil
}

func generateDemoReactions(ctx context.Context, queries *sqlc.Queries, ticketCount int) error {
	created, updated := dates(ticketCount)

	_, err := queries.InsertReaction(ctx, sqlc.InsertReactionParams{
		ID:          "r-schedule",
		Name:        "Create New Ticket",
		Trigger:     "schedule",
		Triggerdata: marshal(map[string]any{"expression": "12 * * * *"}),
		Action:      "python",
		Actiondata: marshal(map[string]any{
			"requirements": "requests",
			"script":       createTicketPy,
		}),
		Created: created,
		Updated: updated,
	})
	if err != nil {
		return fmt.Errorf("failed to create reaction for schedule trigger: %w", err)
	}

	created, updated = dates(ticketCount)

	_, err = queries.InsertReaction(ctx, sqlc.InsertReactionParams{
		ID:          "r-webhook",
		Name:        "Alert Ingest Webhook",
		Trigger:     "webhook",
		Triggerdata: marshal(map[string]any{"token": "1234567890", "path": "webhook"}),
		Action:      "python",
		Actiondata: marshal(map[string]any{
			"requirements": "requests",
			"script":       alertIngestPy,
		}),
		Created: created,
		Updated: updated,
	})
	if err != nil {
		return fmt.Errorf("failed to create reaction for webhook trigger: %w", err)
	}

	created, updated = dates(ticketCount)

	_, err = queries.InsertReaction(ctx, sqlc.InsertReactionParams{
		ID:          "r-hook",
		Name:        "Assign new Tickets",
		Trigger:     "hook",
		Triggerdata: marshal(map[string]any{"collections": []any{"tickets"}, "events": []any{"create"}}),
		Action:      "python",
		Actiondata: marshal(map[string]any{
			"requirements": "requests",
			"script":       assignTicketsPy,
		}),
		Created: created,
		Updated: updated,
	})
	if err != nil {
		return fmt.Errorf("failed to create reaction for hook trigger: %w", err)
	}

	return nil
}

func generateDemoGroups(ctx context.Context, queries *sqlc.Queries, users []sqlc.User, ticketCount int) error { //nolint:cyclop
	created, updated := dates(ticketCount)

	_, err := queries.InsertGroup(ctx, sqlc.InsertGroupParams{
		ID:          "team-ir",
		Name:        "IR Team",
		Permissions: auth.ToJSONArray(ctx, []string{}),
		Created:     created,
		Updated:     updated,
	})
	if err != nil {
		return fmt.Errorf("failed to create IR team group: %w", err)
	}

	created, updated = dates(ticketCount)

	_, err = queries.InsertGroup(ctx, sqlc.InsertGroupParams{
		ID:          "team-seceng",
		Name:        "Security Engineering Team",
		Permissions: auth.ToJSONArray(ctx, []string{}),
		Created:     created,
		Updated:     updated,
	})
	if err != nil {
		return fmt.Errorf("failed to create IR team group: %w", err)
	}

	created, updated = dates(ticketCount)

	_, err = queries.InsertGroup(ctx, sqlc.InsertGroupParams{
		ID:          "team-security",
		Name:        "Security Team",
		Permissions: auth.ToJSONArray(ctx, []string{}),
		Created:     created,
		Updated:     updated,
	})
	if err != nil {
		return fmt.Errorf("failed to create security team group: %w", err)
	}

	created, updated = dates(ticketCount)

	_, err = queries.InsertGroup(ctx, sqlc.InsertGroupParams{
		ID:          "g-engineer",
		Name:        "Engineer",
		Permissions: auth.ToJSONArray(ctx, []string{"reaction:read", "reaction:write"}),
		Created:     created,
		Updated:     updated,
	})
	if err != nil {
		return fmt.Errorf("failed to create analyst group: %w", err)
	}

	for _, user := range users {
		group := gofakeit.RandomString([]string{"team-seceng", "team-ir"})
		if user.ID == "u_test" {
			group = "admin"
		}

		if err := queries.AssignGroupToUser(ctx, sqlc.AssignGroupToUserParams{
			UserID:  user.ID,
			GroupID: group,
		}); err != nil {
			return fmt.Errorf("failed to assign group %s to user %s: %w", group, user.ID, err)
		}
	}

	err = queries.AssignParentGroup(ctx, sqlc.AssignParentGroupParams{
		ParentGroupID: "team-ir",
		ChildGroupID:  "analyst",
	})
	if err != nil {
		return fmt.Errorf("failed to assign parent group: %w", err)
	}

	err = queries.AssignParentGroup(ctx, sqlc.AssignParentGroupParams{
		ParentGroupID: "team-seceng",
		ChildGroupID:  "g-engineer",
	})
	if err != nil {
		return fmt.Errorf("failed to assign parent group: %w", err)
	}

	err = queries.AssignParentGroup(ctx, sqlc.AssignParentGroupParams{
		ParentGroupID: "team-ir",
		ChildGroupID:  "team-security",
	})
	if err != nil {
		return fmt.Errorf("failed to assign parent group: %w", err)
	}

	err = queries.AssignParentGroup(ctx, sqlc.AssignParentGroupParams{
		ParentGroupID: "team-seceng",
		ChildGroupID:  "team-security",
	})
	if err != nil {
		return fmt.Errorf("failed to assign parent group: %w", err)
	}

	return nil
}

func weeksAgo(c int) time.Time {
	return time.Now().UTC().AddDate(0, 0, -7*c)
}

func dates(ticketCount int) (time.Time, time.Time) {
	const ticketsPerWeek = 10
	weeks := ticketCount / ticketsPerWeek

	created := gofakeit.DateRange(weeksAgo(1), weeksAgo(weeks+1)).UTC()
	updated := gofakeit.DateRange(created, time.Now()).UTC()

	return created, updated
}
