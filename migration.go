package migrationrunner

// Migration is a struct that represents a single migration.
type Migration struct {
	// The required timestamp string for the migration.
	// The format is flexible, but must be orderable using standard string comparison.
	//
	// eg. "20200628151601".
	// This reads as 2020, June 28, at 15:16, migration #1.
	Timestamp string

	// An optional string that will be displayed alongside the timestamp in logs.
	Description string

	// Implementation of the Migrator interface.
	Migrator Migrator
}

// Migrator is an interface for running a migration.
type Migrator interface {
	// Up runs the migration and returns any errors.
	Up() error

	// Down runs the inverse migration and returns any errors.
	Down() error
}
