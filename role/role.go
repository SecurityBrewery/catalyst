package role

import (
	"errors"
	"sort"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/SecurityBrewery/catalyst/generated/model"
)

type Role string

const (
	Analyst  string = "analyst"
	Engineer string = "engineer"
	Admin    string = "admin"

	AutomationRead       Role = "analyst:automation:read"
	CurrentuserRead      Role = "analyst:currentuser:read"
	CurrentuserdataRead  Role = "analyst:currentuserdata:read"
	CurrentuserdataWrite Role = "analyst:currentsettings:write"
	DashboardRead        Role = "analyst:dashboard:read"
	FileReadWrite        Role = "analyst:file"
	GroupRead            Role = "analyst:group:read"
	PlaybookRead         Role = "analyst:playbook:read"
	RuleRead             Role = "analyst:rule:read"
	SettingsRead         Role = "analyst:settings:read"
	TemplateRead         Role = "analyst:template:read"
	TicketRead           Role = "analyst:ticket:read"
	TicketWrite          Role = "analyst:ticket:write"
	TickettypeRead       Role = "analyst:tickettype:read"
	UserRead             Role = "analyst:user:read"

	AutomationWrite Role = "engineer:automation:write"
	PlaybookWrite   Role = "engineer:playbook:write"
	RuleWrite       Role = "engineer:rule:write"
	TemplateWrite   Role = "engineer:template:write"
	TickettypeWrite Role = "engineer:tickettype:write"

	BackupRead     Role = "admin:backup:read"
	BackupRestore  Role = "admin:backup:restore"
	DashboardWrite Role = "admin:dashboard:write"
	GroupWrite     Role = "admin:group:write"
	JobRead        Role = "admin:job:read"
	JobWrite       Role = "admin:job:write"
	LogRead        Role = "admin:log:read"
	SettingsWrite  Role = "admin:settings:write"
	TicketDelete   Role = "admin:ticket:delete"
	UserWrite      Role = "admin:user:write"
	UserdataRead   Role = "admin:userdata:read"
	UserdataWrite  Role = "admin:userdata:write"
)

func (p Role) String() string {
	return string(p)
}

func UserHasRoles(user *model.UserResponse, roles []Role) bool {
	hasRoles := true
	for _, role := range roles {
		if !UserHasRole(user, role) {
			hasRoles = false

			break
		}
	}

	return hasRoles
}

func UserHasRole(user *model.UserResponse, role Role) bool {
	return slices.Contains(FromStrings(user.Roles), role)
}

func Explodes(s []string) []Role {
	var roles []Role
	for _, e := range s {
		roles = append(roles, Explode(e)...)
	}
	sort.Slice(roles, func(i, j int) bool {
		return roles[i].String() < roles[j].String()
	})
	roles = slices.Compact(roles)

	return roles
}

func Explode(s string) []Role {
	var roles []Role

	switch s {
	case Admin:
		roles = append(roles, listPrefix(Admin)...)

		fallthrough
	case Engineer:
		roles = append(roles, listPrefix(Engineer)...)

		fallthrough
	case Analyst:
		roles = append(roles, listPrefix(Analyst)...)

		return roles
	}

	for _, role := range List() {
		if role.String() == s {
			roles = append(roles, role)
		}
	}

	return roles
}

func listPrefix(s string) []Role {
	var roles []Role

	for _, role := range List() {
		if strings.HasPrefix(role.String(), s+":") {
			roles = append(roles, role)
		}
	}

	return roles
}

func List() []Role {
	return []Role{
		AutomationRead, CurrentuserdataRead, CurrentuserdataWrite,
		CurrentuserRead, FileReadWrite, GroupRead, PlaybookRead, RuleRead,
		UserdataRead, SettingsRead, TemplateRead, TicketRead, TickettypeRead,
		TicketWrite, UserRead, AutomationWrite, PlaybookWrite, RuleWrite,
		TemplateWrite, TickettypeWrite, BackupRead, BackupRestore, GroupWrite,
		LogRead, UserdataWrite, TicketDelete, UserWrite, JobRead, JobWrite,
		SettingsWrite, DashboardRead, DashboardWrite,
	}
}

func fromString(s string) (Role, error) {
	for _, role := range List() {
		if role.String() == s {
			return role, nil
		}
	}

	return "", errors.New("unknown role")
}

func Strings(roles []Role) []string {
	var s []string
	for _, role := range roles {
		s = append(s, role.String())
	}

	return s
}

func FromStrings(s []string) []Role {
	var roles []Role
	for _, e := range s {
		role, err := fromString(e)
		if err != nil {
			continue
		}
		roles = append(roles, role)
	}

	return roles
}
