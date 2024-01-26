// Code generated by mockery v2.38.0. DO NOT EDIT.

package port

import (
	context "context"

	domain "github.com/ebisaan/inventory/internal/application/core/domain"
	mock "github.com/stretchr/testify/mock"
)

// MockProductDB is an autogenerated mock type for the ProductDB type
type MockProductDB struct {
	mock.Mock
}

type MockProductDB_Expecter struct {
	mock *mock.Mock
}

func (_m *MockProductDB) EXPECT() *MockProductDB_Expecter {
	return &MockProductDB_Expecter{mock: &_m.Mock}
}

// GetProductByID provides a mock function with given fields: ctx, id
func (_m *MockProductDB) GetProductByID(ctx context.Context, id int64) (*domain.Product, error) {
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

// MockProductDB_GetProductByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetProductByID'
type MockProductDB_GetProductByID_Call struct {
	*mock.Call
}

// GetProductByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id int64
func (_e *MockProductDB_Expecter) GetProductByID(ctx interface{}, id interface{}) *MockProductDB_GetProductByID_Call {
	return &MockProductDB_GetProductByID_Call{Call: _e.mock.On("GetProductByID", ctx, id)}
}

func (_c *MockProductDB_GetProductByID_Call) Run(run func(ctx context.Context, id int64)) *MockProductDB_GetProductByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *MockProductDB_GetProductByID_Call) Return(_a0 *domain.Product, _a1 error) *MockProductDB_GetProductByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockProductDB_GetProductByID_Call) RunAndReturn(run func(context.Context, int64) (*domain.Product, error)) *MockProductDB_GetProductByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetProducts provides a mock function with given fields: ctx, filter
func (_m *MockProductDB) GetProducts(ctx context.Context, filter domain.Filter) (int64, []*domain.Product, error) {
	ret := _m.Called(ctx, filter)

	if len(ret) == 0 {
		panic("no return value specified for GetProducts")
	}

	var r0 int64
	var r1 []*domain.Product
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Filter) (int64, []*domain.Product, error)); ok {
		return rf(ctx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.Filter) int64); ok {
		r0 = rf(ctx, filter)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.Filter) []*domain.Product); ok {
		r1 = rf(ctx, filter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]*domain.Product)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, domain.Filter) error); ok {
		r2 = rf(ctx, filter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockProductDB_GetProducts_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetProducts'
type MockProductDB_GetProducts_Call struct {
	*mock.Call
}

// GetProducts is a helper method to define mock.On call
//   - ctx context.Context
//   - filter domain.Filter
func (_e *MockProductDB_Expecter) GetProducts(ctx interface{}, filter interface{}) *MockProductDB_GetProducts_Call {
	return &MockProductDB_GetProducts_Call{Call: _e.mock.On("GetProducts", ctx, filter)}
}

func (_c *MockProductDB_GetProducts_Call) Run(run func(ctx context.Context, filter domain.Filter)) *MockProductDB_GetProducts_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.Filter))
	})
	return _c
}

func (_c *MockProductDB_GetProducts_Call) Return(_a0 int64, _a1 []*domain.Product, _a2 error) *MockProductDB_GetProducts_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockProductDB_GetProducts_Call) RunAndReturn(run func(context.Context, domain.Filter) (int64, []*domain.Product, error)) *MockProductDB_GetProducts_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockProductDB creates a new instance of MockProductDB. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockProductDB(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockProductDB {
	mock := &MockProductDB{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
