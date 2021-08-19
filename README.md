# migrationrunner ![CI Status](https://github.com/mhogar/migrationrunner/actions/workflows/CI.yml/badge.svg) [![Coverage Status](https://coveralls.io/repos/github/mhogar/migrationrunner/badge.svg)](https://coveralls.io/github/mhogar/migrationrunner) [![GoDoc](https://godoc.org/github.com/mhogar/migrationrunner?status.svg)](https://pkg.go.dev/github.com/mhogar/migrationrunner)

`migrationrunner` is a golang package for managing and running data migrations. 

A "migration" is defined as a series of one-time operations that need to be performed to set up a data solution. For example, in the case of a SQL database, a migration may create tables and insert some initial data.

A "data solution" is defined as the specific data storage being used. This may be a SQL or NoSQL database, cloud solution, some ad-hoc in-memory DB, floppy disk, or whatever floats your boat.

## Usage

The runner is designed to work on an abstract level and is independent of the specific data solution being used. As long as the interfaces are implemented as specified, the migration algorithm ~~should~~ will work correctly.

View the [GoDocs](https://pkg.go.dev/github.com/mhogar/migrationrunner) for specific interface and usage details, but in general three types will be needed to be implemented:
- `Migration`: this represents a single migration and implements the actual logic for running your migration
- `MigrationRepository`: this interface simply fetches a slice of Migrations to be run by the runner
- `MigrationCRUD`: this interface acts as a wrapper for your data solution enabling the runner to act abstractly

Once the interfaces have been implemnted to adhere to your needs, the `MigrationRunner` struct can be created. Basic usage example:
```go
func main() {
    crud := MyMigrationCRUD{} //your MigrationCRUD implementation
    repo := MyMigrationRepository{} //your MigrationRepository implementation

    runner := MigrationRunner{
        MigrationCRUD: crud,
        MigrationRepository: repo,
    }

    err := runner.MigrateUp()
    if err != nil {
        log.Fatal(err)
    }
}
```
