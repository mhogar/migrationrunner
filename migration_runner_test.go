package migrationrunner_test

import (
	"errors"
	"testing"

	"github.com/mhogar/migrationrunner"
	"github.com/mhogar/migrationrunner/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MigrationRunnerTestSuite struct {
	suite.Suite
	MigrationRepositoryMock mocks.MigrationRepository
	MigrationCRUDMock       mocks.MigrationCRUD
	MigrationRunner         migrationrunner.MigrationRunner
}

func (suite *MigrationRunnerTestSuite) SetupTest() {
	suite.MigrationRepositoryMock = mocks.MigrationRepository{}
	suite.MigrationCRUDMock = mocks.MigrationCRUD{}
	suite.MigrationRunner = migrationrunner.MigrationRunner{
		MigrationRepository: &suite.MigrationRepositoryMock,
		MigrationCRUD:       &suite.MigrationCRUDMock,
	}
}

func (suite *MigrationRunnerTestSuite) TestMigrateUp_WithErrorRunningSetup_ReturnsError() {
	//arrange
	errMessage := "Setup mock error"
	suite.MigrationCRUDMock.On("Setup").Return(errors.New(errMessage))

	//act
	err := suite.MigrationRunner.MigrateUp()

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), errMessage)
}

func (suite *MigrationRunnerTestSuite) TestMigrateUp_WithNoLatestTimestamp_RunsAllMigrations() {
	//arrange
	migrations := createMigrations("01", "04", "08", "10")

	suite.MigrationCRUDMock.On("Setup").Return(nil)
	suite.MigrationRepositoryMock.On("GetMigrations").Return(migrations)
	suite.MigrationCRUDMock.On("GetLatestTimestamp").Return("100: this can be anything", false, nil)
	suite.MigrationCRUDMock.On("CreateMigration", mock.Anything).Return(nil)

	//act
	err := suite.MigrationRunner.MigrateUp()

	//assert
	suite.NoError(err)

	getMigratorMock(migrations[0]).AssertCalled(suite.T(), "Up")
	getMigratorMock(migrations[1]).AssertCalled(suite.T(), "Up")
	getMigratorMock(migrations[2]).AssertCalled(suite.T(), "Up")
	getMigratorMock(migrations[3]).AssertCalled(suite.T(), "Up")

	suite.MigrationCRUDMock.AssertCalled(suite.T(), "CreateMigration", mock.MatchedBy(func(timestamp string) bool {
		return timestamp == "01"
	}))
	suite.MigrationCRUDMock.AssertCalled(suite.T(), "CreateMigration", mock.MatchedBy(func(timestamp string) bool {
		return timestamp == "04"
	}))
	suite.MigrationCRUDMock.AssertCalled(suite.T(), "CreateMigration", mock.MatchedBy(func(timestamp string) bool {
		return timestamp == "08"
	}))
	suite.MigrationCRUDMock.AssertCalled(suite.T(), "CreateMigration", mock.MatchedBy(func(timestamp string) bool {
		return timestamp == "10"
	}))
}

func (suite *MigrationRunnerTestSuite) TestMigrateUp_RunsAllMigrationsWithTimestampGreaterThanLatest() {
	//arrange
	migrations := createMigrations("01", "04", "08", "10")

	suite.MigrationCRUDMock.On("Setup").Return(nil)
	suite.MigrationRepositoryMock.On("GetMigrations").Return(migrations)
	suite.MigrationCRUDMock.On("GetLatestTimestamp").Return("05", true, nil)
	suite.MigrationCRUDMock.On("CreateMigration", mock.Anything).Return(nil)

	//act
	err := suite.MigrationRunner.MigrateUp()

	//assert
	suite.NoError(err)

	getMigratorMock(migrations[0]).AssertNotCalled(suite.T(), "Up")
	getMigratorMock(migrations[1]).AssertNotCalled(suite.T(), "Up")
	getMigratorMock(migrations[2]).AssertCalled(suite.T(), "Up")
	getMigratorMock(migrations[3]).AssertCalled(suite.T(), "Up")

	suite.MigrationCRUDMock.AssertNotCalled(suite.T(), "CreateMigration", mock.MatchedBy(func(timestamp string) bool {
		return timestamp == "01"
	}))
	suite.MigrationCRUDMock.AssertNotCalled(suite.T(), "CreateMigration", mock.MatchedBy(func(timestamp string) bool {
		return timestamp == "04"
	}))
	suite.MigrationCRUDMock.AssertCalled(suite.T(), "CreateMigration", mock.MatchedBy(func(timestamp string) bool {
		return timestamp == "08"
	}))
	suite.MigrationCRUDMock.AssertCalled(suite.T(), "CreateMigration", mock.MatchedBy(func(timestamp string) bool {
		return timestamp == "10"
	}))
}

func (suite *MigrationRunnerTestSuite) TestMigrateUp_WhereGetLatestTimestampReturnsError_ReturnsError() {
	//arrange
	errMessage := "GetLatestTimestamp mock error"

	suite.MigrationCRUDMock.On("Setup").Return(nil)
	suite.MigrationRepositoryMock.On("GetMigrations").Return(nil)
	suite.MigrationCRUDMock.On("GetLatestTimestamp").Return("", false, errors.New(errMessage))

	//act
	err := suite.MigrationRunner.MigrateUp()

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), errMessage)
}

func (suite *MigrationRunnerTestSuite) TestMigrateUp_WhereMigrationUpReturnsError_ReturnsError() {
	//arrange
	errMessage := "Up mock error"
	migration := createMigration("01", errors.New(errMessage), nil)

	migrations := []migrationrunner.Migration{
		migration,
	}

	suite.MigrationCRUDMock.On("Setup").Return(nil)
	suite.MigrationRepositoryMock.On("GetMigrations").Return(migrations)
	suite.MigrationCRUDMock.On("GetLatestTimestamp").Return("0", true, nil)

	//act
	err := suite.MigrationRunner.MigrateUp()

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), errMessage)
}

func (suite *MigrationRunnerTestSuite) TestMigrateUp_WhereCreateMigrationReturnsError_ReturnsError() {
	//arrange
	errMessage := "CreateMigration mock error"

	migration := createMigration("01", nil, nil)
	migrations := []migrationrunner.Migration{
		migration,
	}

	suite.MigrationCRUDMock.On("Setup").Return(nil)
	suite.MigrationRepositoryMock.On("GetMigrations").Return(migrations)
	suite.MigrationCRUDMock.On("GetLatestTimestamp").Return("0", true, nil)
	suite.MigrationCRUDMock.On("CreateMigration", mock.Anything).Return(errors.New(errMessage))

	//act
	err := suite.MigrationRunner.MigrateUp()

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), errMessage)
}

func (suite *MigrationRunnerTestSuite) TestMigrateDown_WithErrorRunningSetup_ReturnsError() {
	//arrange
	errMessage := "Setup mock error"
	suite.MigrationCRUDMock.On("Setup").Return(errors.New(errMessage))

	//act
	err := suite.MigrationRunner.MigrateDown()

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), errMessage)
}

func (suite *MigrationRunnerTestSuite) TestMigrateDown_WhereGetLatestTimestampReturnsError_ReturnsError() {
	//arrange
	errMessage := "GetLatestTimestamp mock error"

	suite.MigrationCRUDMock.On("Setup").Return(nil)
	suite.MigrationRepositoryMock.On("GetMigrations").Return(nil)
	suite.MigrationCRUDMock.On("GetLatestTimestamp").Return("", false, errors.New(errMessage))

	//act
	err := suite.MigrationRunner.MigrateDown()

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), errMessage)
}

func (suite *MigrationRunnerTestSuite) TestMigrateDown_WithNoLatestTimestamp_ReturnsError() {
	//arrange
	suite.MigrationCRUDMock.On("Setup").Return(nil)
	suite.MigrationRepositoryMock.On("GetMigrations").Return(nil)
	suite.MigrationCRUDMock.On("GetLatestTimestamp").Return("", false, nil)

	//act
	err := suite.MigrationRunner.MigrateDown()

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), "no migrations")
}

func (suite *MigrationRunnerTestSuite) TestMigrateDown_WhereLatestTimestampNotFoundInMigrations_ReturnsError() {
	//arrange
	migrations := createMigrations("01", "04", "08", "10")

	suite.MigrationCRUDMock.On("Setup").Return(nil)
	suite.MigrationRepositoryMock.On("GetMigrations").Return(migrations)
	suite.MigrationCRUDMock.On("GetLatestTimestamp").Return("05", true, nil)

	//act
	err := suite.MigrationRunner.MigrateDown()

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), "could not find migration")
}

func (suite *MigrationRunnerTestSuite) TestMigrateDown_WhereMigrationDownReturnsError_ReturnError() {
	//arrange
	errMessage := "Down mock error"
	migration := createMigration("01", nil, errors.New(errMessage))

	migrations := []migrationrunner.Migration{
		migration,
	}

	suite.MigrationCRUDMock.On("Setup").Return(nil)
	suite.MigrationRepositoryMock.On("GetMigrations").Return(migrations)
	suite.MigrationCRUDMock.On("GetLatestTimestamp").Return("01", true, nil)
	suite.MigrationCRUDMock.On("DeleteMigrationByTimestamp", mock.Anything).Return(nil)

	//act
	err := suite.MigrationRunner.MigrateDown()

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), errMessage)
}

func (suite *MigrationRunnerTestSuite) TestMigrateDown_WhereDeleteMigrationReturnsError_ReturnError() {
	//arrange
	errMessage := "DeleteMigrationByTimestamp mock error"

	migrations := createMigrations("01", "04", "08", "10")

	suite.MigrationCRUDMock.On("Setup").Return(nil)
	suite.MigrationRepositoryMock.On("GetMigrations").Return(migrations)
	suite.MigrationCRUDMock.On("GetLatestTimestamp").Return("08", true, nil)
	suite.MigrationCRUDMock.On("DeleteMigrationByTimestamp", mock.Anything).Return(errors.New(errMessage))

	//act
	err := suite.MigrationRunner.MigrateDown()

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), errMessage)
}

func (suite *MigrationRunnerTestSuite) TestMigrateDown_RunsDownFunctionForMigrationMatchingLatestTimestamp() {
	//arrange
	latestTimestamp := "08"
	migrations := createMigrations("01", "04", "08", "10")

	suite.MigrationCRUDMock.On("Setup").Return(nil)
	suite.MigrationRepositoryMock.On("GetMigrations").Return(migrations)
	suite.MigrationCRUDMock.On("GetLatestTimestamp").Return(latestTimestamp, true, nil)
	suite.MigrationCRUDMock.On("DeleteMigrationByTimestamp", mock.Anything).Return(nil)

	//act
	err := suite.MigrationRunner.MigrateDown()

	//assert
	suite.NoError(err)

	getMigratorMock(migrations[0]).AssertNotCalled(suite.T(), "Down")
	getMigratorMock(migrations[1]).AssertNotCalled(suite.T(), "Down")
	getMigratorMock(migrations[2]).AssertCalled(suite.T(), "Down")
	getMigratorMock(migrations[3]).AssertNotCalled(suite.T(), "Down")

	suite.MigrationCRUDMock.AssertCalled(suite.T(), "DeleteMigrationByTimestamp", mock.MatchedBy(func(timestamp string) bool {
		return timestamp == latestTimestamp
	}))
}

func TestMigrationRunnerTestSuite(t *testing.T) {
	suite.Run(t, &MigrationRunnerTestSuite{})
}

func createMigration(timestamp string, upErr error, downErr error) migrationrunner.Migration {
	migratorMock := mocks.Migrator{}
	migratorMock.On("Up").Return(upErr)
	migratorMock.On("Down").Return(downErr)

	return migrationrunner.Migration{
		Timestamp:   timestamp,
		Description: "Does some migration stuff",
		Migrator:    &migratorMock,
	}
}

func createMigrations(timestamps ...string) []migrationrunner.Migration {
	migrations := make([]migrationrunner.Migration, len(timestamps))

	for i, timestamp := range timestamps {
		migrations[i] = createMigration(timestamp, nil, nil)
	}

	return migrations
}

func getMigratorMock(migration migrationrunner.Migration) *mocks.Migrator {
	return migration.Migrator.(*mocks.Migrator)
}
