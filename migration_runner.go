package migrationrunner

import (
	"errors"
	"log"
)

// MigrateUp runs the Up method for all migrations returned by migrationRepo that are newer than the latest timestamp fetched by db;
// and will create a new migration for every one that's run. If there is no latest timestamp, all migrations will be run.
// If any errors are encountered, the whole function will be aborted and any migrations yet to run will not be run.
// Returns the errors encountered.
func MigrateUp(migrationRepo MigrationRepository, db MigrationCRUD) error {
	log.Println("Migrating Up")

	migrations := migrationRepo.GetMigrations()

	//get latest timestamp
	latestTimestamp, hasLatest, err := db.GetLatestTimestamp()
	if err != nil {
		return ChainError("error getting latest timestamp", err)
	}

	//print the timestamp if it exists
	if !hasLatest {
		log.Println("No timestamps found.")
	} else {
		log.Println("Latest timestamp:", latestTimestamp)
	}

	//run all migrations that are newer
	for _, migration := range migrations {
		timestamp := migration.GetTimestamp()

		if !hasLatest || timestamp > latestTimestamp {
			log.Println("Running", timestamp)

			err = migration.Up()
			if err != nil {
				return ChainError("error running migration", err)
			}

			//save the migration to db to mark it as run
			err = db.CreateMigration(timestamp)
			if err != nil {
				return ChainError("error saving migration", err)
			}
		} else {
			log.Println("Skipping", timestamp)
		}
	}

	log.Println("Finished running migrations.")
	return nil
}

// MigrateDown runs the Down method for the migration whose timestamp matches the latest timestamp returned by db.
// If there is no latest timestamp, an error will be returned. Will return any other errors that are encountered.
func MigrateDown(migrationRepo MigrationRepository, db MigrationCRUD) error {
	log.Println("Migrating Down")

	migrations := migrationRepo.GetMigrations()

	//get latest timestamp
	latestTimestamp, hasLatest, err := db.GetLatestTimestamp()
	if err != nil {
		return ChainError("error getting latest timestamp", err)
	}

	//exit if no latest
	if !hasLatest {
		return errors.New("no migrations to migrate down")
	}

	var latestMigration Migration = nil

	//find migration that matches the latest timestamp
	for _, migration := range migrations {
		if migration.GetTimestamp() == latestTimestamp {
			latestMigration = migration
			break
		}
	}

	if latestMigration == nil {
		return errors.New("could not find migration with timestamp " + latestTimestamp)
	}

	log.Println("Running " + latestTimestamp)

	//run the down function
	err = latestMigration.Down()
	if err != nil {
		return ChainError("error running migration", err)
	}

	//remove migration from database
	err = db.DeleteMigrationByTimestamp(latestTimestamp)
	if err != nil {
		return ChainError("error deleting migration", err)
	}

	log.Println("Finished running migrations.")
	return nil
}
