// Code generated by mockery v2.15.0. DO NOT EDIT.

package driver

import (
	context "context"

	repository "github.com/pokt-foundation/portal-db/repository"
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
func (_m *MockDriver) NotificationChannel() <-chan *repository.Notification {
	ret := _m.Called()

	var r0 <-chan *repository.Notification
	if rf, ok := ret.Get(0).(func() <-chan *repository.Notification); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan *repository.Notification)
		}
	}

	return r0
}

// ReadApplications provides a mock function with given fields: ctx
func (_m *MockDriver) ReadApplications(ctx context.Context) ([]*repository.Application, error) {
	ret := _m.Called(ctx)

	var r0 []*repository.Application
	if rf, ok := ret.Get(0).(func(context.Context) []*repository.Application); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*repository.Application)
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
func (_m *MockDriver) ReadBlockchains(ctx context.Context) ([]*repository.Blockchain, error) {
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

// ReadLoadBalancers provides a mock function with given fields: ctx
func (_m *MockDriver) ReadLoadBalancers(ctx context.Context) ([]*repository.LoadBalancer, error) {
	ret := _m.Called(ctx)

	var r0 []*repository.LoadBalancer
	if rf, ok := ret.Get(0).(func(context.Context) []*repository.LoadBalancer); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*repository.LoadBalancer)
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
func (_m *MockDriver) ReadPayPlans(ctx context.Context) ([]*repository.PayPlan, error) {
	ret := _m.Called(ctx)

	var r0 []*repository.PayPlan
	if rf, ok := ret.Get(0).(func(context.Context) []*repository.PayPlan); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*repository.PayPlan)
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

// UpdateAppFirstDateSurpassed provides a mock function with given fields: ctx, update
func (_m *MockDriver) UpdateAppFirstDateSurpassed(ctx context.Context, update *repository.UpdateFirstDateSurpassed) error {
	ret := _m.Called(ctx, update)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *repository.UpdateFirstDateSurpassed) error); ok {
		r0 = rf(ctx, update)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateApplication provides a mock function with given fields: ctx, id, update
func (_m *MockDriver) UpdateApplication(ctx context.Context, id string, update *repository.UpdateApplication) error {
	ret := _m.Called(ctx, id, update)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *repository.UpdateApplication) error); ok {
		r0 = rf(ctx, id, update)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateLoadBalancer provides a mock function with given fields: ctx, id, options
func (_m *MockDriver) UpdateLoadBalancer(ctx context.Context, id string, options *repository.UpdateLoadBalancer) error {
	ret := _m.Called(ctx, id, options)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *repository.UpdateLoadBalancer) error); ok {
		r0 = rf(ctx, id, options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WriteApplication provides a mock function with given fields: ctx, app
func (_m *MockDriver) WriteApplication(ctx context.Context, app *repository.Application) (*repository.Application, error) {
	ret := _m.Called(ctx, app)

	var r0 *repository.Application
	if rf, ok := ret.Get(0).(func(context.Context, *repository.Application) *repository.Application); ok {
		r0 = rf(ctx, app)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*repository.Application)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *repository.Application) error); ok {
		r1 = rf(ctx, app)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WriteBlockchain provides a mock function with given fields: ctx, blockchain
func (_m *MockDriver) WriteBlockchain(ctx context.Context, blockchain *repository.Blockchain) (*repository.Blockchain, error) {
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

// WriteLoadBalancer provides a mock function with given fields: ctx, loadBalancer
func (_m *MockDriver) WriteLoadBalancer(ctx context.Context, loadBalancer *repository.LoadBalancer) (*repository.LoadBalancer, error) {
	ret := _m.Called(ctx, loadBalancer)

	var r0 *repository.LoadBalancer
	if rf, ok := ret.Get(0).(func(context.Context, *repository.LoadBalancer) *repository.LoadBalancer); ok {
		r0 = rf(ctx, loadBalancer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*repository.LoadBalancer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *repository.LoadBalancer) error); ok {
		r1 = rf(ctx, loadBalancer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WriteRedirect provides a mock function with given fields: ctx, redirect
func (_m *MockDriver) WriteRedirect(ctx context.Context, redirect *repository.Redirect) (*repository.Redirect, error) {
	ret := _m.Called(ctx, redirect)

	var r0 *repository.Redirect
	if rf, ok := ret.Get(0).(func(context.Context, *repository.Redirect) *repository.Redirect); ok {
		r0 = rf(ctx, redirect)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*repository.Redirect)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *repository.Redirect) error); ok {
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
