// Code generated by mockery v2.53.3. DO NOT EDIT.

package mockcontract

import (
	context "context"

	domain "github.com/rom8726/warden/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// MockEnvelopProducer is an autogenerated mock type for the EnvelopProducer type
type MockEnvelopProducer struct {
	mock.Mock
}

type MockEnvelopProducer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockEnvelopProducer) EXPECT() *MockEnvelopProducer_Expecter {
	return &MockEnvelopProducer_Expecter{mock: &_m.Mock}
}

// SendEnvelope provides a mock function with given fields: ctx, projectID, data
func (_m *MockEnvelopProducer) SendEnvelope(ctx context.Context, projectID domain.ProjectID, data []byte) error {
	ret := _m.Called(ctx, projectID, data)

	if len(ret) == 0 {
		panic("no return value specified for SendEnvelope")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.ProjectID, []byte) error); ok {
		r0 = rf(ctx, projectID, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockEnvelopProducer_SendEnvelope_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendEnvelope'
type MockEnvelopProducer_SendEnvelope_Call struct {
	*mock.Call
}

// SendEnvelope is a helper method to define mock.On call
//   - ctx context.Context
//   - projectID domain.ProjectID
//   - data []byte
func (_e *MockEnvelopProducer_Expecter) SendEnvelope(ctx interface{}, projectID interface{}, data interface{}) *MockEnvelopProducer_SendEnvelope_Call {
	return &MockEnvelopProducer_SendEnvelope_Call{Call: _e.mock.On("SendEnvelope", ctx, projectID, data)}
}

func (_c *MockEnvelopProducer_SendEnvelope_Call) Run(run func(ctx context.Context, projectID domain.ProjectID, data []byte)) *MockEnvelopProducer_SendEnvelope_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.ProjectID), args[2].([]byte))
	})
	return _c
}

func (_c *MockEnvelopProducer_SendEnvelope_Call) Return(_a0 error) *MockEnvelopProducer_SendEnvelope_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockEnvelopProducer_SendEnvelope_Call) RunAndReturn(run func(context.Context, domain.ProjectID, []byte) error) *MockEnvelopProducer_SendEnvelope_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockEnvelopProducer creates a new instance of MockEnvelopProducer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockEnvelopProducer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockEnvelopProducer {
	mock := &MockEnvelopProducer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
