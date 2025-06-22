package data

import (
	_ "embed"
	"testing"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

const (
	AdminEmail   = "admin@catalyst-soar.com"
	AnalystEmail = "analyst@catalyst-soar.com"
)

//go:embed testdata.sql
var testdata string

func DefaultTestData(t *testing.T, queries *sqlc.Queries) {
	t.Helper()

	if _, err := queries.WriteDB.ExecContext(t.Context(), testdata); err != nil {
		t.Fatalf("failed to execute test data: %v", err)
	}
}
