package analysis

import (
	"github.com/pocketbase/pocketbase/daos"

	"github.com/SecurityBrewery/catalyst-analysis/plugin"
)

var _ plugin.Plugin = &Catalyst{}

type Catalyst struct {
	Dao *daos.Dao
}

func (a *Catalyst) Info() plugin.Info {
	return plugin.Info{
		Name: "Catalyst",
		ResourceTypes: []plugin.ResourceType{
			&User{catalyst: a},
			&Other{},
			&Link{},
		},
	}
}
