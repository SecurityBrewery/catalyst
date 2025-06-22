package data

import (
	"context"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/permission"
)

func generateDemoGroups(ctx context.Context, queries *sqlc.Queries, users []sqlc.User) error { //nolint:cyclop
	_, err := queries.CreateGroup(ctx, sqlc.CreateGroupParams{
		ID:          "team-ir",
		Name:        "IR Team",
		Permissions: permission.ToJSONArray(ctx, []string{}),
		Created:     dateTime(gofakeit.PastDate()),
		Updated:     dateTime(gofakeit.PastDate()),
	})
	if err != nil {
		return fmt.Errorf("failed to create IR team group: %w", err)
	}

	_, err = queries.CreateGroup(ctx, sqlc.CreateGroupParams{
		ID:          "team-seceng",
		Name:        "Security Engineering Team",
		Permissions: permission.ToJSONArray(ctx, []string{}),
		Created:     dateTime(gofakeit.PastDate()),
		Updated:     dateTime(gofakeit.PastDate()),
	})
	if err != nil {
		return fmt.Errorf("failed to create IR team group: %w", err)
	}

	_, err = queries.CreateGroup(ctx, sqlc.CreateGroupParams{
		ID:          "team-security",
		Name:        "Security Team",
		Permissions: permission.ToJSONArray(ctx, []string{}),
		Created:     dateTime(gofakeit.PastDate()),
		Updated:     dateTime(gofakeit.PastDate()),
	})
	if err != nil {
		return fmt.Errorf("failed to create security team group: %w", err)
	}

	_, err = queries.CreateGroup(ctx, sqlc.CreateGroupParams{
		ID:          "g-engineer",
		Name:        "Engineer",
		Permissions: permission.ToJSONArray(ctx, []string{"reaction:read", "reaction:write"}),
		Created:     dateTime(gofakeit.PastDate()),
		Updated:     dateTime(gofakeit.PastDate()),
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
