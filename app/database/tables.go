package database

type Table struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var (
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
