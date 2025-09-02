package auth

import (
	"context"
	"encoding/json"
	"log/slog"
)

var (
	TicketReadPermission    = "ticket:read"
	TicketWritePermission   = "ticket:write"
	FileReadPermission      = "file:read"
	FileWritePermission     = "file:write"
	TypeReadPermission      = "type:read"
	TypeWritePermission     = "type:write"
	UserReadPermission      = "user:read"
	UserWritePermission     = "user:write"
	GroupReadPermission     = "group:read"
	GroupWritePermission    = "group:write"
	ReactionReadPermission  = "reaction:read"
	ReactionWritePermission = "reaction:write"
	WebhookReadPermission   = "webhook:read"
	WebhookWritePermission  = "webhook:write"
	SettingsReadPermission  = "settings:read"
	SettingsWritePermission = "settings:write"
)

func All() []string {
	return []string{
		TicketReadPermission,
		TicketWritePermission,
		FileReadPermission,
		FileWritePermission,
		TypeReadPermission,
		TypeWritePermission,
		UserReadPermission,
		UserWritePermission,
		GroupReadPermission,
		GroupWritePermission,
		ReactionReadPermission,
		ReactionWritePermission,
		WebhookReadPermission,
		WebhookWritePermission,
		SettingsReadPermission,
		SettingsWritePermission,
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
