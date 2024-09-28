package analysis

import (
	"context"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"

	"github.com/SecurityBrewery/catalyst-analysis/plugin"
	"github.com/SecurityBrewery/catalyst/migrations"
)

var _ plugin.Suggestor = &User{}

type User struct {
	catalyst *Catalyst
}

func (u *User) Info() plugin.ResourceTypeInfo {
	return plugin.ResourceTypeInfo{
		ID:                 "user",
		Name:               "User",
		Attributes:         []string{"email"},
		EnrichmentPatterns: []string{},
	}
}

func (u *User) Resource(_ context.Context, id string) (*plugin.Resource, error) {
	record, err := u.catalyst.Dao.FindRecordById(migrations.UserCollectionName, id)
	if err != nil {
		return nil, err
	}

	return u.toEnrichment(record), nil
}

func (u *User) Suggest(_ context.Context, s string) []*plugin.Resource {
	var records []*models.Record

	var err error

	if s == "" {
		err = u.catalyst.Dao.RecordQuery(migrations.UserCollectionName).Limit(3).All(&records)
	} else {
		records, err = u.catalyst.Dao.FindRecordsByFilter(
			migrations.UserCollectionName,
			"name ~ {:name}",
			"-created",
			3,
			0,
			dbx.Params{"name": s},
		)
	}

	if err != nil {
		return nil
	}

	if len(records) == 0 {
		return nil
	}

	var suggestions []*plugin.Resource
	for _, record := range records {
		suggestions = append(suggestions, u.toEnrichment(record))
	}

	return suggestions
}

func (u *User) toEnrichment(record *models.Record) *plugin.Resource {
	return &plugin.Resource{
		Type: u.Info().ID,
		ID:   record.GetId(),
		Name: record.GetString("name"),
		Icon: "User",
		Attributes: []plugin.Attribute{
			{ID: "email", Icon: "Mail", Name: "Mail", Value: record.GetString("email")},
		},
	}
}
