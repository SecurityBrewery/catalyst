package busdb

import (
	"context"
	"errors"

	"github.com/arangodb/go-driver"

	"github.com/SecurityBrewery/catalyst/generated/models"
	"github.com/SecurityBrewery/catalyst/time"
)

const LogCollectionName = "logs"

func (db *BusDatabase) LogCreate(ctx context.Context, reference, message string) (*models.LogEntry, error) {
	user, ok := UserFromContext(ctx)
	if !ok {
		return nil, errors.New("no user in context")
	}

	logentry := &models.LogEntry{
		Reference: reference,
		Created:   time.Now(),
		Creator:   user.ID,
		Message:   message,
	}

	doc := models.LogEntry{}
	_, err := db.logCollection.CreateDocument(driver.WithReturnNew(ctx, &doc), logentry)
	if err != nil {
		return nil, err
	}

	return &doc, db.bus.PublishUpdate([]driver.DocumentID{driver.DocumentID(logentry.Reference)})
}

func (db *BusDatabase) LogBatchCreate(ctx context.Context, logEntryForms []*models.LogEntry) error {
	user, ok := UserFromContext(ctx)
	if !ok {
		return errors.New("no user in context")
	}

	var ids []driver.DocumentID
	var logentries []*models.LogEntry
	for _, logEntryForm := range logEntryForms {
		logentry := &models.LogEntry{
			Reference: logEntryForm.Reference,
			Created:   time.Now(),
			Creator:   user.ID,
			Message:   logEntryForm.Message,
		}

		logentries = append(logentries, logentry)
		ids = append(ids, driver.DocumentID(logentry.Reference))
	}

	_, errs, err := db.logCollection.CreateDocuments(ctx, logentries)
	if err != nil {
		return err
	}
	err = errs.FirstNonNil()
	if err != nil {
		return err
	}

	return db.bus.PublishUpdate(ids)
}

func (db *BusDatabase) LogList(ctx context.Context, reference string) ([]*models.LogEntry, error) {
	query := "FOR d IN @@collection FILTER d.reference == @reference SORT d.created DESC RETURN d"
	cursor, err := db.internal.Query(ctx, query, map[string]interface{}{
		"@collection": LogCollectionName,
		"reference":   reference,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	var docs []*models.LogEntry
	for {
		var doc models.LogEntry
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}
		docs = append(docs, &doc)
	}

	return docs, err
}
