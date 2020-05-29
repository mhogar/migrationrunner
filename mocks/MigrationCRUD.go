// Code generated by mockery v1.1.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// MigrationCRUD is an autogenerated mock type for the MigrationCRUD type
type MigrationCRUD struct {
	mock.Mock
}

// CreateMigration provides a mock function with given fields: timestamp
func (_m *MigrationCRUD) CreateMigration(timestamp string) error {
	ret := _m.Called(timestamp)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(timestamp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteMigrationByTimestamp provides a mock function with given fields: timestamp
func (_m *MigrationCRUD) DeleteMigrationByTimestamp(timestamp string) error {
	ret := _m.Called(timestamp)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(timestamp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetLatestTimestamp provides a mock function with given fields:
func (_m *MigrationCRUD) GetLatestTimestamp() (string, bool, error) {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func() bool); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(bool)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func() error); ok {
		r2 = rf()
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}