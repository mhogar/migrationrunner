package migrationrunner

import (
	"errors"
	"log"
)

type MigrationRunner struct {
	// Implementation of the MigrationRepository interface that will fetch the migrations to run.
	MigrationRepository MigrationRepository

	// Implementation of the MigrationCRUD interface that will perform the CRUD operations on the data solution.
	MigrationCRUD MigrationCRUD
}

// MigrateUp runs the Up method for all migrations returned by the MigrationRepository that are newer than the timestamp fetched by MigrationCRUD.GetLatestTimestamp.
// This will trigger a call to MigrationCRUD.CreateMigration for every one that's run.
//
// If there is no latest timestamp, all migrations will be run.
//
// If any errors are encountered, the whole function will be aborted and any migrations yet to run will not be run.
// Returns any such errors encountered.
func (m MigrationRunner) MigrateUp() error {
	log.Println("Migrating Up")

	log.Println("Running setup")
	err := m.MigrationCRUD.Setup()
	if err != nil {
		return chainError("error running migration setup", err)
	}

	migrations := m.MigrationRepository.GetMigrations()

	//get latest timestamp
	latestTimestamp, hasLatest, err := m.MigrationCRUD.GetLatestTimestamp()
	if err != nil {
		return chainError("error getting latest timestamp", err)
	}

	//print the timestamp if it exists
	if !hasLatest {
		log.Println("No latest timestamp found.")
	} else {
		log.Println("Latest timestamp:", latestTimestamp)
	}

	//run all migrations that are newer
	for _, migration := range migrations {
		if !hasLatest || migration.Timestamp > latestTimestamp {
			log.Println("Running", migration.Timestamp, "-", migration.Description)

			err = migration.Migrator.Up()
			if err != nil {
				return chainError("error running migration", err)
			}

			//save the migration to db to mark it as run
			err = m.MigrationCRUD.CreateMigration(migration.Timestamp)
			if err != nil {
				return chainError("error saving migration", err)
			}
		} else {
			log.Println("Skipping", migration.Timestamp)
		}
	}

	log.Println("Finished running migrations.")
	return nil
}

// MigrateDown runs the Down method for the migration whose timestamp matches the latest timestamp returned by MigrationCRUD.GetLatestTimestamp.
//
// If there is no latest timestamp, an error will be returned.
// Will also return any other errors that are encountered.
func (m MigrationRunner) MigrateDown() error {
	log.Println("Migrating Down")

	log.Println("Running setup")
	err := m.MigrationCRUD.Setup()
	if err != nil {
		return chainError("error running migration setup", err)
	}

	migrations := m.MigrationRepository.GetMigrations()

	//get latest timestamp
	latestTimestamp, hasLatest, err := m.MigrationCRUD.GetLatestTimestamp()
	if err != nil {
		return chainError("error getting latest timestamp", err)
	}

	//exit if no latest
	if !hasLatest {
		return errors.New("no migrations to migrate down")
	}

	var latestMigration Migration
	migrationFound := false

	//find migration that matches the latest timestamp
	for _, migration := range migrations {
		if migration.Timestamp == latestTimestamp {
			latestMigration = migration
			migrationFound = true
			break
		}
	}

	if !migrationFound {
		return errors.New("could not find migration with timestamp " + latestTimestamp)
	}

	log.Println("Running", latestTimestamp, "-", latestMigration.Description)

	//run the down function
	err = latestMigration.Migrator.Down()
	if err != nil {
		return chainError("error running migration", err)
	}

	//delete the migration
	err = m.MigrationCRUD.DeleteMigrationByTimestamp(latestTimestamp)
	if err != nil {
		return chainError("error deleting migration", err)
	}

	log.Println("Finished running migrations.")
	return nil
}
