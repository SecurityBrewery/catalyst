package migrations

import (
	"github.com/pocketbase/pocketbase/migrations"

	"github.com/SecurityBrewery/catalyst/ff"
)

func Register() {
	migrations.Register(baseUp, baseDown, "1700000000_base.go")
	migrations.Register(collectionsUp, collectionsDown, "1700000001_collections.go")
	migrations.Register(adminUp, nil, "1700000002_admin.go")

	if ff.HasDummyDataFlag() {
		migrations.Register(testdataUp, nil, "1700000003_dummydata.go")
	}

	migrations.Register(viewsUp, viewsDown, "1700000002_stats_view.go")
}
