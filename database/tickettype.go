package database

import (
	"context"
	"errors"

	"github.com/arangodb/go-driver"
	"github.com/iancoleman/strcase"

	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/model"
)

func toTicketType(doc *model.TicketTypeForm) *model.TicketType {
	return &model.TicketType{
		Name:             doc.Name,
		Icon:             doc.Icon,
		DefaultPlaybooks: doc.DefaultPlaybooks,
		DefaultTemplate:  doc.DefaultTemplate,
		DefaultGroups:    doc.DefaultGroups,
	}
}

func toTicketTypeResponse(key string, doc *model.TicketType) *model.TicketTypeResponse {
	return &model.TicketTypeResponse{
		ID:               key,
		Name:             doc.Name,
		Icon:             doc.Icon,
		DefaultPlaybooks: doc.DefaultPlaybooks,
		DefaultTemplate:  doc.DefaultTemplate,
		DefaultGroups:    doc.DefaultGroups,
	}
}

func (db *Database) TicketTypeCreate(ctx context.Context, tickettype *model.TicketTypeForm) (*model.TicketTypeResponse, error) {
	if tickettype == nil {
		return nil, errors.New("requires ticket type")
	}
	if tickettype.Name == "" {
		return nil, errors.New("requires ticket type name")
	}

	var doc model.TicketType
	newctx := driver.WithReturnNew(ctx, &doc)

	meta, err := db.tickettypeCollection.CreateDocument(ctx, newctx, strcase.ToKebab(tickettype.Name), toTicketType(tickettype))
	if err != nil {
		return nil, err
	}

	return toTicketTypeResponse(meta.Key, &doc), nil
}

func (db *Database) TicketTypeGet(ctx context.Context, id string) (*model.TicketTypeResponse, error) {
	var doc model.TicketType
	meta, err := db.tickettypeCollection.ReadDocument(ctx, id, &doc)
	if err != nil {
		return nil, err
	}

	return toTicketTypeResponse(meta.Key, &doc), nil
}

func (db *Database) TicketTypeUpdate(ctx context.Context, id string, tickettype *model.TicketTypeForm) (*model.TicketTypeResponse, error) {
	var doc model.TicketType
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

func (db *Database) TicketTypeList(ctx context.Context) ([]*model.TicketTypeResponse, error) {
	query := "FOR d IN @@collection RETURN d"
	cursor, _, err := db.Query(ctx, query, map[string]any{"@collection": TicketTypeCollectionName}, busdb.ReadOperation)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	var docs []*model.TicketTypeResponse
	for {
		var doc model.TicketType
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
