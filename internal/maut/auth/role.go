package auth

import "golang.org/x/exp/slices"

const AdminRole = "admin"

type Role struct {
	Name        string
	Permissions []string
}

func (r Role) Contains(p string) bool {
	if r.Name == "admin" {
		return true
	}

	return slices.Contains(r.Permissions, p)
}
