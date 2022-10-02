package catalyst

import maut "github.com/jonas-plum/maut/auth"

var Admin = &maut.Role{
	Name: "admin",
	Permissions: append(engineer.Permissions,
		"backup:create",
		"backup:restore",
		"dashboard:write",
		"job:read",
		"job:write",
		"log:read",
		"settings:write",
		"ticket:delete",
		"user:write",
		"userdata:write",
	),
}

var engineer = &maut.Role{
	Name: "engineer",
	Permissions: append(analyst.Permissions,
		"automation:write",
		"playbook:write",
		"template:write",
		"tickettype:write",
	),
}

var analyst = &maut.Role{
	Name: "analyst",
	Permissions: []string{
		"automation:read",
		"currentuser:read",
		"currentuserdata:read",
		"currentuserdata:write",
		"dashboard:read",
		"file:read",
		"file:write",
		"playbook:read",
		"settings:read",
		"template:read",
		"ticket:read",
		"ticket:write",
		"tickettype:read",
		"user:read",
		"userdata:read",
	},
}
