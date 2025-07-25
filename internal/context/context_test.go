package context

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rom8726/warden/internal/domain"
)

func TestWithProjectID(t *testing.T) {
	t.Parallel()

	// Create a base context
	ctx := context.Background()

	// Add a project ID to the context
	projectID := domain.ProjectID(123)
	ctxWithProjectID := WithProjectID(ctx, projectID)

	// Verify the context is not the same as the original
	require.NotEqual(t, ctx, ctxWithProjectID)

	// Retrieve the project ID from the context
	retrievedID := ProjectID(ctxWithProjectID)

	// Verify the retrieved ID matches the original
	require.Equal(t, projectID, retrievedID)
}

func TestProjectID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setupCtx    func() context.Context
		shouldPanic bool
	}{
		{
			name: "Context with project ID",
			setupCtx: func() context.Context {
				return WithProjectID(context.Background(), domain.ProjectID(123))
			},
			shouldPanic: false,
		},
		{
			name: "Context without project ID",
			setupCtx: func() context.Context {
				return context.Background()
			},
			shouldPanic: true,
		},
		{
			name: "Context with wrong type for project ID",
			setupCtx: func() context.Context {
				return context.WithValue(context.Background(), ctxKeyProjectID, "not a project ID")
			},
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := tt.setupCtx()

			if tt.shouldPanic {
				require.Panics(t, func() {
					_ = ProjectID(ctx)
				})
			} else {
				require.NotPanics(t, func() {
					_ = ProjectID(ctx)
				})
			}
		})
	}
}

func TestWithUserID(t *testing.T) {
	t.Parallel()

	// Create a base context
	ctx := context.Background()

	// Add a user ID to the context
	userID := domain.UserID(456)
	ctxWithUserID := WithUserID(ctx, userID)

	// Verify the context is not the same as the original
	require.NotEqual(t, ctx, ctxWithUserID)

	// Retrieve the user ID from the context
	retrievedID := UserID(ctxWithUserID)

	// Verify the retrieved ID matches the original
	require.Equal(t, userID, retrievedID)
}

func TestUserID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setupCtx       func() context.Context
		expectedUserID domain.UserID
	}{
		{
			name: "Context with user ID",
			setupCtx: func() context.Context {
				return WithUserID(context.Background(), domain.UserID(456))
			},
			expectedUserID: domain.UserID(456),
		},
		{
			name: "Context without user ID",
			setupCtx: func() context.Context {
				return context.Background()
			},
			expectedUserID: domain.UserID(0),
		},
		{
			name: "Context with wrong type for user ID",
			setupCtx: func() context.Context {
				return context.WithValue(context.Background(), ctxKeyUserID, "not a user ID")
			},
			expectedUserID: domain.UserID(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := tt.setupCtx()
			userID := UserID(ctx)

			require.Equal(t, tt.expectedUserID, userID)
		})
	}
}

func TestWithIsSuper(t *testing.T) {
	t.Parallel()

	// Create a base context
	ctx := context.Background()

	// Add isSuper flag to the context
	isSuper := true
	ctxWithIsSuper := WithIsSuper(ctx, isSuper)

	// Verify the context is not the same as the original
	require.NotEqual(t, ctx, ctxWithIsSuper)

	// Retrieve the isSuper flag from the context
	retrievedIsSuper := IsSuper(ctxWithIsSuper)

	// Verify the retrieved flag matches the original
	require.Equal(t, isSuper, retrievedIsSuper)
}

func TestIsSuper(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setupCtx      func() context.Context
		expectedValue bool
	}{
		{
			name: "Context with isSuper=true",
			setupCtx: func() context.Context {
				return WithIsSuper(context.Background(), true)
			},
			expectedValue: true,
		},
		{
			name: "Context with isSuper=false",
			setupCtx: func() context.Context {
				return WithIsSuper(context.Background(), false)
			},
			expectedValue: false,
		},
		{
			name: "Context without isSuper",
			setupCtx: func() context.Context {
				return context.Background()
			},
			expectedValue: false,
		},
		{
			name: "Context with wrong type for isSuper",
			setupCtx: func() context.Context {
				return context.WithValue(context.Background(), ctxKeyIsSuper, "not a bool")
			},
			expectedValue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := tt.setupCtx()
			isSuper := IsSuper(ctx)

			require.Equal(t, tt.expectedValue, isSuper)
		})
	}
}
