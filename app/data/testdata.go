package data

import (
	_ "embed"
	"log/slog"
	"os"
	"path"
	"testing"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

const (
	AdminEmail   = "admin@catalyst-soar.com"
	AnalystEmail = "analyst@catalyst-soar.com"
)

//go:embed testdata.sql
var testdata string

func DefaultTestData(t *testing.T, dir string, queries *sqlc.Queries) {
	t.Helper()

	if _, err := queries.WriteDB.ExecContext(t.Context(), testdata); err != nil {
		t.Fatalf("failed to execute test data: %v", err)
	}

	files, err := queries.ListFiles(t.Context(), sqlc.ListFilesParams{
		Limit: 1000, // TODO
	})
	if err != nil {
		t.Fatalf("failed to list files: %v", err)
	}

	for _, file := range files {
		_ = os.MkdirAll(path.Join(dir, "uploads", file.ID), 0o755)

		infoFilePath := path.Join(dir, "uploads", file.ID+".info")
		slog.InfoContext(t.Context(), "Creating file info", "path", infoFilePath)

		if err := os.WriteFile(infoFilePath, []byte(`{"MetaData":{"filetype":"text/plain"}}`), 0o644); err != nil {
			t.Fatalf("failed to write file %s: %v", file.Name, err)
		}

		if err := os.WriteFile(path.Join(dir, "uploads", file.ID, file.Blob), []byte("hello"), 0o644); err != nil {
			t.Fatalf("failed to write file %s: %v", file.Name, err)
		}
	}
}
