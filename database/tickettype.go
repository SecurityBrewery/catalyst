package database

import (
	"context"
	"errors"

	"github.com/arangodb/go-driver"
	"github.com/iancoleman/strcase"

	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/models"
)

func toTicketType(doc *models.TicketTypeForm) *models.TicketType {
	return &models.TicketType{
		Name:             doc.Name,
		Icon:             doc.Icon,
		DefaultPlaybooks: doc.DefaultPlaybooks,
		DefaultTemplate:  doc.DefaultTemplate,
		DefaultGroups:    doc.DefaultGroups,
	}
}

func toTicketTypeResponse(key string, doc *models.TicketType) *models.TicketTypeResponse {
	return &models.TicketTypeResponse{
		ID:               key,
		Name:             doc.Name,
		Icon:             doc.Icon,
		DefaultPlaybooks: doc.DefaultPlaybooks,
		DefaultTemplate:  doc.DefaultTemplate,
		DefaultGroups:    doc.DefaultGroups,
	}
}

func (db *Database) TicketTypeCreate(ctx context.Context, tickettype *models.TicketTypeForm) (*models.TicketTypeResponse, error) {
	if tickettype == nil {
		return nil, errors.New("requires ticket type")
	}
	if tickettype.Name == "" {
		return nil, errors.New("requires ticket type name")
	}

	var doc models.TicketType
	newctx := driver.WithReturnNew(ctx, &doc)

	meta, err := db.tickettypeCollection.CreateDocument(ctx, newctx, strcase.ToKebab(tickettype.Name), toTicketType(tickettype))
	if err != nil {
		return nil, err
	}

	return toTicketTypeResponse(meta.Key, &doc), nil
}

func (db *Database) TicketTypeGet(ctx context.Context, id string) (*models.TicketTypeResponse, error) {
	var doc models.TicketType
	meta, err := db.tickettypeCollection.ReadDocument(ctx, id, &doc)
	if err != nil {
		return nil, err
	}

	return toTicketTypeResponse(meta.Key, &doc), nil
}

func (db *Database) TicketTypeUpdate(ctx context.Context, id string, tickettype *models.TicketTypeForm) (*models.TicketTypeResponse, error) {
	var doc models.TicketType
	ctx = driver.WithReturnNew(ctx, &doc)

	meta, err := db.tickettypeCollection.ReplaceDocument(ctx, id, toTicketType(tickettype))
	if err != nil {
		return nil, err
	}

	return toTicketTypeResponse(meta.Key, &doc), nil
}

func (db *Database) TicketTypeDelete(ctx context.Context, id string) error {
	_, err := db.tickettypeCollection.RemoveDocument(ctx, id)
	return err
}

func (db *Database) TicketTypeList(ctx context.Context) ([]*models.TicketTypeResponse, error) {
	query := "FOR d IN @@collection RETURN d"
	cursor, _, err := db.Query(ctx, query, map[string]interface{}{"@collection": TicketTypeCollectionName}, busdb.ReadOperation)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	var docs []*models.TicketTypeResponse
	for {
		var doc models.TicketType
		meta, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}
		docs = append(docs, toTicketTypeResponse(meta.Key, &doc))
	}

	return docs, err
}
