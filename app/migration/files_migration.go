package migration

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/upload/uploader"
)

type filesMigration struct{}

func newFilesMigration() func() (migration, error) {
	return func() (migration, error) {
		return filesMigration{}, nil
	}
}

func (filesMigration) name() string { return "005_pocketbase_files_to_tusd" }

func (filesMigration) up(ctx context.Context, queries *sqlc.Queries, dir string, uploader *uploader.Uploader) error {
	oldUploadDir := filepath.Join(dir, "storage")
	if _, err := os.Stat(oldUploadDir); os.IsNotExist(err) {
		// If the old upload directory does not exist, we assume no migration is needed.
		return nil
	}

	oldUploadRoot, err := os.OpenRoot(oldUploadDir)
	if err != nil {
		return fmt.Errorf("open old uploads root: %w", err)
	}

	files, err := queries.ListFiles(ctx, sqlc.ListFilesParams{
		Limit:  1000,
		Offset: 0,
	})
	if err != nil {
		return fmt.Errorf("list files: %w", err)
	}

	for _, file := range files {
		data, err := fs.ReadFile(oldUploadRoot.FS(), filepath.Join(file.ID, file.Blob))
		if err != nil {
			return fmt.Errorf("read file %s: %w", file.Blob, err)
		}

		if _, err := uploader.CreateFile(file.ID, file.Name, data); err != nil {
			return err
		}
	}

	return nil
}
