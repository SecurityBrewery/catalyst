package migration

import "fmt"

var migrationGenerators = []func() (migration, error){
	newSQLMigration("000_create_tables"),
	newSQLMigration("002_create_defaultdata"),
	newSQLMigration("003_create_groups"),
	newSQLMigration("004_cleanup_pocketbase"),
}

func migrations(version int) ([]migration, error) {
	var migrations []migration

	if version < 0 || version > len(migrationGenerators) {
		return nil, fmt.Errorf("invalid migration version: %d", version)
	}

	if version == len(migrationGenerators) {
		return migrations, nil // No migrations to apply
	}

	for _, migrationFunc := range migrationGenerators[version:] {
		migration, err := migrationFunc()
		if err != nil {
			return nil, fmt.Errorf("failed to create migration: %w", err)
		}

		migrations = append(migrations, migration)
	}

	return migrations, nil
}
