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

func TestAuthMiddleware(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		setupMocks      func(mockTokenizer *mockcontract.MockTokenizer, mockUsersSrv *mockcontract.MockUsersUseCase)
		authHeader      string
		checkContext    bool
		expectedUserID  domain.UserID
		expectedIsSuper bool
	}{
		{
			name: "No auth header passes through",
			setupMocks: func(mockTokenizer *mockcontract.MockTokenizer, mockUsersSrv *mockcontract.MockUsersUseCase) {
				// No expectations, as the middleware should bypass the check
			},
			authHeader:      "",
			checkContext:    false,
			expectedUserID:  0,
			expectedIsSuper: false,
		},
		{
			name: "Non-bearer token passes through",
			setupMocks: func(mockTokenizer *mockcontract.MockTokenizer, mockUsersSrv *mockcontract.MockUsersUseCase) {
				// No expectations, as the middleware should bypass the check
			},
			authHeader:      "Basic dXNlcjpwYXNzd29yZA==",
			checkContext:    false,
			expectedUserID:  0,
			expectedIsSuper: false,
		},
		{
			name: "Invalid token passes through",
			setupMocks: func(mockTokenizer *mockcontract.MockTokenizer, mockUsersSrv *mockcontract.MockUsersUseCase) {
				mockTokenizer.EXPECT().VerifyToken("invalid-token", domain.TokenTypeAccess).
					Return((*domain.TokenClaims)(nil), errors.New("invalid token"))
			},
			authHeader:      "Bearer invalid-token",
			checkContext:    false,
			expectedUserID:  0,
			expectedIsSuper: false,
		},
		{
			name: "User not found passes through",
			setupMocks: func(mockTokenizer *mockcontract.MockTokenizer, mockUsersSrv *mockcontract.MockUsersUseCase) {
				mockTokenizer.EXPECT().VerifyToken("valid-token", domain.TokenTypeAccess).
					Return(&domain.TokenClaims{UserID: 123}, nil)
				mockUsersSrv.EXPECT().GetByID(mock.Anything, domain.UserID(123)).
					Return(domain.User{}, errors.New("user not found"))
			},
			authHeader:      "Bearer valid-token",
			checkContext:    false,
			expectedUserID:  0,
			expectedIsSuper: false,
		},
		{
			name: "Valid token and user sets context",
			setupMocks: func(mockTokenizer *mockcontract.MockTokenizer, mockUsersSrv *mockcontract.MockUsersUseCase) {
				mockTokenizer.EXPECT().VerifyToken("valid-token", domain.TokenTypeAccess).
					Return(&domain.TokenClaims{UserID: 123}, nil)
				mockUsersSrv.EXPECT().GetByID(mock.Anything, domain.UserID(123)).
					Return(domain.User{ID: 123, IsSuperuser: false}, nil)
			},
			authHeader:      "Bearer valid-token",
			checkContext:    true,
			expectedUserID:  123,
			expectedIsSuper: false,
		},
		{
			name: "Superuser token sets superuser flag",
			setupMocks: func(mockTokenizer *mockcontract.MockTokenizer, mockUsersSrv *mockcontract.MockUsersUseCase) {
				mockTokenizer.EXPECT().VerifyToken("super-token", domain.TokenTypeAccess).
					Return(&domain.TokenClaims{UserID: 456}, nil)
				mockUsersSrv.EXPECT().GetByID(mock.Anything, domain.UserID(456)).
					Return(domain.User{ID: 456, IsSuperuser: true}, nil)
			},
			authHeader:      "Bearer super-token",
			checkContext:    true,
			expectedUserID:  456,
			expectedIsSuper: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockTokenizer := mockcontract.NewMockTokenizer(t)
			mockUsersSrv := mockcontract.NewMockUsersUseCase(t)
			tt.setupMocks(mockTokenizer, mockUsersSrv)

			// Create a test handler that will be wrapped by the middleware
			var userIDFromContext domain.UserID
			var isSuperFromContext bool
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.checkContext {
					userIDFromContext = wardencontext.UserID(r.Context())
					isSuperFromContext = wardencontext.IsSuper(r.Context())
				}
				w.WriteHeader(http.StatusOK)
			})

			// Create the middleware
			middleware := AuthMiddleware(mockTokenizer, mockUsersSrv)
			handler := middleware(testHandler)

			// Create a test request
			req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()

			// Call the handler
			handler.ServeHTTP(rec, req)

			// Check the response
			require.Equal(t, http.StatusOK, rec.Code)

			// Check that the user ID and superuser flag were set in the context if expected
			if tt.checkContext {
				require.Equal(t, tt.expectedUserID, userIDFromContext)
				require.Equal(t, tt.expectedIsSuper, isSuperFromContext)
			}
		})
	}
}
