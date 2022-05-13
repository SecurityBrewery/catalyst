package busdb

import (
	"context"
	"errors"

	"github.com/arangodb/go-driver"

	"github.com/SecurityBrewery/catalyst/bus"
	"github.com/SecurityBrewery/catalyst/generated/api"
	"github.com/SecurityBrewery/catalyst/generated/model"
)

// BusDatabase
//  1. Save entry to log
//  2. Send update ticket to bus
//  3. Add document to index
type BusDatabase struct {
	internal      driver.Database
	logCollection driver.Collection
	bus           *bus.Bus
	// index         *index.Index
}

func NewDatabase(ctx context.Context, internal driver.Database, b *bus.Bus) (*BusDatabase, error) {
	logCollection, err := internal.Collection(ctx, LogCollectionName)
	if err != nil {
		return nil, err
	}

	return &BusDatabase{
		internal:      internal,
		logCollection: logCollection,
		bus:           b,
	}, nil
}

type Operation struct {
	Type bus.DatabaseUpdateType
	Ids  []driver.DocumentID
}

var (
	CreateOperation = &Operation{Type: bus.DatabaseEntryCreated}
	ReadOperation   = &Operation{Type: bus.DatabaseEntryRead}
)

func (db *BusDatabase) Query(ctx context.Context, query string, vars map[string]any, operation *Operation) (cur driver.Cursor, logs *model.LogEntry, err error) {
	defer func() { err = toHTTPErr(err) }()

	cur, err = db.internal.Query(ctx, query, vars)
	if err != nil {
		return nil, nil, err
	}

	switch {
	case operation.Type == bus.DatabaseEntryCreated, operation.Type == bus.DatabaseEntryUpdated:
		db.bus.DatabaseChannel.Publish(&bus.DatabaseUpdateMsg{IDs: operation.Ids, Type: operation.Type})
	}

	return cur, logs, err
}

func (db *BusDatabase) Remove(ctx context.Context) (err error) {
	defer func() { err = toHTTPErr(err) }()

	return db.internal.Remove(ctx)
}

func (db *BusDatabase) Collection(ctx context.Context, name string) (col driver.Collection, err error) {
	defer func() { err = toHTTPErr(err) }()

	return db.internal.Collection(ctx, name)
}

type Collection[T any] struct {
	internal driver.Collection
	db       *BusDatabase
}

func NewCollection[T any](internal driver.Collection, db *BusDatabase) *Collection[T] {
	return &Collection[T]{internal: internal, db: db}
}

func (c *Collection[T]) CreateDocument(ctx, newctx context.Context, key string, document *T) (meta driver.DocumentMeta, err error) {
	defer func() { err = toHTTPErr(err) }()

	meta, err = c.internal.CreateDocument(newctx, &Keyed[T]{Key: key, Doc: document})
	if err != nil {
		return meta, err
	}

	c.db.bus.DatabaseChannel.Publish(&bus.DatabaseUpdateMsg{IDs: []driver.DocumentID{meta.ID}, Type: bus.DatabaseEntryCreated})

	return meta, nil
}

func (c *Collection[T]) CreateEdge(ctx, newctx context.Context, edge *driver.EdgeDocument) (meta driver.DocumentMeta, err error) {
	defer func() { err = toHTTPErr(err) }()

	meta, err = c.internal.CreateDocument(newctx, edge)
	if err != nil {
		return meta, err
	}

	c.db.bus.DatabaseChannel.Publish(&bus.DatabaseUpdateMsg{IDs: []driver.DocumentID{meta.ID}, Type: bus.DatabaseEntryCreated})

	return meta, nil
}

func (c *Collection[T]) CreateEdges(ctx context.Context, edges []*driver.EdgeDocument) (meta driver.DocumentMetaSlice, err error) {
	defer func() { err = toHTTPErr(err) }()

	metas, errs, err := c.internal.CreateDocuments(ctx, edges)
	if err != nil {
		return nil, err
	}
	if errs.FirstNonNil() != nil {
		return nil, errs.FirstNonNil()
	}

	var ids []driver.DocumentID
	for _, meta := range metas {
		ids = append(ids, meta.ID)
	}

	c.db.bus.DatabaseChannel.Publish(&bus.DatabaseUpdateMsg{IDs: ids, Type: bus.DatabaseEntryCreated})

	return metas, nil
}

func (c *Collection[T]) DocumentExists(ctx context.Context, id string) (exists bool, err error) {
	defer func() { err = toHTTPErr(err) }()

	return c.internal.DocumentExists(ctx, id)
}

func (c *Collection[T]) ReadDocument(ctx context.Context, key string, result *T) (meta driver.DocumentMeta, err error) {
	defer func() { err = toHTTPErr(err) }()

	meta, err = c.internal.ReadDocument(ctx, key, result)

	return
}

func (c *Collection[T]) UpdateDocument(ctx context.Context, key string, update any) (meta driver.DocumentMeta, err error) {
	defer func() { err = toHTTPErr(err) }()

	meta, err = c.internal.UpdateDocument(ctx, key, update)
	if err != nil {
		return meta, err
	}

	c.db.bus.DatabaseChannel.Publish(&bus.DatabaseUpdateMsg{IDs: []driver.DocumentID{meta.ID}, Type: bus.DatabaseEntryUpdated})

	return meta, nil
}

func (c *Collection[T]) ReplaceDocument(ctx context.Context, key string, document *T) (meta driver.DocumentMeta, err error) {
	defer func() { err = toHTTPErr(err) }()

	meta, err = c.internal.ReplaceDocument(ctx, key, document)
	if err != nil {
		return meta, err
	}

	c.db.bus.DatabaseChannel.Publish(&bus.DatabaseUpdateMsg{IDs: []driver.DocumentID{meta.ID}, Type: bus.DatabaseEntryUpdated})

	return meta, nil
}

func (c *Collection[T]) RemoveDocument(ctx context.Context, formatInt string) (meta driver.DocumentMeta, err error) {
	defer func() { err = toHTTPErr(err) }()

	return c.internal.RemoveDocument(ctx, formatInt)
}

func (c *Collection[T]) Truncate(ctx context.Context) (err error) {
	defer func() { err = toHTTPErr(err) }()

	return c.internal.Truncate(ctx)
}

func toHTTPErr(err error) error {
	if err != nil {
		ae := driver.ArangoError{}
		if errors.As(err, &ae) {
			return &api.HTTPError{Status: ae.Code, Internal: err}
		}

		return err
	}

	return nil
}
