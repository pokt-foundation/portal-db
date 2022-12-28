// Code generated by mockery v2.15.0. DO NOT EDIT.

package driver

import (
	context "context"

	types "github.com/pokt-foundation/portal-db/types"
	mock "github.com/stretchr/testify/mock"
)

// MockDriver is an autogenerated mock type for the Driver type
type MockDriver struct {
	mock.Mock
}

// ActivateChain provides a mock function with given fields: ctx, id, active
func (_m *MockDriver) ActivateChain(ctx context.Context, id string, active bool) error {
	ret := _m.Called(ctx, id, active)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, bool) error); ok {
		r0 = rf(ctx, id, active)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NotificationChannel provides a mock function with given fields:
func (_m *MockDriver) NotificationChannel() <-chan *types.Notification {
	ret := _m.Called()

	var r0 <-chan *types.Notification
	if rf, ok := ret.Get(0).(func() <-chan *types.Notification); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan *types.Notification)
		}
	}

	return r0
}

// ReadApplications provides a mock function with given fields: ctx
func (_m *MockDriver) ReadApplications(ctx context.Context) ([]*types.Application, error) {
	ret := _m.Called(ctx)

	var r0 []*types.Application
	if rf, ok := ret.Get(0).(func(context.Context) []*types.Application); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Application)
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

// ReadBlockchains provides a mock function with given fields: ctx
func (_m *MockDriver) ReadBlockchains(ctx context.Context) ([]*types.Blockchain, error) {
	ret := _m.Called(ctx)

	var r0 []*types.Blockchain
	if rf, ok := ret.Get(0).(func(context.Context) []*types.Blockchain); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Blockchain)
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

// ReadLoadBalancers provides a mock function with given fields: ctx
func (_m *MockDriver) ReadLoadBalancers(ctx context.Context) ([]*types.LoadBalancer, error) {
	ret := _m.Called(ctx)

	var r0 []*types.LoadBalancer
	if rf, ok := ret.Get(0).(func(context.Context) []*types.LoadBalancer); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.LoadBalancer)
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

// ReadPayPlans provides a mock function with given fields: ctx
func (_m *MockDriver) ReadPayPlans(ctx context.Context) ([]*types.PayPlan, error) {
	ret := _m.Called(ctx)

	var r0 []*types.PayPlan
	if rf, ok := ret.Get(0).(func(context.Context) []*types.PayPlan); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.PayPlan)
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

// RemoveApplication provides a mock function with given fields: ctx, id
func (_m *MockDriver) RemoveApplication(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveLoadBalancer provides a mock function with given fields: ctx, id
func (_m *MockDriver) RemoveLoadBalancer(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveRedirect provides a mock function with given fields: ctx, blockchainID, domain
func (_m *MockDriver) RemoveRedirect(ctx context.Context, blockchainID string, domain string) error {
	ret := _m.Called(ctx, blockchainID, domain)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, blockchainID, domain)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveUserAccess provides a mock function with given fields: ctx, userID, lbID
func (_m *MockDriver) RemoveUserAccess(ctx context.Context, userID string, lbID string) error {
	ret := _m.Called(ctx, userID, lbID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, userID, lbID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateAppFirstDateSurpassed provides a mock function with given fields: ctx, update
func (_m *MockDriver) UpdateAppFirstDateSurpassed(ctx context.Context, update *types.UpdateFirstDateSurpassed) error {
	ret := _m.Called(ctx, update)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.UpdateFirstDateSurpassed) error); ok {
		r0 = rf(ctx, update)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateApplication provides a mock function with given fields: ctx, id, update
func (_m *MockDriver) UpdateApplication(ctx context.Context, id string, update *types.UpdateApplication) error {
	ret := _m.Called(ctx, id, update)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *types.UpdateApplication) error); ok {
		r0 = rf(ctx, id, update)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateChain provides a mock function with given fields: ctx, update
func (_m *MockDriver) UpdateChain(ctx context.Context, update *types.UpdateBlockchain) error {
	ret := _m.Called(ctx, update)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.UpdateBlockchain) error); ok {
		r0 = rf(ctx, update)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateLoadBalancer provides a mock function with given fields: ctx, id, options
func (_m *MockDriver) UpdateLoadBalancer(ctx context.Context, id string, options *types.UpdateLoadBalancer) error {
	ret := _m.Called(ctx, id, options)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *types.UpdateLoadBalancer) error); ok {
		r0 = rf(ctx, id, options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateUserAccessRole provides a mock function with given fields: ctx, userID, lbID, roleName
func (_m *MockDriver) UpdateUserAccessRole(ctx context.Context, userID string, lbID string, roleName types.RoleName) error {
	ret := _m.Called(ctx, userID, lbID, roleName)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, types.RoleName) error); ok {
		r0 = rf(ctx, userID, lbID, roleName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WriteApplication provides a mock function with given fields: ctx, app
func (_m *MockDriver) WriteApplication(ctx context.Context, app *types.Application) (*types.Application, error) {
	ret := _m.Called(ctx, app)

	var r0 *types.Application
	if rf, ok := ret.Get(0).(func(context.Context, *types.Application) *types.Application); ok {
		r0 = rf(ctx, app)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Application)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *types.Application) error); ok {
		r1 = rf(ctx, app)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WriteBlockchain provides a mock function with given fields: ctx, blockchain
func (_m *MockDriver) WriteBlockchain(ctx context.Context, blockchain *types.Blockchain) (*types.Blockchain, error) {
	ret := _m.Called(ctx, blockchain)

	var r0 *types.Blockchain
	if rf, ok := ret.Get(0).(func(context.Context, *types.Blockchain) *types.Blockchain); ok {
		r0 = rf(ctx, blockchain)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Blockchain)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *types.Blockchain) error); ok {
		r1 = rf(ctx, blockchain)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WriteLoadBalancer provides a mock function with given fields: ctx, loadBalancer
func (_m *MockDriver) WriteLoadBalancer(ctx context.Context, loadBalancer *types.LoadBalancer) (*types.LoadBalancer, error) {
	ret := _m.Called(ctx, loadBalancer)

	var r0 *types.LoadBalancer
	if rf, ok := ret.Get(0).(func(context.Context, *types.LoadBalancer) *types.LoadBalancer); ok {
		r0 = rf(ctx, loadBalancer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.LoadBalancer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *types.LoadBalancer) error); ok {
		r1 = rf(ctx, loadBalancer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WriteLoadBalancerUser provides a mock function with given fields: ctx, lbID, userAccess
func (_m *MockDriver) WriteLoadBalancerUser(ctx context.Context, lbID string, userAccess types.UserAccess) error {
	ret := _m.Called(ctx, lbID, userAccess)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, types.UserAccess) error); ok {
		r0 = rf(ctx, lbID, userAccess)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WriteRedirect provides a mock function with given fields: ctx, redirect
func (_m *MockDriver) WriteRedirect(ctx context.Context, redirect *types.Redirect) (*types.Redirect, error) {
	ret := _m.Called(ctx, redirect)

	var r0 *types.Redirect
	if rf, ok := ret.Get(0).(func(context.Context, *types.Redirect) *types.Redirect); ok {
		r0 = rf(ctx, redirect)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Redirect)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *types.Redirect) error); ok {
		r1 = rf(ctx, redirect)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockDriver interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockDriver creates a new instance of MockDriver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockDriver(t mockConstructorTestingTNewMockDriver) *MockDriver {
	mock := &MockDriver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
