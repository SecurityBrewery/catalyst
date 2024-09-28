package analysis

import (
	"context"
	"github.com/SecurityBrewery/catalyst-analysis/config"
	"github.com/SecurityBrewery/catalyst/migrations"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/tidwall/gjson"
)

type PocketbaseCredentialProvider struct {
	dao *daos.Dao
}

func (p *PocketbaseCredentialProvider) Config(_ context.Context) (*config.Config, error) {
	records, err := p.dao.FindRecordsByExpr(migrations.IntegrationCollectionName)
	if err != nil {
		return nil, err
	}

	services := make([]*config.ServiceConfig, 0, len(records))
	for _, record := range records {
		services = append(services, &config.ServiceConfig{
			ID:     record.GetString("name"),
			Plugin: record.GetString("plugin"),
		})
	}

	return &config.Config{
		Services: services,
	}, nil
}

func (p *PocketbaseCredentialProvider) Get(_ context.Context, service, key string) (string, bool, error) {
	records, err := p.dao.FindRecordsByExpr(migrations.IntegrationCollectionName, dbx.HashExp{"name": service})
	if err != nil {
		return "", false, err
	}

	if len(records) == 0 {
		return "", false, nil
	}

	config := records[0].GetString("config")
	if config == "" {
		return "", false, nil
	}

	return gjson.Get(config, key).String(), true, nil
}
