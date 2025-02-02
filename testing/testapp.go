package testing

import (
	"fmt"
	"os"
	"testing"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"

	"github.com/SecurityBrewery/catalyst/app"
	"github.com/SecurityBrewery/catalyst/migrations"
)

func App(t *testing.T) (*pocketbase.PocketBase, *Counter, func()) {
	t.Helper()

	temp, err := os.MkdirTemp("", "catalyst_test_data")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	baseApp, err := app.App(temp, true)
	if err != nil {
		t.Fatalf("failed to create app: %v", err)
	}

	if err := baseApp.Bootstrap(); err != nil {
		t.Fatalf("failed to bootstrap: %v", err)
	}

	baseApp.Settings().Logs.MaxDays = 0

	defaultTestData(t, baseApp)

	counter := countEvents(baseApp)

	return baseApp, counter, func() { _ = os.RemoveAll(temp) }
}

func generateAdminToken(t *testing.T, baseApp core.App, email string) (string, error) {
	t.Helper()

	record, err := baseApp.FindAuthRecordByEmail(core.CollectionNameSuperusers, email)
	if err != nil {
		return "", fmt.Errorf("failed to find admin: %w", err)
	}

	return record.NewAuthToken()
}

func generateRecordToken(t *testing.T, baseApp core.App, email string) (string, error) {
	t.Helper()

	record, err := baseApp.FindAuthRecordByEmail(migrations.UserCollectionID, email)
	if err != nil {
		return "", fmt.Errorf("failed to find record: %w", err)
	}

	return record.NewAuthToken()
}

func countEvents(t *pocketbase.PocketBase) *Counter {
	c := NewCounter()

	t.OnModelAfterCreateSuccess().BindFunc(count[*core.ModelEvent](c, "OnModelAfterCreateSuccess"))
	t.OnModelAfterCreateError().BindFunc(count[*core.ModelErrorEvent](c, "OnModelAfterCreateError"))
	t.OnModelAfterUpdateSuccess().BindFunc(count[*core.ModelEvent](c, "OnModelAfterUpdateSuccess"))
	t.OnModelAfterUpdateError().BindFunc(count[*core.ModelErrorEvent](c, "OnModelAfterUpdateError"))
	t.OnModelAfterDeleteSuccess().BindFunc(count[*core.ModelEvent](c, "OnModelAfterDeleteSuccess"))
	t.OnModelAfterDeleteError().BindFunc(count[*core.ModelErrorEvent](c, "OnModelAfterDeleteError"))

	t.OnRecordsListRequest().BindFunc(count[*core.RecordsListRequestEvent](c, "OnRecordsListRequest"))
	t.OnRecordViewRequest().BindFunc(count[*core.RecordRequestEvent](c, "OnRecordViewRequest"))
	t.OnRecordCreateRequest().BindFunc(count[*core.RecordRequestEvent](c, "OnRecordCreateRequest"))
	t.OnRecordUpdateRequest().BindFunc(count[*core.RecordRequestEvent](c, "OnRecordUpdateRequest"))
	t.OnRecordDeleteRequest().BindFunc(count[*core.RecordRequestEvent](c, "OnRecordDeleteRequest"))
	t.OnBatchRequest().BindFunc(count[*core.BatchRequestEvent](c, "OnBatchRequest"))

	return c
}

type Nexter interface {
	Next() error
}

func count[T Nexter](c *Counter, name string) func(_ T) error {
	return func(t T) error {
		defer c.Increment(name)

		return t.Next()
	}
}
