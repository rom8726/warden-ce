// Code generated by mockery v2.53.3. DO NOT EDIT.

package mockcontract

import (
	context "context"

	contract "github.com/rom8726/warden/internal/envelope-consumer/contract"
	mock "github.com/stretchr/testify/mock"
)

// MockCacheManager is an autogenerated mock type for the CacheManager type
type MockCacheManager struct {
	mock.Mock
}

type MockCacheManager_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCacheManager) EXPECT() *MockCacheManager_Expecter {
	return &MockCacheManager_Expecter{mock: &_m.Mock}
}

// Clear provides a mock function with given fields: ctx
func (_m *MockCacheManager) Clear(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Clear")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockCacheManager_Clear_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Clear'
type MockCacheManager_Clear_Call struct {
	*mock.Call
}

// Clear is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockCacheManager_Expecter) Clear(ctx interface{}) *MockCacheManager_Clear_Call {
	return &MockCacheManager_Clear_Call{Call: _e.mock.On("Clear", ctx)}
}

func (_c *MockCacheManager_Clear_Call) Run(run func(ctx context.Context)) *MockCacheManager_Clear_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockCacheManager_Clear_Call) Return(_a0 error) *MockCacheManager_Clear_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCacheManager_Clear_Call) RunAndReturn(run func(context.Context) error) *MockCacheManager_Clear_Call {
	_c.Call.Return(run)
	return _c
}

// Close provides a mock function with given fields: ctx
func (_m *MockCacheManager) Close(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockCacheManager_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type MockCacheManager_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockCacheManager_Expecter) Close(ctx interface{}) *MockCacheManager_Close_Call {
	return &MockCacheManager_Close_Call{Call: _e.mock.On("Close", ctx)}
}

func (_c *MockCacheManager_Close_Call) Run(run func(ctx context.Context)) *MockCacheManager_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockCacheManager_Close_Call) Return(_a0 error) *MockCacheManager_Close_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCacheManager_Close_Call) RunAndReturn(run func(context.Context) error) *MockCacheManager_Close_Call {
	_c.Call.Return(run)
	return _c
}

// GetIssue provides a mock function with given fields: ctx, fingerprint
func (_m *MockCacheManager) GetIssue(ctx context.Context, fingerprint string) (contract.IssueValue, bool) {
	ret := _m.Called(ctx, fingerprint)

	if len(ret) == 0 {
		panic("no return value specified for GetIssue")
	}

	var r0 contract.IssueValue
	var r1 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) (contract.IssueValue, bool)); ok {
		return rf(ctx, fingerprint)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) contract.IssueValue); ok {
		r0 = rf(ctx, fingerprint)
	} else {
		r0 = ret.Get(0).(contract.IssueValue)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) bool); ok {
		r1 = rf(ctx, fingerprint)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// MockCacheManager_GetIssue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetIssue'
type MockCacheManager_GetIssue_Call struct {
	*mock.Call
}

// GetIssue is a helper method to define mock.On call
//   - ctx context.Context
//   - fingerprint string
func (_e *MockCacheManager_Expecter) GetIssue(ctx interface{}, fingerprint interface{}) *MockCacheManager_GetIssue_Call {
	return &MockCacheManager_GetIssue_Call{Call: _e.mock.On("GetIssue", ctx, fingerprint)}
}

func (_c *MockCacheManager_GetIssue_Call) Run(run func(ctx context.Context, fingerprint string)) *MockCacheManager_GetIssue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockCacheManager_GetIssue_Call) Return(_a0 contract.IssueValue, _a1 bool) *MockCacheManager_GetIssue_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCacheManager_GetIssue_Call) RunAndReturn(run func(context.Context, string) (contract.IssueValue, bool)) *MockCacheManager_GetIssue_Call {
	_c.Call.Return(run)
	return _c
}

// GetIssueRelease provides a mock function with given fields: ctx, issueID, releaseID
func (_m *MockCacheManager) GetIssueRelease(ctx context.Context, issueID uint, releaseID uint) (contract.IssueReleaseValue, bool) {
	ret := _m.Called(ctx, issueID, releaseID)

	if len(ret) == 0 {
		panic("no return value specified for GetIssueRelease")
	}

	var r0 contract.IssueReleaseValue
	var r1 bool
	if rf, ok := ret.Get(0).(func(context.Context, uint, uint) (contract.IssueReleaseValue, bool)); ok {
		return rf(ctx, issueID, releaseID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint, uint) contract.IssueReleaseValue); ok {
		r0 = rf(ctx, issueID, releaseID)
	} else {
		r0 = ret.Get(0).(contract.IssueReleaseValue)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint, uint) bool); ok {
		r1 = rf(ctx, issueID, releaseID)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// MockCacheManager_GetIssueRelease_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetIssueRelease'
type MockCacheManager_GetIssueRelease_Call struct {
	*mock.Call
}

// GetIssueRelease is a helper method to define mock.On call
//   - ctx context.Context
//   - issueID uint
//   - releaseID uint
func (_e *MockCacheManager_Expecter) GetIssueRelease(ctx interface{}, issueID interface{}, releaseID interface{}) *MockCacheManager_GetIssueRelease_Call {
	return &MockCacheManager_GetIssueRelease_Call{Call: _e.mock.On("GetIssueRelease", ctx, issueID, releaseID)}
}

func (_c *MockCacheManager_GetIssueRelease_Call) Run(run func(ctx context.Context, issueID uint, releaseID uint)) *MockCacheManager_GetIssueRelease_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint), args[2].(uint))
	})
	return _c
}

func (_c *MockCacheManager_GetIssueRelease_Call) Return(_a0 contract.IssueReleaseValue, _a1 bool) *MockCacheManager_GetIssueRelease_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCacheManager_GetIssueRelease_Call) RunAndReturn(run func(context.Context, uint, uint) (contract.IssueReleaseValue, bool)) *MockCacheManager_GetIssueRelease_Call {
	_c.Call.Return(run)
	return _c
}

// GetRelease provides a mock function with given fields: ctx, projectID, version
func (_m *MockCacheManager) GetRelease(ctx context.Context, projectID uint, version string) (contract.ReleaseValue, bool) {
	ret := _m.Called(ctx, projectID, version)

	if len(ret) == 0 {
		panic("no return value specified for GetRelease")
	}

	var r0 contract.ReleaseValue
	var r1 bool
	if rf, ok := ret.Get(0).(func(context.Context, uint, string) (contract.ReleaseValue, bool)); ok {
		return rf(ctx, projectID, version)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint, string) contract.ReleaseValue); ok {
		r0 = rf(ctx, projectID, version)
	} else {
		r0 = ret.Get(0).(contract.ReleaseValue)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint, string) bool); ok {
		r1 = rf(ctx, projectID, version)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// MockCacheManager_GetRelease_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRelease'
type MockCacheManager_GetRelease_Call struct {
	*mock.Call
}

// GetRelease is a helper method to define mock.On call
//   - ctx context.Context
//   - projectID uint
//   - version string
func (_e *MockCacheManager_Expecter) GetRelease(ctx interface{}, projectID interface{}, version interface{}) *MockCacheManager_GetRelease_Call {
	return &MockCacheManager_GetRelease_Call{Call: _e.mock.On("GetRelease", ctx, projectID, version)}
}

func (_c *MockCacheManager_GetRelease_Call) Run(run func(ctx context.Context, projectID uint, version string)) *MockCacheManager_GetRelease_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint), args[2].(string))
	})
	return _c
}

func (_c *MockCacheManager_GetRelease_Call) Return(_a0 contract.ReleaseValue, _a1 bool) *MockCacheManager_GetRelease_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCacheManager_GetRelease_Call) RunAndReturn(run func(context.Context, uint, string) (contract.ReleaseValue, bool)) *MockCacheManager_GetRelease_Call {
	_c.Call.Return(run)
	return _c
}

// SetIssue provides a mock function with given fields: ctx, fingerprint, issueID
func (_m *MockCacheManager) SetIssue(ctx context.Context, fingerprint string, issueID uint) error {
	ret := _m.Called(ctx, fingerprint, issueID)

	if len(ret) == 0 {
		panic("no return value specified for SetIssue")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, uint) error); ok {
		r0 = rf(ctx, fingerprint, issueID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockCacheManager_SetIssue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetIssue'
type MockCacheManager_SetIssue_Call struct {
	*mock.Call
}

// SetIssue is a helper method to define mock.On call
//   - ctx context.Context
//   - fingerprint string
//   - issueID uint
func (_e *MockCacheManager_Expecter) SetIssue(ctx interface{}, fingerprint interface{}, issueID interface{}) *MockCacheManager_SetIssue_Call {
	return &MockCacheManager_SetIssue_Call{Call: _e.mock.On("SetIssue", ctx, fingerprint, issueID)}
}

func (_c *MockCacheManager_SetIssue_Call) Run(run func(ctx context.Context, fingerprint string, issueID uint)) *MockCacheManager_SetIssue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(uint))
	})
	return _c
}

func (_c *MockCacheManager_SetIssue_Call) Return(_a0 error) *MockCacheManager_SetIssue_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCacheManager_SetIssue_Call) RunAndReturn(run func(context.Context, string, uint) error) *MockCacheManager_SetIssue_Call {
	_c.Call.Return(run)
	return _c
}

// SetIssueRelease provides a mock function with given fields: ctx, issueID, releaseID, issueReleaseID, firstSeenIn
func (_m *MockCacheManager) SetIssueRelease(ctx context.Context, issueID uint, releaseID uint, issueReleaseID uint, firstSeenIn bool) error {
	ret := _m.Called(ctx, issueID, releaseID, issueReleaseID, firstSeenIn)

	if len(ret) == 0 {
		panic("no return value specified for SetIssueRelease")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, uint, uint, bool) error); ok {
		r0 = rf(ctx, issueID, releaseID, issueReleaseID, firstSeenIn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockCacheManager_SetIssueRelease_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetIssueRelease'
type MockCacheManager_SetIssueRelease_Call struct {
	*mock.Call
}

// SetIssueRelease is a helper method to define mock.On call
//   - ctx context.Context
//   - issueID uint
//   - releaseID uint
//   - issueReleaseID uint
//   - firstSeenIn bool
func (_e *MockCacheManager_Expecter) SetIssueRelease(ctx interface{}, issueID interface{}, releaseID interface{}, issueReleaseID interface{}, firstSeenIn interface{}) *MockCacheManager_SetIssueRelease_Call {
	return &MockCacheManager_SetIssueRelease_Call{Call: _e.mock.On("SetIssueRelease", ctx, issueID, releaseID, issueReleaseID, firstSeenIn)}
}

func (_c *MockCacheManager_SetIssueRelease_Call) Run(run func(ctx context.Context, issueID uint, releaseID uint, issueReleaseID uint, firstSeenIn bool)) *MockCacheManager_SetIssueRelease_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint), args[2].(uint), args[3].(uint), args[4].(bool))
	})
	return _c
}

func (_c *MockCacheManager_SetIssueRelease_Call) Return(_a0 error) *MockCacheManager_SetIssueRelease_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCacheManager_SetIssueRelease_Call) RunAndReturn(run func(context.Context, uint, uint, uint, bool) error) *MockCacheManager_SetIssueRelease_Call {
	_c.Call.Return(run)
	return _c
}

// SetRelease provides a mock function with given fields: ctx, projectID, version, releaseID
func (_m *MockCacheManager) SetRelease(ctx context.Context, projectID uint, version string, releaseID uint) error {
	ret := _m.Called(ctx, projectID, version, releaseID)

	if len(ret) == 0 {
		panic("no return value specified for SetRelease")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, string, uint) error); ok {
		r0 = rf(ctx, projectID, version, releaseID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockCacheManager_SetRelease_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetRelease'
type MockCacheManager_SetRelease_Call struct {
	*mock.Call
}

// SetRelease is a helper method to define mock.On call
//   - ctx context.Context
//   - projectID uint
//   - version string
//   - releaseID uint
func (_e *MockCacheManager_Expecter) SetRelease(ctx interface{}, projectID interface{}, version interface{}, releaseID interface{}) *MockCacheManager_SetRelease_Call {
	return &MockCacheManager_SetRelease_Call{Call: _e.mock.On("SetRelease", ctx, projectID, version, releaseID)}
}

func (_c *MockCacheManager_SetRelease_Call) Run(run func(ctx context.Context, projectID uint, version string, releaseID uint)) *MockCacheManager_SetRelease_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint), args[2].(string), args[3].(uint))
	})
	return _c
}

func (_c *MockCacheManager_SetRelease_Call) Return(_a0 error) *MockCacheManager_SetRelease_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCacheManager_SetRelease_Call) RunAndReturn(run func(context.Context, uint, string, uint) error) *MockCacheManager_SetRelease_Call {
	_c.Call.Return(run)
	return _c
}

// Stats provides a mock function with no fields
func (_m *MockCacheManager) Stats() map[string]interface{} {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Stats")
	}

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func() map[string]interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	return r0
}

// MockCacheManager_Stats_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Stats'
type MockCacheManager_Stats_Call struct {
	*mock.Call
}

// Stats is a helper method to define mock.On call
func (_e *MockCacheManager_Expecter) Stats() *MockCacheManager_Stats_Call {
	return &MockCacheManager_Stats_Call{Call: _e.mock.On("Stats")}
}

func (_c *MockCacheManager_Stats_Call) Run(run func()) *MockCacheManager_Stats_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCacheManager_Stats_Call) Return(_a0 map[string]interface{}) *MockCacheManager_Stats_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCacheManager_Stats_Call) RunAndReturn(run func() map[string]interface{}) *MockCacheManager_Stats_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockCacheManager creates a new instance of MockCacheManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCacheManager(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCacheManager {
	mock := &MockCacheManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
