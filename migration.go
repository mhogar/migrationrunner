package migrationrunner

// Migration is a struct that represents a single migration.
type Migration struct {
	// The required timestamp string for the migration.
	// The format is arbitrary but must be orderable using standard string comparison.
	Timestamp string

	// An optional string that will be displayed alongside the timestamp in logs.
	Description string

	// Implementation of the Migrator interface.
	Migrator Migrator
}

// Migrator is an interface for running a migration.
//
// A new struct that implements this interface should be created for each migration.
// Any additional dependencies required to run your migrations can be added as fields to this struct.
type Migrator interface {
	// Up runs the migration and returns any errors.
	Up() error

	// Down runs the inverse migration and returns any errors.
	Down() error
}
