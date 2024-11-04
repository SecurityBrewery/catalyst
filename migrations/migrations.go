package migrations

import (
	"github.com/pocketbase/pocketbase/migrations"
)

func Register() {
	migrations.Register(baseUp, baseDown, "1700000000_base.go")
	migrations.Register(collectionsUp, collectionsDown, "1700000001_collections.go")
	migrations.Register(defaultDataUp, nil, "1700000003_defaultdata.go")
	migrations.Register(viewsUp, viewsDown, "1700000004_views.go")
	migrations.Register(reactionsUp, reactionsDown, "1700000005_reactions.go")
	migrations.Register(systemuserUp, systemuserDown, "1700000006_systemuser.go")
	migrations.Register(searchViewUp, searchViewDown, "1700000007_search_view.go")
	migrations.Register(dashboardCountsViewUpdateUp, dashboardCountsViewUpdateDown, "1700000008_dashboardview.go")
	migrations.Register(reactionsUpdateUp, nil, "1700000009_reactions_update.go")
}
