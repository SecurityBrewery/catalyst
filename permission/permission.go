package permission

import (
	"context"
	"encoding/json"
	"log/slog"
)

type Permission = string

const (
	TicketRead    Permission = "ticket:read"
	TicketWrite   Permission = "ticket:write"
	FileRead      Permission = "file:read"
	FileWrite     Permission = "file:write"
	TypeRead      Permission = "type:read"
	TypeWrite     Permission = "type:write"
	UserRead      Permission = "user:read"
	UserWrite     Permission = "user:write"
	RoleRead      Permission = "role:read"
	RoleWrite     Permission = "role:write"
	ReactionRead  Permission = "reaction:read"
	ReactionWrite Permission = "reaction:write"
	WebhookRead   Permission = "webhook:read"
	WebhookWrite  Permission = "webhook:write"
)

func Default() []Permission {
	return []Permission{
		TypeRead,
		FileRead,
		TicketRead,
		TicketWrite,
	}
}

func AllPermissions() []Permission {
	return []Permission{
		TicketRead,
		TicketWrite,
		TypeRead,
		TypeWrite,
		UserRead,
		UserWrite,
		RoleRead,
		RoleWrite,
		ReactionRead,
		ReactionWrite,
		FileRead,
		FileWrite,
		WebhookRead,
		WebhookWrite,
	}
}

func FromJSONArray(ctx context.Context, permissions string) []string {
	var result []string
	if err := json.Unmarshal([]byte(permissions), &result); err != nil {
		slog.ErrorContext(ctx, "Failed to unmarshal permissions", "error", err)

		return nil
	}

	return result
}

func ToJSONArray(ctx context.Context, permissions []string) string {
	if len(permissions) == 0 {
		return "[]"
	}

	data, err := json.Marshal(permissions)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to marshal permissions", "error", err)

		return "[]"
	}

	return string(data)
}
