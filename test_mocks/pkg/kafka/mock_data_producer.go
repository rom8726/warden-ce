// Code generated by mockery v2.53.3. DO NOT EDIT.

package mockkafka

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockDataProducer is an autogenerated mock type for the DataProducer type
type MockDataProducer struct {
	mock.Mock
}

type MockDataProducer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDataProducer) EXPECT() *MockDataProducer_Expecter {
	return &MockDataProducer_Expecter{mock: &_m.Mock}
}

// Produce provides a mock function with given fields: ctx, data
func (_m *MockDataProducer) Produce(ctx context.Context, data []byte) error {
	ret := _m.Called(ctx, data)

	if len(ret) == 0 {
		panic("no return value specified for Produce")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []byte) error); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDataProducer_Produce_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Produce'
type MockDataProducer_Produce_Call struct {
	*mock.Call
}

// Produce is a helper method to define mock.On call
//   - ctx context.Context
//   - data []byte
func (_e *MockDataProducer_Expecter) Produce(ctx interface{}, data interface{}) *MockDataProducer_Produce_Call {
	return &MockDataProducer_Produce_Call{Call: _e.mock.On("Produce", ctx, data)}
}

func (_c *MockDataProducer_Produce_Call) Run(run func(ctx context.Context, data []byte)) *MockDataProducer_Produce_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]byte))
	})
	return _c
}

func (_c *MockDataProducer_Produce_Call) Return(_a0 error) *MockDataProducer_Produce_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDataProducer_Produce_Call) RunAndReturn(run func(context.Context, []byte) error) *MockDataProducer_Produce_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDataProducer creates a new instance of MockDataProducer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDataProducer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDataProducer {
	mock := &MockDataProducer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
