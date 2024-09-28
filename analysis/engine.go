package analysis

import (
	"context"

	"github.com/pocketbase/pocketbase/daos"

	"github.com/SecurityBrewery/catalyst-analysis/analysis"
)

type Engine struct {
	engine         *analysis.Engine
	configProvider *PocketbaseCredentialProvider
}

func NewEngine(ctx context.Context, dao *daos.Dao) *Engine {
	services := analysis.LoadPlugins(ctx, &PocketbaseCredentialProvider{dao: dao})

	services = append(services, &analysis.Service{
		ID:     "catalyst",
		Plugin: &Catalyst{Dao: dao},
	})

	return &Engine{
		engine:         analysis.NewEngine(services),
		configProvider: &PocketbaseCredentialProvider{dao: dao},
	}
}

func (e *Engine) SetDao(dao *daos.Dao) error {
	service, err := e.engine.Service("catalyst")
	if err != nil {
		return err
	}

	if catalyst, ok := service.Plugin.(*Catalyst); ok {
		catalyst.Dao = dao
	}

	e.configProvider.dao = dao

	return nil
}

func (e *Engine) Engine() *analysis.Engine {
	return e.engine
}
