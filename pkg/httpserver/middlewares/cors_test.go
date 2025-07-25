package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCORSMdw(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedNext   bool
	}{
		{
			name:           "OPTIONS request returns 204",
			method:         http.MethodOptions,
			expectedStatus: http.StatusNoContent,
			expectedNext:   false,
		},
		{
			name:           "GET request passes through",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedNext:   true,
		},
		{
			name:           "POST request passes through",
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
			expectedNext:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a test handler that will be wrapped by the middleware
			nextCalled := false
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			})

			// Create the middleware
			handler := CORSMdw(testHandler)

			// Create a test request
			req := httptest.NewRequest(tt.method, "/api/test", nil)
			rec := httptest.NewRecorder()

			// Call the handler
			handler.ServeHTTP(rec, req)

			// Check the response
			require.Equal(t, tt.expectedStatus, rec.Code)
			require.Equal(t, tt.expectedNext, nextCalled)

			// Check CORS headers
			require.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
			require.Equal(t, "GET, POST, PUT, PATCH, DELETE, OPTIONS", rec.Header().Get("Access-Control-Allow-Methods"))
			require.Equal(t, "Content-Type, Authorization", rec.Header().Get("Access-Control-Allow-Headers"))
		})
	}
}
