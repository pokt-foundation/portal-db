// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	context "context"

	driver "github.com/pokt-foundation/portal-db/driver"
	mock "github.com/stretchr/testify/mock"

	repository "github.com/pokt-foundation/portal-db/repository"
)

// IPostgresDriver is an autogenerated mock type for the IPostgresDriver type
type IPostgresDriver struct {
	mock.Mock
}

// ActivateBlockchain provides a mock function with given fields: ctx, arg
func (_m *IPostgresDriver) ActivateBlockchain(ctx context.Context, arg driver.ActivateBlockchainParams) error {
	ret := _m.Called(ctx, arg)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, driver.ActivateBlockchainParams) error); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ActivateChain provides a mock function with given fields: ctx, id, active
func (_m *IPostgresDriver) ActivateChain(ctx context.Context, id string, active bool) error {
	ret := _m.Called(ctx, id, active)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, bool) error); ok {
		r0 = rf(ctx, id, active)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertBlockchain provides a mock function with given fields: ctx, arg
func (_m *IPostgresDriver) InsertBlockchain(ctx context.Context, arg driver.InsertBlockchainParams) error {
	ret := _m.Called(ctx, arg)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, driver.InsertBlockchainParams) error); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertSyncCheckOptions provides a mock function with given fields: ctx, arg
func (_m *IPostgresDriver) InsertSyncCheckOptions(ctx context.Context, arg driver.InsertSyncCheckOptionsParams) error {
	ret := _m.Called(ctx, arg)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, driver.InsertSyncCheckOptionsParams) error); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReadBlockchains provides a mock function with given fields: ctx
func (_m *IPostgresDriver) ReadBlockchains(ctx context.Context) ([]*repository.Blockchain, error) {
	ret := _m.Called(ctx)

	var r0 []*repository.Blockchain
	if rf, ok := ret.Get(0).(func(context.Context) []*repository.Blockchain); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*repository.Blockchain)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SelectBlockchains provides a mock function with given fields: ctx
func (_m *IPostgresDriver) SelectBlockchains(ctx context.Context) ([]driver.SelectBlockchainsRow, error) {
	ret := _m.Called(ctx)

	var r0 []driver.SelectBlockchainsRow
	if rf, ok := ret.Get(0).(func(context.Context) []driver.SelectBlockchainsRow); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]driver.SelectBlockchainsRow)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WriteBlockchain provides a mock function with given fields: ctx, blockchain
func (_m *IPostgresDriver) WriteBlockchain(ctx context.Context, blockchain *repository.Blockchain) (*repository.Blockchain, error) {
	ret := _m.Called(ctx, blockchain)

	var r0 *repository.Blockchain
	if rf, ok := ret.Get(0).(func(context.Context, *repository.Blockchain) *repository.Blockchain); ok {
		r0 = rf(ctx, blockchain)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*repository.Blockchain)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *repository.Blockchain) error); ok {
		r1 = rf(ctx, blockchain)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewIPostgresDriver interface {
	mock.TestingT
	Cleanup(func())
}

// NewIPostgresDriver creates a new instance of IPostgresDriver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIPostgresDriver(t mockConstructorTestingTNewIPostgresDriver) *IPostgresDriver {
	mock := &IPostgresDriver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
