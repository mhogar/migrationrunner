package migrationrunner

// MigrationCRUD is an interface for preforming CRUD (Create, Read, Update, Delete) operations on a migration.
type MigrationCRUD interface {
	// Setup performs setup operations needed before other CRUD operations can be used.
	// Returns any errors.
	Setup() error

	// CreateMigration creates a new migration with the given timestamp and returns any errors.
	CreateMigration(timestamp string) error

	// GetLatestTimestamp returns the latest timestamp of all migrations.
	// If no timestamps are found then hasLatest should be false, else it should be true.
	// Also returns any errors encountered.
	GetLatestTimestamp() (timestamp string, hasLatest bool, err error)

	// DeleteMigrationByTimestamp deletes the migration with the given timestamp and returns any errors.
	DeleteMigrationByTimestamp(timestamp string) error
}

// MigrationRepository is an interface for fetching migrations to run.
type MigrationRepository interface {
	// GetMigrations returns a slice of migrations to run.
	GetMigrations() []Migration
}
