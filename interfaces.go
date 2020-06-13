package migrationrunner

// Migration is an interface for running a migration.
type Migration interface {
	// GetTimestamp should return the timestamp for this migration.
	GetTimestamp() string

	// Up runs the migration and returns any errors.
	Up() error

	// Down runs the inverse migration and returns any errors.
	Down() error
}

// MigrationCRUD is an interface for preforming CRUD operations on a migration.
type MigrationCRUD interface {
	// Setup performs setup operations needed before other CRUD operations can be used.
	// Returns any errors.
	Setup() error

	// CreateMigration creates a new migration with the given timestamp and returns any errors.
	CreateMigration(timestamp string) error

	// GetLatestTimestamp returns the latest timestamp of all migrations. If no timestamps are found then
	// hasLatest should be false, else it should be true. Also returns any errors encountered.
	GetLatestTimestamp() (timestamp string, hasLatest bool, err error)

	// DeleteMigrationByTimestamp deletes the migration with the given timestamp and returns any errors.
	DeleteMigrationByTimestamp(timestamp string) error
}

// MigrationRepository is an interface for fetching migrations to run.
type MigrationRepository interface {
	// GetMigrations returns the migrations to run
	GetMigrations() []Migration
}
