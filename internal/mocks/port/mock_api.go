// Code generated by mockery v2.38.0. DO NOT EDIT.

package port

import (
	context "context"

	domain "github.com/ebisaan/inventory/internal/application/core/domain"
	mock "github.com/stretchr/testify/mock"
)

// MockAPI is an autogenerated mock type for the API type
type MockAPI struct {
	mock.Mock
}

type MockAPI_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAPI) EXPECT() *MockAPI_Expecter {
	return &MockAPI_Expecter{mock: &_m.Mock}
}

// GetProductByID provides a mock function with given fields: ctx, id
func (_m *MockAPI) GetProductByID(ctx context.Context, id int64) (*domain.Product, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetProductByID")
	}

	var r0 *domain.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*domain.Product, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *domain.Product); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAPI_GetProductByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetProductByID'
type MockAPI_GetProductByID_Call struct {
	*mock.Call
}

// GetProductByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id int64
func (_e *MockAPI_Expecter) GetProductByID(ctx interface{}, id interface{}) *MockAPI_GetProductByID_Call {
	return &MockAPI_GetProductByID_Call{Call: _e.mock.On("GetProductByID", ctx, id)}
}

func (_c *MockAPI_GetProductByID_Call) Run(run func(ctx context.Context, id int64)) *MockAPI_GetProductByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *MockAPI_GetProductByID_Call) Return(_a0 *domain.Product, _a1 error) *MockAPI_GetProductByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAPI_GetProductByID_Call) RunAndReturn(run func(context.Context, int64) (*domain.Product, error)) *MockAPI_GetProductByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetProducts provides a mock function with given fields: ctx, filter
func (_m *MockAPI) GetProducts(ctx context.Context, filter domain.Filter) ([]*domain.Product, domain.Metadata, error) {
	ret := _m.Called(ctx, filter)

	if len(ret) == 0 {
		panic("no return value specified for GetProducts")
	}

	var r0 []*domain.Product
	var r1 domain.Metadata
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Filter) ([]*domain.Product, domain.Metadata, error)); ok {
		return rf(ctx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.Filter) []*domain.Product); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.Filter) domain.Metadata); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Get(1).(domain.Metadata)
	}

	if rf, ok := ret.Get(2).(func(context.Context, domain.Filter) error); ok {
		r2 = rf(ctx, filter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockAPI_GetProducts_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetProducts'
type MockAPI_GetProducts_Call struct {
	*mock.Call
}

// GetProducts is a helper method to define mock.On call
//   - ctx context.Context
//   - filter domain.Filter
func (_e *MockAPI_Expecter) GetProducts(ctx interface{}, filter interface{}) *MockAPI_GetProducts_Call {
	return &MockAPI_GetProducts_Call{Call: _e.mock.On("GetProducts", ctx, filter)}
}

func (_c *MockAPI_GetProducts_Call) Run(run func(ctx context.Context, filter domain.Filter)) *MockAPI_GetProducts_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.Filter))
	})
	return _c
}

func (_c *MockAPI_GetProducts_Call) Return(_a0 []*domain.Product, _a1 domain.Metadata, _a2 error) *MockAPI_GetProducts_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockAPI_GetProducts_Call) RunAndReturn(run func(context.Context, domain.Filter) ([]*domain.Product, domain.Metadata, error)) *MockAPI_GetProducts_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAPI creates a new instance of MockAPI. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAPI(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAPI {
	mock := &MockAPI{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}