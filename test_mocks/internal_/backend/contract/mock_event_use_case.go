// Code generated by mockery v2.53.3. DO NOT EDIT.

package mockcontract

import (
	context "context"

	domain "github.com/rom8726/warden/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// MockEventUseCase is an autogenerated mock type for the EventUseCase type
type MockEventUseCase struct {
	mock.Mock
}

type MockEventUseCase_Expecter struct {
	mock *mock.Mock
}

func (_m *MockEventUseCase) EXPECT() *MockEventUseCase_Expecter {
	return &MockEventUseCase_Expecter{mock: &_m.Mock}
}

// IssueTimeseries provides a mock function with given fields: ctx, filter
func (_m *MockEventUseCase) IssueTimeseries(ctx context.Context, filter *domain.IssueEventsTimeseriesFilter) ([]domain.Timeseries, error) {
	ret := _m.Called(ctx, filter)

	if len(ret) == 0 {
		panic("no return value specified for IssueTimeseries")
	}

	var r0 []domain.Timeseries
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.IssueEventsTimeseriesFilter) ([]domain.Timeseries, error)); ok {
		return rf(ctx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *domain.IssueEventsTimeseriesFilter) []domain.Timeseries); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Timeseries)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *domain.IssueEventsTimeseriesFilter) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockEventUseCase_IssueTimeseries_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IssueTimeseries'
type MockEventUseCase_IssueTimeseries_Call struct {
	*mock.Call
}

// IssueTimeseries is a helper method to define mock.On call
//   - ctx context.Context
//   - filter *domain.IssueEventsTimeseriesFilter
func (_e *MockEventUseCase_Expecter) IssueTimeseries(ctx interface{}, filter interface{}) *MockEventUseCase_IssueTimeseries_Call {
	return &MockEventUseCase_IssueTimeseries_Call{Call: _e.mock.On("IssueTimeseries", ctx, filter)}
}

func (_c *MockEventUseCase_IssueTimeseries_Call) Run(run func(ctx context.Context, filter *domain.IssueEventsTimeseriesFilter)) *MockEventUseCase_IssueTimeseries_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*domain.IssueEventsTimeseriesFilter))
	})
	return _c
}

func (_c *MockEventUseCase_IssueTimeseries_Call) Return(_a0 []domain.Timeseries, _a1 error) *MockEventUseCase_IssueTimeseries_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockEventUseCase_IssueTimeseries_Call) RunAndReturn(run func(context.Context, *domain.IssueEventsTimeseriesFilter) ([]domain.Timeseries, error)) *MockEventUseCase_IssueTimeseries_Call {
	_c.Call.Return(run)
	return _c
}

// Timeseries provides a mock function with given fields: ctx, filter
func (_m *MockEventUseCase) Timeseries(ctx context.Context, filter *domain.EventTimeseriesFilter) ([]domain.Timeseries, error) {
	ret := _m.Called(ctx, filter)

	if len(ret) == 0 {
		panic("no return value specified for Timeseries")
	}

	var r0 []domain.Timeseries
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.EventTimeseriesFilter) ([]domain.Timeseries, error)); ok {
		return rf(ctx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *domain.EventTimeseriesFilter) []domain.Timeseries); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Timeseries)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *domain.EventTimeseriesFilter) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockEventUseCase_Timeseries_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Timeseries'
type MockEventUseCase_Timeseries_Call struct {
	*mock.Call
}

// Timeseries is a helper method to define mock.On call
//   - ctx context.Context
//   - filter *domain.EventTimeseriesFilter
func (_e *MockEventUseCase_Expecter) Timeseries(ctx interface{}, filter interface{}) *MockEventUseCase_Timeseries_Call {
	return &MockEventUseCase_Timeseries_Call{Call: _e.mock.On("Timeseries", ctx, filter)}
}

func (_c *MockEventUseCase_Timeseries_Call) Run(run func(ctx context.Context, filter *domain.EventTimeseriesFilter)) *MockEventUseCase_Timeseries_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*domain.EventTimeseriesFilter))
	})
	return _c
}

func (_c *MockEventUseCase_Timeseries_Call) Return(_a0 []domain.Timeseries, _a1 error) *MockEventUseCase_Timeseries_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockEventUseCase_Timeseries_Call) RunAndReturn(run func(context.Context, *domain.EventTimeseriesFilter) ([]domain.Timeseries, error)) *MockEventUseCase_Timeseries_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockEventUseCase creates a new instance of MockEventUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockEventUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockEventUseCase {
	mock := &MockEventUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
