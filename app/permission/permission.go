package permission

import (
	"context"
	"encoding/json"
	"log/slog"
)

type Table struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

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

	TicketsTable   = Table{ID: "tickets", Name: "Tickets"}
	CommentsTable  = Table{ID: "comments", Name: "Comments"}
	LinksTable     = Table{ID: "links", Name: "Links"}
	TasksTable     = Table{ID: "tasks", Name: "Tasks"}
	TimelinesTable = Table{ID: "timeline", Name: "Timeline"}
	FilesTable     = Table{ID: "files", Name: "Files"}
	TypesTable     = Table{ID: "types", Name: "Types"}
	UsersTable     = Table{ID: "users", Name: "Users"}
	GroupsTable    = Table{ID: "groups", Name: "Groups"}
	ReactionsTable = Table{ID: "reactions", Name: "Reactions"}
	WebhooksTable  = Table{ID: "webhooks", Name: "Webhooks"}

	DashboardCountsTable = Table{ID: "dashboard_counts", Name: "Dashboard Counts"}
	SidebarTable         = Table{ID: "sidebar", Name: "Sidebar"}
	UserPermissionTable  = Table{ID: "user_permissions", Name: "User Permissions"}
	UserGroupTable       = Table{ID: "user_groups", Name: "User Groups"}
	GroupUserTable       = Table{ID: "group_users", Name: "Group Users"}
	GroupPermissionTable = Table{ID: "group_permissions", Name: "Group Permissions"}
	GroupParentTable     = Table{ID: "group_parents", Name: "Group Parents"}
	GroupChildTable      = Table{ID: "group_children", Name: "Group Children"}

	CreateAction = "create"
	UpdateAction = "update"
	DeleteAction = "delete"
)

func Tables() []Table {
	return []Table{
		TicketsTable,
		FilesTable,
		TypesTable,
		UsersTable,
		GroupsTable,
		ReactionsTable,
		WebhooksTable,
	}
}

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
