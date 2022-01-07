package database

import (
	"context"
	"errors"

	"github.com/arangodb/go-driver"
	"github.com/iancoleman/strcase"

	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/model"
)

func toTicketTemplate(doc *model.TicketTemplateForm) *model.TicketTemplate {
	return &model.TicketTemplate{Name: doc.Name, Schema: doc.Schema}
}

func toTicketTemplateResponse(key string, doc *model.TicketTemplate) *model.TicketTemplateResponse {
	return &model.TicketTemplateResponse{ID: key, Name: doc.Name, Schema: doc.Schema}
}

func (db *Database) TemplateCreate(ctx context.Context, template *model.TicketTemplateForm) (*model.TicketTemplateResponse, error) {
	if template == nil {
		return nil, errors.New("requires template")
	}
	if template.Name == "" {
		return nil, errors.New("requires template name")
	}

	var doc model.TicketTemplate
	newctx := driver.WithReturnNew(ctx, &doc)

	meta, err := db.templateCollection.CreateDocument(ctx, newctx, strcase.ToKebab(template.Name), toTicketTemplate(template))
	if err != nil {
		return nil, err
	}

	return toTicketTemplateResponse(meta.Key, &doc), nil
}

func (db *Database) TemplateGet(ctx context.Context, id string) (*model.TicketTemplateResponse, error) {
	var doc model.TicketTemplate
	meta, err := db.templateCollection.ReadDocument(ctx, id, &doc)
	if err != nil {
		return nil, err
	}

	return toTicketTemplateResponse(meta.Key, &doc), nil
}

func (db *Database) TemplateUpdate(ctx context.Context, id string, template *model.TicketTemplateForm) (*model.TicketTemplateResponse, error) {
	var doc model.TicketTemplate
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

func (db *Database) TemplateList(ctx context.Context) ([]*model.TicketTemplateResponse, error) {
	query := "FOR d IN @@collection RETURN d"
	cursor, _, err := db.Query(ctx, query, map[string]interface{}{"@collection": TemplateCollectionName}, busdb.ReadOperation)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	var docs []*model.TicketTemplateResponse
	for {
		var doc model.TicketTemplate
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
