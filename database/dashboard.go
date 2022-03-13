package database

import (
	"context"
	"errors"

	"github.com/arangodb/go-driver"
	"github.com/iancoleman/strcase"

	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/model"
)

func toDashboardResponse(key string, doc *model.Dashboard) *model.DashboardResponse {
	return &model.DashboardResponse{
		ID:      key,
		Name:    doc.Name,
		Widgets: doc.Widgets,
	}
}

func (db *Database) DashboardCreate(ctx context.Context, dashboard *model.Dashboard) (*model.DashboardResponse, error) {
	if dashboard == nil {
		return nil, errors.New("requires dashboard")
	}
	if dashboard.Name == "" {
		return nil, errors.New("requires dashboard name")
	}

	var doc model.Dashboard
	newctx := driver.WithReturnNew(ctx, &doc)

	meta, err := db.dashboardCollection.CreateDocument(ctx, newctx, strcase.ToKebab(dashboard.Name), dashboard)
	if err != nil {
		return nil, err
	}

	return toDashboardResponse(meta.Key, &doc), nil
}

func (db *Database) DashboardGet(ctx context.Context, id string) (*model.DashboardResponse, error) {
	var doc model.Dashboard
	meta, err := db.dashboardCollection.ReadDocument(ctx, id, &doc)
	if err != nil {
		return nil, err
	}

	return toDashboardResponse(meta.Key, &doc), nil
}

func (db *Database) DashboardUpdate(ctx context.Context, id string, dashboard *model.Dashboard) (*model.DashboardResponse, error) {
	var doc model.Dashboard
	ctx = driver.WithReturnNew(ctx, &doc)

	meta, err := db.dashboardCollection.ReplaceDocument(ctx, id, dashboard)
	if err != nil {
		return nil, err
	}

	return toDashboardResponse(meta.Key, &doc), nil
}

func (db *Database) DashboardDelete(ctx context.Context, id string) error {
	_, err := db.dashboardCollection.RemoveDocument(ctx, id)
	return err
}

func (db *Database) DashboardList(ctx context.Context) ([]*model.DashboardResponse, error) {
	query := "FOR d IN @@collection RETURN d"
	cursor, _, err := db.Query(ctx, query, map[string]interface{}{"@collection": DashboardCollectionName}, busdb.ReadOperation)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	var docs []*model.DashboardResponse
	for {
		var doc model.Dashboard
		meta, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}
		docs = append(docs, toDashboardResponse(meta.Key, &doc))
	}

	return docs, err
}
