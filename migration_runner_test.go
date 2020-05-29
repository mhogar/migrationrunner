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
}

func (suite *MigrationRunnerTestSuite) SetupTest() {
	suite.MigrationRepositoryMock = mocks.MigrationRepository{}
	suite.MigrationCRUDMock = mocks.MigrationCRUD{}
}

func (suite *MigrationRunnerTestSuite) TestMigrateUp_WithNoLatestTimestamp_RunsAllMigrations() {
	//arrange
	migrationMocks := createMigrationMocks("01", "04", "08", "10")

	migrations := make([]migrationrunner.Migration, len(migrationMocks))
	for i, _ := range migrationMocks {
		migrations[i] = &migrationMocks[i]
	}

	suite.MigrationRepositoryMock.On("GetMigrations").Return(migrations)
	suite.MigrationCRUDMock.On("GetLatestTimestamp").Return("100: this can be anything", false, nil)
	suite.MigrationCRUDMock.On("CreateMigration", mock.Anything).Return(nil)

	//act
	err := migrationrunner.MigrateUp(&suite.MigrationRepositoryMock, &suite.MigrationCRUDMock)

	//assert
	suite.NoError(err)

	migrationMocks[0].AssertCalled(suite.T(), "Up")
	migrationMocks[1].AssertCalled(suite.T(), "Up")
	migrationMocks[2].AssertCalled(suite.T(), "Up")
	migrationMocks[3].AssertCalled(suite.T(), "Up")

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
	migrationMocks := createMigrationMocks("01", "04", "08", "10")

	migrations := make([]migrationrunner.Migration, len(migrationMocks))
	for i, _ := range migrationMocks {
		migrations[i] = &migrationMocks[i]
	}

	suite.MigrationRepositoryMock.On("GetMigrations").Return(migrations)
	suite.MigrationCRUDMock.On("GetLatestTimestamp").Return("05", true, nil)
	suite.MigrationCRUDMock.On("CreateMigration", mock.Anything).Return(nil)

	//act
	err := migrationrunner.MigrateUp(&suite.MigrationRepositoryMock, &suite.MigrationCRUDMock)

	//assert
	suite.NoError(err)

	migrationMocks[0].AssertNotCalled(suite.T(), "Up")
	migrationMocks[1].AssertNotCalled(suite.T(), "Up")
	migrationMocks[2].AssertCalled(suite.T(), "Up")
	migrationMocks[3].AssertCalled(suite.T(), "Up")

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

	suite.MigrationRepositoryMock.On("GetMigrations").Return(nil)
	suite.MigrationCRUDMock.On("GetLatestTimestamp").Return("", false, errors.New(errMessage))

	//act
	err := migrationrunner.MigrateUp(&suite.MigrationRepositoryMock, &suite.MigrationCRUDMock)

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), errMessage)
}

func (suite *MigrationRunnerTestSuite) TestMigrateUp_WhereMigrationUpReturnsError_ReturnsError() {
	//arrange
	errMessage := "Up mock error"

	migrationMock := mocks.Migration{}
	migrationMock.On("GetTimestamp").Return("1")
	migrationMock.On("Up").Return(errors.New(errMessage))

	migrations := []migrationrunner.Migration{
		&migrationMock,
	}

	suite.MigrationRepositoryMock.On("GetMigrations").Return(migrations)
	suite.MigrationCRUDMock.On("GetLatestTimestamp").Return("0", true, nil)

	//act
	err := migrationrunner.MigrateUp(&suite.MigrationRepositoryMock, &suite.MigrationCRUDMock)

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), errMessage)
}

func (suite *MigrationRunnerTestSuite) TestMigrateUp_WhereCreateMigrationReturnsError_ReturnsError() {
	//arrange
	errMessage := "CreateMigration mock error"

	migrationMock := mocks.Migration{}
	migrationMock.On("GetTimestamp").Return("1")
	migrationMock.On("Up").Return(nil)

	migrations := []migrationrunner.Migration{
		&migrationMock,
	}

	suite.MigrationRepositoryMock.On("GetMigrations").Return(migrations)
	suite.MigrationCRUDMock.On("GetLatestTimestamp").Return("0", true, nil)
	suite.MigrationCRUDMock.On("CreateMigration", mock.Anything).Return(errors.New(errMessage))

	//act
	err := migrationrunner.MigrateUp(&suite.MigrationRepositoryMock, &suite.MigrationCRUDMock)

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), errMessage)
}

func (suite *MigrationRunnerTestSuite) TestMigrateDown_WhereGetLatestTimestampReturnsError_ReturnsError() {
	//arrange
	errMessage := "GetLatestTimestamp mock error"

	suite.MigrationRepositoryMock.On("GetMigrations").Return(nil)
	suite.MigrationCRUDMock.On("GetLatestTimestamp").Return("", false, errors.New(errMessage))

	//act
	err := migrationrunner.MigrateDown(&suite.MigrationRepositoryMock, &suite.MigrationCRUDMock)

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), errMessage)
}

func (suite *MigrationRunnerTestSuite) TestMigrateDown_WithNoLatestTimestamp_ReturnsError() {
	//arrange
	suite.MigrationRepositoryMock.On("GetMigrations").Return(nil)
	suite.MigrationCRUDMock.On("GetLatestTimestamp").Return("", false, nil)

	//act
	err := migrationrunner.MigrateDown(&suite.MigrationRepositoryMock, &suite.MigrationCRUDMock)

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), "no migrations")
}

func (suite *MigrationRunnerTestSuite) TestMigrateDown_WhereLatestTimestampNotFoundInMigrations_ReturnsError() {
	//arrange
	migrationMocks := createMigrationMocks("01", "04", "08", "10")

	migrations := make([]migrationrunner.Migration, len(migrationMocks))
	for i, _ := range migrationMocks {
		migrations[i] = &migrationMocks[i]
	}

	suite.MigrationRepositoryMock.On("GetMigrations").Return(migrations)
	suite.MigrationCRUDMock.On("GetLatestTimestamp").Return("05", true, nil)

	//act
	err := migrationrunner.MigrateDown(&suite.MigrationRepositoryMock, &suite.MigrationCRUDMock)

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), "could not find migration")
}

func (suite *MigrationRunnerTestSuite) TestMigrateDown_WhereMigrationDownReturnsError_ReturnError() {
	//arrange
	errMessage := "Down mock error"

	migrationMock := mocks.Migration{}
	migrationMock.On("GetTimestamp").Return("1")
	migrationMock.On("Down").Return(errors.New(errMessage))

	migrations := []migrationrunner.Migration{
		&migrationMock,
	}

	suite.MigrationRepositoryMock.On("GetMigrations").Return(migrations)
	suite.MigrationCRUDMock.On("GetLatestTimestamp").Return("1", true, nil)
	suite.MigrationCRUDMock.On("DeleteMigrationByTimestamp", mock.Anything).Return(nil)

	//act
	err := migrationrunner.MigrateDown(&suite.MigrationRepositoryMock, &suite.MigrationCRUDMock)

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), errMessage)
}

func (suite *MigrationRunnerTestSuite) TestMigrateDown_WhereDeleteMigrationReturnsError_ReturnError() {
	//arrange
	errMessage := "DeleteMigrationByTimestamp mock error"

	migrationMocks := createMigrationMocks("01", "04", "08", "10")

	migrations := make([]migrationrunner.Migration, len(migrationMocks))
	for i, _ := range migrationMocks {
		migrations[i] = &migrationMocks[i]
	}

	suite.MigrationRepositoryMock.On("GetMigrations").Return(migrations)
	suite.MigrationCRUDMock.On("GetLatestTimestamp").Return("08", true, nil)
	suite.MigrationCRUDMock.On("DeleteMigrationByTimestamp", mock.Anything).Return(errors.New(errMessage))

	//act
	err := migrationrunner.MigrateDown(&suite.MigrationRepositoryMock, &suite.MigrationCRUDMock)

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), errMessage)
}

func (suite *MigrationRunnerTestSuite) TestMigrateDown_RunsDownFunctionForMigrationMatchingLatestTimestamp() {
	//arrange
	latestTimestamp := "08"
	migrationMocks := createMigrationMocks("01", "04", "08", "10")

	migrations := make([]migrationrunner.Migration, len(migrationMocks))
	for i, _ := range migrationMocks {
		migrations[i] = &migrationMocks[i]
	}

	suite.MigrationRepositoryMock.On("GetMigrations").Return(migrations)
	suite.MigrationCRUDMock.On("GetLatestTimestamp").Return(latestTimestamp, true, nil)
	suite.MigrationCRUDMock.On("DeleteMigrationByTimestamp", mock.Anything).Return(nil)

	//act
	err := migrationrunner.MigrateDown(&suite.MigrationRepositoryMock, &suite.MigrationCRUDMock)

	//assert
	suite.NoError(err)

	migrationMocks[0].AssertNotCalled(suite.T(), "Down")
	migrationMocks[1].AssertNotCalled(suite.T(), "Down")
	migrationMocks[2].AssertCalled(suite.T(), "Down")
	migrationMocks[3].AssertNotCalled(suite.T(), "Down")

	suite.MigrationCRUDMock.AssertCalled(suite.T(), "DeleteMigrationByTimestamp", mock.MatchedBy(func(timestamp string) bool {
		return timestamp == latestTimestamp
	}))
}

func TestMigrationRunnerTestSuite(t *testing.T) {
	suite.Run(t, &MigrationRunnerTestSuite{})
}

func createMigrationMocks(timestamps ...string) []mocks.Migration {
	migrationMocks := make([]mocks.Migration, len(timestamps))

	for i, timestamp := range timestamps {
		migrationMocks[i] = mocks.Migration{}
		migrationMocks[i].On("GetTimestamp").Return(timestamp)
		migrationMocks[i].On("Up").Return(nil)
		migrationMocks[i].On("Down").Return(nil)
	}

	return migrationMocks
}
