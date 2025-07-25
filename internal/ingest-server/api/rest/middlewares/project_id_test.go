package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	wardencontext "github.com/rom8726/warden/internal/context"
	"github.com/rom8726/warden/internal/domain"
)

func TestWithProjectID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		checkContext   bool
		expectedID     domain.ProjectID
	}{
		{
			name:           "Short path passes through",
			path:           "/api/test",
			expectedStatus: http.StatusOK,
			checkContext:   false,
		},
		{
			name:           "Non-matching path passes through",
			path:           "/api/123/other/endpoint",
			expectedStatus: http.StatusOK,
			checkContext:   false,
		},
		{
			name:           "Invalid project ID returns 400",
			path:           "/api/invalid/envelope/test",
			expectedStatus: http.StatusBadRequest,
			checkContext:   false,
		},
		{
			name:           "Valid envelope path sets project ID",
			path:           "/api/123/envelope/test",
			expectedStatus: http.StatusOK,
			checkContext:   true,
			expectedID:     123,
		},
		{
			name:           "Valid store path sets project ID",
			path:           "/api/456/store/test",
			expectedStatus: http.StatusOK,
			checkContext:   true,
			expectedID:     456,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a test handler that will be wrapped by the middleware
			var projectIDFromContext domain.ProjectID
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.checkContext {
					projectIDFromContext = wardencontext.ProjectID(r.Context())
				}
				w.WriteHeader(http.StatusOK)
			})

			// Create the middleware
			handler := WithProjectID(testHandler)

			// Create a test request
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()

			// Call the handler
			handler.ServeHTTP(rec, req)

			// Check the response
			require.Equal(t, tt.expectedStatus, rec.Code)

			// Check that the project ID was set in the context if expected
			if tt.checkContext && tt.expectedStatus == http.StatusOK {
				require.Equal(t, tt.expectedID, projectIDFromContext)
			}
		})
	}
}
