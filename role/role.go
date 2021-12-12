package role

import (
	"errors"
	"sort"
	"strings"

	"github.com/SecurityBrewery/catalyst/generated/models"
)

type Role string

const (
	Analyst  string = "analyst"
	Engineer string = "engineer"
	Admin    string = "admin"

	AutomationRead       Role = "analyst:automation:read"
	CurrentuserdataRead  Role = "analyst:currentuserdata:read"
	CurrentuserdataWrite Role = "analyst:currentsettings:write"
	CurrentuserRead      Role = "analyst:currentuser:read"
	FileReadWrite        Role = "analyst:file"
	GroupRead            Role = "analyst:group:read"
	PlaybookRead         Role = "analyst:playbook:read"
	RuleRead             Role = "analyst:rule:read"
	SettingsRead         Role = "analyst:settings:read"
	TemplateRead         Role = "analyst:template:read"
	TicketRead           Role = "analyst:ticket:read"
	TickettypeRead       Role = "analyst:tickettype:read"
	TicketWrite          Role = "analyst:ticket:write"
	UserRead             Role = "analyst:user:read"

	AutomationWrite Role = "engineer:automation:write"
	PlaybookWrite   Role = "engineer:playbook:write"
	RuleWrite       Role = "engineer:rule:write"
	TemplateWrite   Role = "engineer:template:write"
	TickettypeWrite Role = "engineer:tickettype:write"

	BackupRead    Role = "admin:backup:read"
	BackupRestore Role = "admin:backup:restore"
	GroupWrite    Role = "admin:group:write"
	JobWrite      Role = "admin:job:write"
	JobRead       Role = "admin:job:read"
	LogRead       Role = "admin:log:read"
	UserdataRead  Role = "admin:userdata:read"
	UserdataWrite Role = "admin:userdata:write"
	TicketDelete  Role = "admin:ticket:delete"
	UserWrite     Role = "admin:user:write"
)

func (p Role) String() string {
	return string(p)
}

func UserHasRoles(user *models.UserResponse, roles []Role) bool {
	hasRoles := true
	for _, role := range roles {
		if !UserHasRole(user, role) {
			hasRoles = false
			break
		}
	}
	return hasRoles
}

func UserHasRole(user *models.UserResponse, role Role) bool {
	return ContainsRole(FromStrings(user.Roles), role)
}

func ContainsRole(roles []Role, role Role) bool {
	for _, r := range roles {
		if r.String() == role.String() { //  || strings.HasPrefix(role.String(), r.String()+":")
			return true
		}
	}
	return false
}

func Explodes(s []string) []Role {
	var roles []Role
	for _, e := range s {
		roles = append(roles, Explode(e)...)
	}
	roles = unique(roles)
	sort.Slice(roles, func(i, j int) bool {
		return roles[i].String() < roles[j].String()
	})

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

func unique(l []Role) []Role {
	keys := make(map[Role]bool)
	var list []Role
	for _, entry := range l {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func List() []Role {
	return []Role{
		AutomationRead, CurrentuserdataRead, CurrentuserdataWrite,
		CurrentuserRead, FileReadWrite, GroupRead, PlaybookRead, RuleRead,
		UserdataRead, SettingsRead, TemplateRead, TicketRead, TickettypeRead,
		TicketWrite, UserRead, AutomationWrite, PlaybookWrite, RuleWrite,
		TemplateWrite, TickettypeWrite, BackupRead, BackupRestore, GroupWrite,
		LogRead, UserdataWrite, TicketDelete, UserWrite, JobRead, JobWrite,
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
