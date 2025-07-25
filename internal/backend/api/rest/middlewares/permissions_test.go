package middlewares

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	wardencontext "github.com/rom8726/warden/internal/context"
	"github.com/rom8726/warden/internal/domain"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/backend/contract"
)

func TestProjectAccess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setupMock      func(mockSvc *mockcontract.MockPermissionsService)
		path           string
		method         string
		expectedStatus int
		checkContext   bool
	}{
		{
			name: "POST request bypasses permission check",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				// No expectations, as the middleware should bypass the check
			},
			path:           "/api/projects/123/something",
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
			checkContext:   false,
		},
		{
			name: "Empty project ID bypasses permission check",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				// No expectations, as the middleware should bypass the check
			},
			path:           "/api/projects//something",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			checkContext:   false,
		},
		{
			name: "Project ID 'add' bypasses permission check",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				// No expectations, as the middleware should bypass the check
			},
			path:           "/api/projects/add/something",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			checkContext:   false,
		},
		{
			name: "Invalid project ID returns 400",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				// No expectations, as the middleware should return before calling the service
			},
			path:           "/api/projects/invalid/something",
			method:         http.MethodGet,
			expectedStatus: http.StatusBadRequest,
			checkContext:   false,
		},
		{
			name: "Permission denied returns 403",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				mockSvc.EXPECT().CanAccessProject(mock.Anything, domain.ProjectID(123)).
					Return(domain.ErrPermissionDenied)
			},
			path:           "/api/projects/123/something",
			method:         http.MethodGet,
			expectedStatus: http.StatusForbidden,
			checkContext:   false,
		},
		{
			name: "User not found returns 401",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				mockSvc.EXPECT().CanAccessProject(mock.Anything, domain.ProjectID(123)).
					Return(domain.ErrUserNotFound)
			},
			path:           "/api/projects/123/something",
			method:         http.MethodGet,
			expectedStatus: http.StatusUnauthorized,
			checkContext:   false,
		},
		{
			name: "Other error returns 500",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				mockSvc.EXPECT().CanAccessProject(mock.Anything, domain.ProjectID(123)).
					Return(errors.New("some error"))
			},
			path:           "/api/projects/123/something",
			method:         http.MethodGet,
			expectedStatus: http.StatusInternalServerError,
			checkContext:   false,
		},
		{
			name: "Successful access sets project ID in context",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				mockSvc.EXPECT().CanAccessProject(mock.Anything, domain.ProjectID(123)).
					Return(nil)
			},
			path:           "/api/projects/123/something",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			checkContext:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a mock permissions service
			mockSvc := mockcontract.NewMockPermissionsService(t)
			tt.setupMock(mockSvc)

			// Create a test handler that will be wrapped by the middleware
			var projectIDFromContext domain.ProjectID
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.checkContext {
					projectIDFromContext = wardencontext.ProjectID(r.Context())
				}
				w.WriteHeader(http.StatusOK)
			})

			// Create the middleware
			middleware := ProjectAccess(mockSvc)
			handler := middleware(testHandler)

			// Create a test request
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			// Call the handler
			handler.ServeHTTP(rec, req)

			// Check the response
			require.Equal(t, tt.expectedStatus, rec.Code)

			// Check that the project ID was set in the context if expected
			if tt.checkContext && tt.expectedStatus == http.StatusOK {
				require.Equal(t, domain.ProjectID(123), projectIDFromContext)
			}
		})
	}
}

func TestProjectManagement(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setupMock      func(mockSvc *mockcontract.MockPermissionsService)
		path           string
		method         string
		expectedStatus int
		checkContext   bool
	}{
		{
			name: "GET request bypasses permission check",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				// No expectations, as the middleware should bypass the check
			},
			path:           "/api/projects/123/something",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			checkContext:   false,
		},
		{
			name: "Empty project ID bypasses permission check",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				// No expectations, as the middleware should bypass the check
			},
			path:           "/api/projects//something",
			method:         http.MethodPut,
			expectedStatus: http.StatusOK,
			checkContext:   false,
		},
		{
			name: "Project ID 'add' bypasses permission check",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				// No expectations, as the middleware should bypass the check
			},
			path:           "/api/projects/add/something",
			method:         http.MethodPut,
			expectedStatus: http.StatusOK,
			checkContext:   false,
		},
		{
			name: "Stats endpoint bypasses permission check",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				// No expectations, as the middleware should bypass the check
			},
			path:           "/api/projects/123/something/stats",
			method:         http.MethodPut,
			expectedStatus: http.StatusOK,
			checkContext:   false,
		},
		{
			name: "Invalid project ID returns 400",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				// No expectations, as the middleware should return before calling the service
			},
			path:           "/api/projects/invalid/something",
			method:         http.MethodPut,
			expectedStatus: http.StatusBadRequest,
			checkContext:   false,
		},
		{
			name: "Permission denied returns 403",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				mockSvc.EXPECT().CanManageProject(mock.Anything, domain.ProjectID(123), false).
					Return(domain.ErrPermissionDenied)
			},
			path:           "/api/projects/123/something",
			method:         http.MethodPut,
			expectedStatus: http.StatusForbidden,
			checkContext:   false,
		},
		{
			name: "User not found returns 401",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				mockSvc.EXPECT().CanManageProject(mock.Anything, domain.ProjectID(123), false).
					Return(domain.ErrUserNotFound)
			},
			path:           "/api/projects/123/something",
			method:         http.MethodPut,
			expectedStatus: http.StatusUnauthorized,
			checkContext:   false,
		},
		{
			name: "Other error returns 500",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				mockSvc.EXPECT().CanManageProject(mock.Anything, domain.ProjectID(123), false).
					Return(errors.New("some error"))
			},
			path:           "/api/projects/123/something",
			method:         http.MethodPut,
			expectedStatus: http.StatusInternalServerError,
			checkContext:   false,
		},
		{
			name: "Successful management sets project ID in context",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				mockSvc.EXPECT().CanManageProject(mock.Anything, domain.ProjectID(123), false).
					Return(nil)
			},
			path:           "/api/projects/123/something",
			method:         http.MethodPut,
			expectedStatus: http.StatusOK,
			checkContext:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a mock permissions service
			mockSvc := mockcontract.NewMockPermissionsService(t)
			tt.setupMock(mockSvc)

			// Create a test handler that will be wrapped by the middleware
			var projectIDFromContext domain.ProjectID
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.checkContext {
					projectIDFromContext = wardencontext.ProjectID(r.Context())
				}
				w.WriteHeader(http.StatusOK)
			})

			// Create the middleware
			middleware := ProjectManagement(mockSvc)
			handler := middleware(testHandler)

			// Create a test request
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			// Call the handler
			handler.ServeHTTP(rec, req)

			// Check the response
			require.Equal(t, tt.expectedStatus, rec.Code)

			// Check that the project ID was set in the context if expected
			if tt.checkContext && tt.expectedStatus == http.StatusOK {
				require.Equal(t, domain.ProjectID(123), projectIDFromContext)
			}
		})
	}
}

func TestIssueAccess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setupMock      func(mockSvc *mockcontract.MockPermissionsService)
		path           string
		method         string
		expectedStatus int
	}{
		{
			name: "Empty issue ID bypasses permission check",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				// No expectations, as the middleware should bypass the check
			},
			path:           "/api/projects/123/issues//something",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Issue ID 'recent' bypasses permission check",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				// No expectations, as the middleware should bypass the check
			},
			path:           "/api/projects/123/issues/recent/something",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Issue ID 'timeseries' bypasses permission check",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				// No expectations, as the middleware should bypass the check
			},
			path:           "/api/projects/123/issues/timeseries/something",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid issue ID returns 400",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				// No expectations, as the middleware should return before calling the service
			},
			path:           "/api/projects/123/issues/invalid/something",
			method:         http.MethodGet,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Permission denied returns 403",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				mockSvc.EXPECT().CanAccessIssue(mock.Anything, domain.IssueID(456)).
					Return(domain.ErrPermissionDenied)
			},
			path:           "/api/projects/123/issues/456/something",
			method:         http.MethodGet,
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "User not found returns 401",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				mockSvc.EXPECT().CanAccessIssue(mock.Anything, domain.IssueID(456)).
					Return(domain.ErrUserNotFound)
			},
			path:           "/api/projects/123/issues/456/something",
			method:         http.MethodGet,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Other error returns 500",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				mockSvc.EXPECT().CanAccessIssue(mock.Anything, domain.IssueID(456)).
					Return(errors.New("some error"))
			},
			path:           "/api/projects/123/issues/456/something",
			method:         http.MethodGet,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Successful access passes through",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				mockSvc.EXPECT().CanAccessIssue(mock.Anything, domain.IssueID(456)).
					Return(nil)
			},
			path:           "/api/projects/123/issues/456/something",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a mock permissions service
			mockSvc := mockcontract.NewMockPermissionsService(t)
			tt.setupMock(mockSvc)

			// Create a test handler that will be wrapped by the middleware
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Create the middleware
			middleware := IssueAccess(mockSvc)
			handler := middleware(testHandler)

			// Create a test request
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			// Call the handler
			handler.ServeHTTP(rec, req)

			// Check the response
			require.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestIssueManagement(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setupMock      func(mockSvc *mockcontract.MockPermissionsService)
		path           string
		method         string
		expectedStatus int
	}{
		{
			name: "GET request bypasses permission check",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				// No expectations, as the middleware should bypass the check
			},
			path:           "/api/projects/123/issues/456/something",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Empty issue ID bypasses permission check",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				// No expectations, as the middleware should bypass the check
			},
			path:           "/api/projects/123/issues//something",
			method:         http.MethodPut,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Issue ID 'recent' bypasses permission check",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				// No expectations, as the middleware should bypass the check
			},
			path:           "/api/projects/123/issues/recent/something",
			method:         http.MethodPut,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Issue ID 'timeseries' bypasses permission check",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				// No expectations, as the middleware should bypass the check
			},
			path:           "/api/projects/123/issues/timeseries/something",
			method:         http.MethodPut,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid issue ID returns 400",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				// No expectations, as the middleware should return before calling the service
			},
			path:           "/api/projects/123/issues/invalid/something",
			method:         http.MethodPut,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Permission denied returns 403",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				mockSvc.EXPECT().CanManageIssue(mock.Anything, domain.IssueID(456)).
					Return(domain.ErrPermissionDenied)
			},
			path:           "/api/projects/123/issues/456/something",
			method:         http.MethodPut,
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "User not found returns 401",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				mockSvc.EXPECT().CanManageIssue(mock.Anything, domain.IssueID(456)).
					Return(domain.ErrUserNotFound)
			},
			path:           "/api/projects/123/issues/456/something",
			method:         http.MethodPut,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Other error returns 500",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				mockSvc.EXPECT().CanManageIssue(mock.Anything, domain.IssueID(456)).
					Return(errors.New("some error"))
			},
			path:           "/api/projects/123/issues/456/something",
			method:         http.MethodPut,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Successful management passes through",
			setupMock: func(mockSvc *mockcontract.MockPermissionsService) {
				mockSvc.EXPECT().CanManageIssue(mock.Anything, domain.IssueID(456)).
					Return(nil)
			},
			path:           "/api/projects/123/issues/456/something",
			method:         http.MethodPut,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a mock permissions service
			mockSvc := mockcontract.NewMockPermissionsService(t)
			tt.setupMock(mockSvc)

			// Create a test handler that will be wrapped by the middleware
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Create the middleware
			middleware := IssueManagement(mockSvc)
			handler := middleware(testHandler)

			// Create a test request
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			// Call the handler
			handler.ServeHTTP(rec, req)

			// Check the response
			require.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
