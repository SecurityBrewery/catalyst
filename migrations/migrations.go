package migrations

import (
	"github.com/pocketbase/pocketbase/migrations"
)

func Register() {
	migrations.Register(baseUp, baseDown, "1700000000_base.go")
	migrations.Register(collectionsUp, collectionsDown, "1700000001_collections.go")
	migrations.Register(defaultDataUp, nil, "1700000003_defaultdata.go")
	migrations.Register(viewsUp, viewsDown, "1700000004_views.go")
	migrations.Register(playbookUp, playbookDown, "1700000005_playbooks.go")
}
