package cachemanager

import (
	"context"
	"testing"

	commonconfig "github.com/rom8726/warden/internal/common/config"
	"github.com/rom8726/warden/internal/envelope-consumer/contract"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewManager(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		config  commonconfig.CacheConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: commonconfig.CacheConfig{
				ReleaseCacheSize:      100,
				IssueCacheSize:        100,
				IssueReleaseCacheSize: 100,
			},
			wantErr: false,
		},
		{
			name: "zero sizes",
			config: commonconfig.CacheConfig{
				ReleaseCacheSize:      0,
				IssueCacheSize:        0,
				IssueReleaseCacheSize: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			manager, err := New(&tt.config)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, manager)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, manager)
			}
		})
	}
}

func TestManager_ReleaseCache(t *testing.T) {
	t.Parallel()

	cacheConfig := commonconfig.CacheConfig{
		ReleaseCacheSize:      100,
		IssueCacheSize:        100,
		IssueReleaseCacheSize: 100,
	}

	manager, err := New(&cacheConfig)
	require.NoError(t, err)

	ctx := context.Background()

	// Test GetRelease - should return false for non-existent key
	value, found := manager.GetRelease(ctx, 1, "1.0.0")
	assert.False(t, found)
	assert.Equal(t, contract.ReleaseValue{}, value)

	// Test SetRelease
	err = manager.SetRelease(ctx, 1, "1.0.0", 123)
	assert.NoError(t, err)

	// Test GetRelease - should return true for existing key
	value, found = manager.GetRelease(ctx, 1, "1.0.0")
	assert.True(t, found)
	assert.Equal(t, uint(123), value.ReleaseID)
}

func TestManager_IssueCache(t *testing.T) {
	t.Parallel()

	cacheConfig := commonconfig.CacheConfig{
		ReleaseCacheSize:      100,
		IssueCacheSize:        100,
		IssueReleaseCacheSize: 100,
	}

	manager, err := New(&cacheConfig)
	require.NoError(t, err)

	ctx := context.Background()

	// Test GetIssue - should return false for non-existent key
	value, found := manager.GetIssue(ctx, "fingerprint123")
	assert.False(t, found)
	assert.Equal(t, contract.IssueValue{}, value)

	// Test SetIssue
	err = manager.SetIssue(ctx, "fingerprint123", 456)
	assert.NoError(t, err)

	// Test GetIssue - should return true for existing key
	value, found = manager.GetIssue(ctx, "fingerprint123")
	assert.True(t, found)
	assert.Equal(t, uint(456), value.IssueID)
}

func TestManager_IssueReleaseCache(t *testing.T) {
	t.Parallel()

	cacheConfig := commonconfig.CacheConfig{
		ReleaseCacheSize:      100,
		IssueCacheSize:        100,
		IssueReleaseCacheSize: 100,
	}

	manager, err := New(&cacheConfig)
	require.NoError(t, err)

	ctx := context.Background()

	// Test GetIssueRelease - should return false for non-existent key
	value, found := manager.GetIssueRelease(ctx, 1, 2)
	assert.False(t, found)
	assert.Equal(t, contract.IssueReleaseValue{}, value)

	// Test SetIssueRelease
	err = manager.SetIssueRelease(ctx, 1, 2, 789, true)
	assert.NoError(t, err)

	// Test GetIssueRelease - should return true for existing key
	value, found = manager.GetIssueRelease(ctx, 1, 2)
	assert.True(t, found)
	assert.Equal(t, uint(789), value.IssueReleaseID)
	assert.True(t, value.FirstSeenIn)
}

func TestManager_Stats(t *testing.T) {
	t.Parallel()

	cacheConfig := commonconfig.CacheConfig{
		ReleaseCacheSize:      100,
		IssueCacheSize:        200,
		IssueReleaseCacheSize: 300,
	}

	manager, err := New(&cacheConfig)
	require.NoError(t, err)

	ctx := context.Background()

	// Add some data
	err = manager.SetRelease(ctx, 1, "1.0.0", 123)
	require.NoError(t, err)

	err = manager.SetIssue(ctx, "fp1", 456)
	require.NoError(t, err)

	err = manager.SetIssueRelease(ctx, 1, 2, 789, true)
	require.NoError(t, err)

	// Get stats
	stats := manager.Stats()

	// Check status
	assert.Equal(t, "active", stats["status"])

	// Check release cache stats
	releaseStats, ok := stats["release_cache"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, 1, releaseStats["size"])
	assert.Equal(t, 100, releaseStats["capacity"])

	// Check issue cache stats
	issueStats, ok := stats["issue_cache"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, 1, issueStats["size"])
	assert.Equal(t, 200, issueStats["capacity"])

	// Check issue release cache stats
	issueReleaseStats, ok := stats["issue_release_cache"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, 1, issueReleaseStats["size"])
	assert.Equal(t, 300, issueReleaseStats["capacity"])
}

func TestManager_Clear(t *testing.T) {
	t.Parallel()

	cacheConfig := commonconfig.CacheConfig{
		ReleaseCacheSize:      100,
		IssueCacheSize:        100,
		IssueReleaseCacheSize: 100,
	}

	manager, err := New(&cacheConfig)
	require.NoError(t, err)

	ctx := context.Background()

	// Add some data
	err = manager.SetRelease(ctx, 1, "1.0.0", 123)
	require.NoError(t, err)

	err = manager.SetIssue(ctx, "fp1", 456)
	require.NoError(t, err)

	err = manager.SetIssueRelease(ctx, 1, 2, 789, true)
	require.NoError(t, err)

	// Verify data exists
	_, found := manager.GetRelease(ctx, 1, "1.0.0")
	assert.True(t, found)

	_, found = manager.GetIssue(ctx, "fp1")
	assert.True(t, found)

	_, found = manager.GetIssueRelease(ctx, 1, 2)
	assert.True(t, found)

	// Clear all caches
	err = manager.Clear(ctx)
	assert.NoError(t, err)

	// Verify data is gone
	_, found = manager.GetRelease(ctx, 1, "1.0.0")
	assert.False(t, found)

	_, found = manager.GetIssue(ctx, "fp1")
	assert.False(t, found)

	_, found = manager.GetIssueRelease(ctx, 1, 2)
	assert.False(t, found)
}

func TestManager_Close(t *testing.T) {
	t.Parallel()

	cacheConfig := commonconfig.CacheConfig{
		ReleaseCacheSize:      100,
		IssueCacheSize:        100,
		IssueReleaseCacheSize: 100,
	}

	manager, err := New(&cacheConfig)
	require.NoError(t, err)

	ctx := context.Background()

	// Add some data
	err = manager.SetRelease(ctx, 1, "1.0.0", 123)
	require.NoError(t, err)

	// Close manager
	err = manager.Close(ctx)
	assert.NoError(t, err)

	// Try to get data after close - should return false
	_, found := manager.GetRelease(ctx, 1, "1.0.0")
	assert.False(t, found)

	// Try to set data after close - should return error
	err = manager.SetRelease(ctx, 2, "2.0.0", 456)
	assert.Error(t, err)

	// Stats should show closed status
	stats := manager.Stats()
	assert.Equal(t, "closed", stats["status"])

	// Close again should not error
	err = manager.Close(ctx)
	assert.NoError(t, err)
}

func TestCacheKey_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		key  interface{ String() string }
		want string
	}{
		{
			name: "release key",
			key:  ReleaseKey{ProjectID: 123, Version: "1.0.0"},
			want: "123:1.0.0",
		},
		{
			name: "issue key",
			key:  IssueKey{Fingerprint: "abc123"},
			want: "abc123",
		},
		{
			name: "issue release key",
			key:  IssueReleaseKey{IssueID: 456, ReleaseID: 789},
			want: "456:789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.key.String())
		})
	}
}

func TestCacheValue_IsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value interface{ IsValid() bool }
		want  bool
	}{
		{
			name:  "valid release value",
			value: ReleaseValue{ReleaseID: 123},
			want:  true,
		},
		{
			name:  "invalid release value",
			value: ReleaseValue{ReleaseID: 0},
			want:  false,
		},
		{
			name:  "valid issue value",
			value: IssueValue{IssueID: 456},
			want:  true,
		},
		{
			name:  "invalid issue value",
			value: IssueValue{IssueID: 0},
			want:  false,
		},
		{
			name:  "valid issue release value",
			value: IssueReleaseValue{IssueReleaseID: 789, FirstSeenIn: true},
			want:  true,
		},
		{
			name:  "invalid issue release value",
			value: IssueReleaseValue{IssueReleaseID: 0, FirstSeenIn: true},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.value.IsValid())
		})
	}
}
