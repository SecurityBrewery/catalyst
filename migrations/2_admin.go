package migrations

import (
	"os"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

func adminUp(db dbx.Builder) error {
	adminEmail := os.Getenv("CATALYST_ADMIN_EMAIL")
	adminPassword := os.Getenv("CATALYST_ADMIN_PASSWORD")

	if adminEmail == "" || adminPassword == "" {
		return nil
	}

	admin := &models.Admin{Email: adminEmail}
	if err := admin.SetPassword(adminPassword); err != nil {
		return err
	}

	return daos.New(db).SaveAdmin(admin)
}
