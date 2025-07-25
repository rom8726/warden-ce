// Code generated by mockery v2.53.3. DO NOT EDIT.

package mockcontract

import (
	context "context"

	domain "github.com/rom8726/warden/internal/domain"
	contract "github.com/rom8726/warden/internal/envelope-consumer/contract"

	mock "github.com/stretchr/testify/mock"
)

// MockCacheService is an autogenerated mock type for the CacheService type
type MockCacheService struct {
	mock.Mock
}

type MockCacheService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCacheService) EXPECT() *MockCacheService_Expecter {
	return &MockCacheService_Expecter{mock: &_m.Mock}
}

// Clear provides a mock function with given fields: ctx
func (_m *MockCacheService) Clear(ctx context.Context) error {
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

// MockCacheService_Clear_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Clear'
type MockCacheService_Clear_Call struct {
	*mock.Call
}

// Clear is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockCacheService_Expecter) Clear(ctx interface{}) *MockCacheService_Clear_Call {
	return &MockCacheService_Clear_Call{Call: _e.mock.On("Clear", ctx)}
}

func (_c *MockCacheService_Clear_Call) Run(run func(ctx context.Context)) *MockCacheService_Clear_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockCacheService_Clear_Call) Return(_a0 error) *MockCacheService_Clear_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCacheService_Clear_Call) RunAndReturn(run func(context.Context) error) *MockCacheService_Clear_Call {
	_c.Call.Return(run)
	return _c
}

// Close provides a mock function with given fields: ctx
func (_m *MockCacheService) Close(ctx context.Context) error {
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

// MockCacheService_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type MockCacheService_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockCacheService_Expecter) Close(ctx interface{}) *MockCacheService_Close_Call {
	return &MockCacheService_Close_Call{Call: _e.mock.On("Close", ctx)}
}

func (_c *MockCacheService_Close_Call) Run(run func(ctx context.Context)) *MockCacheService_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockCacheService_Close_Call) Return(_a0 error) *MockCacheService_Close_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCacheService_Close_Call) RunAndReturn(run func(context.Context) error) *MockCacheService_Close_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrCreateIssue provides a mock function with given fields: ctx, issue, issueRepo
func (_m *MockCacheService) GetOrCreateIssue(ctx context.Context, issue domain.IssueDTO, issueRepo contract.IssuesRepository) (domain.IssueUpsertResult, error) {
	ret := _m.Called(ctx, issue, issueRepo)

	if len(ret) == 0 {
		panic("no return value specified for GetOrCreateIssue")
	}

	var r0 domain.IssueUpsertResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.IssueDTO, contract.IssuesRepository) (domain.IssueUpsertResult, error)); ok {
		return rf(ctx, issue, issueRepo)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.IssueDTO, contract.IssuesRepository) domain.IssueUpsertResult); ok {
		r0 = rf(ctx, issue, issueRepo)
	} else {
		r0 = ret.Get(0).(domain.IssueUpsertResult)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.IssueDTO, contract.IssuesRepository) error); ok {
		r1 = rf(ctx, issue, issueRepo)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCacheService_GetOrCreateIssue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrCreateIssue'
type MockCacheService_GetOrCreateIssue_Call struct {
	*mock.Call
}

// GetOrCreateIssue is a helper method to define mock.On call
//   - ctx context.Context
//   - issue domain.IssueDTO
//   - issueRepo contract.IssuesRepository
func (_e *MockCacheService_Expecter) GetOrCreateIssue(ctx interface{}, issue interface{}, issueRepo interface{}) *MockCacheService_GetOrCreateIssue_Call {
	return &MockCacheService_GetOrCreateIssue_Call{Call: _e.mock.On("GetOrCreateIssue", ctx, issue, issueRepo)}
}

func (_c *MockCacheService_GetOrCreateIssue_Call) Run(run func(ctx context.Context, issue domain.IssueDTO, issueRepo contract.IssuesRepository)) *MockCacheService_GetOrCreateIssue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.IssueDTO), args[2].(contract.IssuesRepository))
	})
	return _c
}

func (_c *MockCacheService_GetOrCreateIssue_Call) Return(_a0 domain.IssueUpsertResult, _a1 error) *MockCacheService_GetOrCreateIssue_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCacheService_GetOrCreateIssue_Call) RunAndReturn(run func(context.Context, domain.IssueDTO, contract.IssuesRepository) (domain.IssueUpsertResult, error)) *MockCacheService_GetOrCreateIssue_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrCreateIssueRelease provides a mock function with given fields: ctx, issueID, releaseID, firstSeenIn, issueReleaseRepo
func (_m *MockCacheService) GetOrCreateIssueRelease(ctx context.Context, issueID domain.IssueID, releaseID domain.ReleaseID, firstSeenIn bool, issueReleaseRepo contract.IssueReleasesRepository) error {
	ret := _m.Called(ctx, issueID, releaseID, firstSeenIn, issueReleaseRepo)

	if len(ret) == 0 {
		panic("no return value specified for GetOrCreateIssueRelease")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.IssueID, domain.ReleaseID, bool, contract.IssueReleasesRepository) error); ok {
		r0 = rf(ctx, issueID, releaseID, firstSeenIn, issueReleaseRepo)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockCacheService_GetOrCreateIssueRelease_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrCreateIssueRelease'
type MockCacheService_GetOrCreateIssueRelease_Call struct {
	*mock.Call
}

// GetOrCreateIssueRelease is a helper method to define mock.On call
//   - ctx context.Context
//   - issueID domain.IssueID
//   - releaseID domain.ReleaseID
//   - firstSeenIn bool
//   - issueReleaseRepo contract.IssueReleasesRepository
func (_e *MockCacheService_Expecter) GetOrCreateIssueRelease(ctx interface{}, issueID interface{}, releaseID interface{}, firstSeenIn interface{}, issueReleaseRepo interface{}) *MockCacheService_GetOrCreateIssueRelease_Call {
	return &MockCacheService_GetOrCreateIssueRelease_Call{Call: _e.mock.On("GetOrCreateIssueRelease", ctx, issueID, releaseID, firstSeenIn, issueReleaseRepo)}
}

func (_c *MockCacheService_GetOrCreateIssueRelease_Call) Run(run func(ctx context.Context, issueID domain.IssueID, releaseID domain.ReleaseID, firstSeenIn bool, issueReleaseRepo contract.IssueReleasesRepository)) *MockCacheService_GetOrCreateIssueRelease_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.IssueID), args[2].(domain.ReleaseID), args[3].(bool), args[4].(contract.IssueReleasesRepository))
	})
	return _c
}

func (_c *MockCacheService_GetOrCreateIssueRelease_Call) Return(_a0 error) *MockCacheService_GetOrCreateIssueRelease_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCacheService_GetOrCreateIssueRelease_Call) RunAndReturn(run func(context.Context, domain.IssueID, domain.ReleaseID, bool, contract.IssueReleasesRepository) error) *MockCacheService_GetOrCreateIssueRelease_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrCreateRelease provides a mock function with given fields: ctx, projectID, version, releaseRepo
func (_m *MockCacheService) GetOrCreateRelease(ctx context.Context, projectID domain.ProjectID, version string, releaseRepo contract.ReleaseRepository) (domain.ReleaseID, error) {
	ret := _m.Called(ctx, projectID, version, releaseRepo)

	if len(ret) == 0 {
		panic("no return value specified for GetOrCreateRelease")
	}

	var r0 domain.ReleaseID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.ProjectID, string, contract.ReleaseRepository) (domain.ReleaseID, error)); ok {
		return rf(ctx, projectID, version, releaseRepo)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.ProjectID, string, contract.ReleaseRepository) domain.ReleaseID); ok {
		r0 = rf(ctx, projectID, version, releaseRepo)
	} else {
		r0 = ret.Get(0).(domain.ReleaseID)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.ProjectID, string, contract.ReleaseRepository) error); ok {
		r1 = rf(ctx, projectID, version, releaseRepo)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCacheService_GetOrCreateRelease_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrCreateRelease'
type MockCacheService_GetOrCreateRelease_Call struct {
	*mock.Call
}

// GetOrCreateRelease is a helper method to define mock.On call
//   - ctx context.Context
//   - projectID domain.ProjectID
//   - version string
//   - releaseRepo contract.ReleaseRepository
func (_e *MockCacheService_Expecter) GetOrCreateRelease(ctx interface{}, projectID interface{}, version interface{}, releaseRepo interface{}) *MockCacheService_GetOrCreateRelease_Call {
	return &MockCacheService_GetOrCreateRelease_Call{Call: _e.mock.On("GetOrCreateRelease", ctx, projectID, version, releaseRepo)}
}

func (_c *MockCacheService_GetOrCreateRelease_Call) Run(run func(ctx context.Context, projectID domain.ProjectID, version string, releaseRepo contract.ReleaseRepository)) *MockCacheService_GetOrCreateRelease_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.ProjectID), args[2].(string), args[3].(contract.ReleaseRepository))
	})
	return _c
}

func (_c *MockCacheService_GetOrCreateRelease_Call) Return(_a0 domain.ReleaseID, _a1 error) *MockCacheService_GetOrCreateRelease_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCacheService_GetOrCreateRelease_Call) RunAndReturn(run func(context.Context, domain.ProjectID, string, contract.ReleaseRepository) (domain.ReleaseID, error)) *MockCacheService_GetOrCreateRelease_Call {
	_c.Call.Return(run)
	return _c
}

// Stats provides a mock function with no fields
func (_m *MockCacheService) Stats() map[string]interface{} {
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

// MockCacheService_Stats_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Stats'
type MockCacheService_Stats_Call struct {
	*mock.Call
}

// Stats is a helper method to define mock.On call
func (_e *MockCacheService_Expecter) Stats() *MockCacheService_Stats_Call {
	return &MockCacheService_Stats_Call{Call: _e.mock.On("Stats")}
}

func (_c *MockCacheService_Stats_Call) Run(run func()) *MockCacheService_Stats_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCacheService_Stats_Call) Return(_a0 map[string]interface{}) *MockCacheService_Stats_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCacheService_Stats_Call) RunAndReturn(run func() map[string]interface{}) *MockCacheService_Stats_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockCacheService creates a new instance of MockCacheService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCacheService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCacheService {
	mock := &MockCacheService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
