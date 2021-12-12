package database

import (
	"context"
	"errors"

	"github.com/arangodb/go-driver"
	"github.com/iancoleman/strcase"

	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/models"
)

func toTicketTemplate(doc *models.TicketTemplateForm) *models.TicketTemplate {
	return &models.TicketTemplate{Name: doc.Name, Schema: doc.Schema}
}

func toTicketTemplateResponse(key string, doc *models.TicketTemplate) *models.TicketTemplateResponse {
	return &models.TicketTemplateResponse{ID: key, Name: doc.Name, Schema: doc.Schema}
}

func (db *Database) TemplateCreate(ctx context.Context, template *models.TicketTemplateForm) (*models.TicketTemplateResponse, error) {
	if template == nil {
		return nil, errors.New("requires template")
	}
	if template.Name == "" {
		return nil, errors.New("requires template name")
	}

	var doc models.TicketTemplate
	newctx := driver.WithReturnNew(ctx, &doc)

	meta, err := db.templateCollection.CreateDocument(ctx, newctx, strcase.ToKebab(template.Name), toTicketTemplate(template))
	if err != nil {
		return nil, err
	}

	return toTicketTemplateResponse(meta.Key, &doc), nil
}

func (db *Database) TemplateGet(ctx context.Context, id string) (*models.TicketTemplateResponse, error) {
	var doc models.TicketTemplate
	meta, err := db.templateCollection.ReadDocument(ctx, id, &doc)
	if err != nil {
		return nil, err
	}

	return toTicketTemplateResponse(meta.Key, &doc), nil
}

func (db *Database) TemplateUpdate(ctx context.Context, id string, template *models.TicketTemplateForm) (*models.TicketTemplateResponse, error) {
	var doc models.TicketTemplate
	ctx = driver.WithReturnNew(ctx, &doc)

	meta, err := db.templateCollection.ReplaceDocument(ctx, id, toTicketTemplate(template))
	if err != nil {
		return nil, err
	}

	return toTicketTemplateResponse(meta.Key, &doc), nil
}

func (db *Database) TemplateDelete(ctx context.Context, id string) error {
	_, err := db.templateCollection.RemoveDocument(ctx, id)
	return err
}

func (db *Database) TemplateList(ctx context.Context) ([]*models.TicketTemplateResponse, error) {
	query := "FOR d IN @@collection RETURN d"
	cursor, _, err := db.Query(ctx, query, map[string]interface{}{"@collection": TemplateCollectionName}, busdb.ReadOperation)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	var docs []*models.TicketTemplateResponse
	for {
		var doc models.TicketTemplate
		meta, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}
		docs = append(docs, toTicketTemplateResponse(meta.Key, &doc))
	}

	return docs, err
}
