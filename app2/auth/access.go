package auth

import (
	"fmt"

	"github.com/SecurityBrewery/catalyst/app2/database/sqlc"
)

func HasAccess(user sqlc.User, permission string) error {
	if !user.Verified {
		return fmt.Errorf("user is not verified")
	}

	/*
		var hasAccess bool

		for _, p := range gjson.ParseBytes(user.Permissions).Array() {
			if p.Get("value").String() == permission {
				hasAccess = true

				break
			}
		}

		if !hasAccess {
			return fmt.Errorf("missing permission: %q", permission)
		}
	*/

	return nil
}
